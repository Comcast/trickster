/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/tricksterproxy/trickster/pkg/proxy/headers"
	"github.com/tricksterproxy/trickster/pkg/timeseries"
	"github.com/tricksterproxy/trickster/pkg/timeseries/dataset"
	"github.com/tricksterproxy/trickster/pkg/timeseries/epoch"
)

// WFDocument the Wire Format Document for the timeseries
type WFDocument struct {
	Status string `json:"status"`
	Data   WFData `json:"data"`
}

// WFData is the data section of the WFD
type WFData struct {
	ResultType string     `json:"resultType"`
	Results    []WFResult `json:"result"`
}

// WFResult is the Result section of the WFD
type WFResult struct {
	Metric dataset.Tags    `json:"metric"`
	Values [][]interface{} `json:"values"`
	Value  []interface{}   `json:"value"`
}

// NewModeler returns a collection of modeling functions for prometheus interoperability
func NewModeler() *timeseries.Modeler {
	return &timeseries.Modeler{
		WireUnmarshalerReader: UnmarshalTimeseriesReader,
		WireMarshaler:         MarshalTimeseries,
		WireMarshalWriter:     MarshalTimeseriesWriter,
		WireUnmarshaler:       UnmarshalTimeseries,
		CacheMarshaler:        dataset.MarshalDataSet,
		CacheUnmarshaler:      dataset.UnmarshalDataSet,
	}
}

// UnmarshalTimeseries converts a JSON blob into a Timeseries
func UnmarshalTimeseries(data []byte, trq *timeseries.TimeRangeQuery) (timeseries.Timeseries, error) {
	buf := bytes.NewReader(data)
	return UnmarshalTimeseriesReader(buf, trq)
}

// UnmarshalTimeseriesReader converts a JSON blob into a Timeseries via io.Reader
func UnmarshalTimeseriesReader(reader io.Reader, trq *timeseries.TimeRangeQuery) (timeseries.Timeseries, error) {
	if trq == nil {
		return nil, timeseries.ErrNoTimerangeQuery
	}
	wfd := &WFDocument{}
	d := json.NewDecoder(reader)
	err := d.Decode(wfd)
	if err != nil {
		return nil, err
	}
	ds := &dataset.DataSet{
		Status:         wfd.Status,
		Results:        []*dataset.Result{{}},
		TimeRangeQuery: trq,
		ExtentList:     timeseries.ExtentList{trq.Extent},
	}
	ds.Results[0].SeriesList = make([]*dataset.Series, len(wfd.Data.Results))

	for i, pr := range wfd.Data.Results {
		sh := dataset.SeriesHeader{
			Tags:           pr.Metric,
			QueryStatement: trq.Statement,
		}
		if n, ok := pr.Metric["__name__"]; ok {
			sh.Name = n
		}
		fd := timeseries.FieldDefinition{
			Name:     "value",
			DataType: timeseries.String,
		}
		sh.FieldsList = []timeseries.FieldDefinition{fd}
		var pts dataset.Points
		l := len(pr.Values)
		var ps int64 = 16
		if wfd.Data.ResultType == "matrix" && l > 0 {
			pts = make(dataset.Points, 0, l)
			var wg sync.WaitGroup
			var mtx sync.Mutex
			for _, v := range pr.Values {
				wg.Add(1)
				go func(vals []interface{}) {
					pt, _ := pointFromValues(vals)
					if pt.Epoch > 0 {
						mtx.Lock()
						ps += int64(pt.Size)
						pts = append(pts, pt)
						mtx.Unlock()
					}
					wg.Done()
				}(v)
			}
			wg.Wait()
		} else if wfd.Data.ResultType == "vector" && len(pr.Value) == 2 {
			pts = make(dataset.Points, 1)
			pt, _ := pointFromValues(pr.Value)
			ps = int64(pt.Size)
			pts[0] = pt
			t := time.Unix(0, int64(pt.Epoch))
			ds.ExtentList = timeseries.ExtentList{timeseries.Extent{Start: t, End: t}}
		}
		sh.CalculateSize()
		s := &dataset.Series{
			Header:    sh,
			Points:    pts,
			PointSize: ps,
		}
		ds.Results[0].SeriesList[i] = s
	}
	return ds, nil
}

func pointFromValues(v []interface{}) (dataset.Point, error) {
	if len(v) != 2 {
		return dataset.Point{}, timeseries.ErrInvalidBody
	}
	var f1 float64
	var s string
	var ok bool
	if f1, ok = v[0].(float64); !ok {
		return dataset.Point{}, timeseries.ErrInvalidBody
	}
	if s, ok = v[1].(string); !ok {
		return dataset.Point{}, timeseries.ErrInvalidBody
	}
	return dataset.Point{
		Epoch:  epoch.Epoch(f1) * 1000000000,
		Size:   len(s) + 32, // 8 bytes for epoch, 8 bytes for size, 16 bytes for s stringHeader
		Values: []interface{}{s},
	}, nil
}

// MarshalTimeseries converts a Timeseries into a JSON blob
func MarshalTimeseries(ts timeseries.Timeseries, rlo *timeseries.RequestOptions, status int) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := MarshalTimeseriesWriter(ts, rlo, status, buf)
	return buf.Bytes(), err
}

// MarshalTimeseriesWriter converts a Timeseries into a JSON blob via an io.Writer
func MarshalTimeseriesWriter(ts timeseries.Timeseries, rlo *timeseries.RequestOptions, status int, w io.Writer) error {

	ds, ok := ts.(*dataset.DataSet)
	if !ok {
		return timeseries.ErrUnknownFormat
	}
	// With Prometheus we presume only one Result per Dataset
	if len(ds.Results) != 1 {
		return timeseries.ErrUnknownFormat
	}
	if rw, ok := w.(http.ResponseWriter); ok {
		h := rw.Header()
		h.Set(headers.NameContentType, headers.ValueApplicationJSON+"; charset=UTF-8")
		rw.WriteHeader(status)
	}

	w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[`)) // todo: always "success" ?

	seriesSep := ""
	for _, s := range ds.Results[0].SeriesList {
		if s == nil || len(s.Points) == 0 {
			continue
		}
		w.Write([]byte(seriesSep + `{"metric":{`))
		sep := ""
		for _, k := range s.Header.Tags.Keys() {
			w.Write([]byte(fmt.Sprintf(`%s"%s":"%s"`, sep, k, s.Header.Tags[k])))
			sep = ","
		}
		w.Write([]byte(`},"values":[`))
		sep = ""
		sort.Sort(s.Points)
		for _, p := range s.Points {
			w.Write([]byte(fmt.Sprintf(`%s[%s,"%s"]`,
				sep,
				strconv.FormatFloat(float64(p.Epoch)/1000000000, 'f', -1, 64),
				p.Values[0]),
			))
			sep = ","
		}
		w.Write([]byte("]}"))
		seriesSep = ","
	}
	w.Write([]byte("]}}"))
	return nil
}
