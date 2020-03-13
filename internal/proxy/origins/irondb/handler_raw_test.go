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

package irondb

import (
	"io/ioutil"
	"net/http"
	"testing"

	oo "github.com/Comcast/trickster/internal/proxy/origins/options"
	"github.com/Comcast/trickster/internal/proxy/request"
	tu "github.com/Comcast/trickster/internal/util/testing"
)

func TestRawHandler(t *testing.T) {

	client := &Client{name: "test"}
	ts, w, r, hc, err := tu.NewTestInstance("", client.DefaultPathConfigs,
		200, "{}", nil, "irondb", "/raw/00112233-4455-6677-8899-aabbccddeeff/metric"+
			"?start_ts=0&end_ts=900", "debug")
	rsc := request.GetResources(r)
	rsc.OriginClient = client
	client.config = rsc.OriginConfig
	client.webClient = hc
	client.config.HTTPClient = hc
	defer ts.Close()
	if err != nil {
		t.Error(err)
	}

	client.RawHandler(w, r)
	resp := w.Result()

	// it should return 200 OK
	if resp.StatusCode != 200 {
		t.Errorf("expected 200 got %d.", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if string(bodyBytes) != "{}" {
		t.Errorf("expected '{}' got %s.", bodyBytes)
	}
}

func TestRawHandlerParseTimeRangeQuery(t *testing.T) {

	// provide bad URL with no TimeRange query params
	hc := tu.NewTestWebClient()
	cfg := oo.NewOptions()
	client := &Client{name: "test", webClient: hc, config: cfg}

	cfg.Paths = client.DefaultPathConfigs(cfg)
	r, err := http.NewRequest(http.MethodGet, "http://0/raw/00112233-4455-6677-8899-aabbccddeeff/metric", nil)
	if err != nil {
		t.Error(err)
	}
	//tr := model.NewRequest("RawHandler", r.Method, r.URL, r.Header, cfg.Timeout, r, hc)

	// case where everything is good
	r.URL.RawQuery = "start_ts=9012&end_ts=3456"
	trq, err := client.rawHandlerParseTimeRangeQuery(r)
	if err != nil {
		t.Error(err)
	}
	if trq == nil {
		t.Errorf("expected value got nil for %s", r.URL.RawQuery)
	}

	// missing start param
	r.URL.RawQuery = "end_ts=3456"
	_, err = client.rawHandlerParseTimeRangeQuery(r)
	if err == nil {
		t.Errorf("expected error for parameter missing")
	}

	// can't parse start param
	r.URL.RawQuery = "start_ts=abc&end_ts=3456"
	_, err = client.rawHandlerParseTimeRangeQuery(r)
	if err == nil {
		t.Errorf("expected error for parameter missing")
	}

	// missing end param
	r.URL.RawQuery = "start_ts=3456"
	_, err = client.rawHandlerParseTimeRangeQuery(r)
	if err == nil {
		t.Errorf("expected error for parameter missing")
	}

	// can't parse end param
	r.URL.RawQuery = "start_ts=9012&end_ts=def"
	_, err = client.rawHandlerParseTimeRangeQuery(r)
	if err == nil {
		t.Errorf("expected error for parameter missing")
	}

}
