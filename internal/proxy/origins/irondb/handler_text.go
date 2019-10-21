/**
* Copyright 2018 Comcast Cable Communications Management, LLC
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package irondb

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Comcast/trickster/internal/proxy/engines"
	"github.com/Comcast/trickster/internal/proxy/errors"
	"github.com/Comcast/trickster/internal/proxy/model"
	"github.com/Comcast/trickster/internal/timeseries"
	"github.com/Comcast/trickster/internal/util/md5"
)

// TextHandler handles requests for text timeseries data and processes them
// through the delta proxy cache.
func (c *Client) TextHandler(w http.ResponseWriter, r *http.Request) {
	u := c.BuildUpstreamURL(r)
	engines.DeltaProxyCacheRequest(
		model.NewRequest("TextHandler",
			r.Method, u, r.Header, c.config.Timeout, r, c.webClient),
		w, c)
}

// textHandlerSetExtent will change the upstream request query to use the
// provided Extent.
func (c Client) textHandlerSetExtent(r *model.Request,
	extent *timeseries.Extent) {
	ps := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 5)
	if len(ps) < 5 || ps[0] != "read" {
		return
	}

	sb := new(strings.Builder)
	if strings.HasPrefix(r.URL.Path, "/") {
		sb.WriteString("/")
	}

	sb.WriteString("read")
	sb.WriteString("/" + strconv.FormatInt(extent.Start.Unix(), 10))
	sb.WriteString("/" + strconv.FormatInt(extent.End.Unix(), 10))
	sb.WriteString("/" + strings.Join(ps[3:], "/"))
	r.URL.Path = sb.String()
}

// textHandlerParseTimeRangeQuery parses the key parts of a TimeRangeQuery
// from the inbound HTTP Request.
func (c *Client) textHandlerParseTimeRangeQuery(
	r *model.Request) (*timeseries.TimeRangeQuery, error) {
	trq := &timeseries.TimeRangeQuery{}
	ps := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 5)
	if len(ps) < 5 || ps[0] != "read" {
		return nil, errors.NotTimeRangeQuery()
	}

	trq.Statement = "/read/" + strings.Join(ps[3:], "/")

	var err error
	if trq.Extent.Start, err = parseTimestamp(ps[1]); err != nil {
		return nil, err
	}

	if trq.Extent.End, err = parseTimestamp(ps[2]); err != nil {
		return nil, err
	}

	return trq, nil
}

// textHandlerDeriveCacheKey calculates a query-specific keyname based on the
// user request.
func (c Client) textHandlerDeriveCacheKey(path string, params url.Values,
	headers http.Header, body io.ReadCloser, extra string) string {
	var sb strings.Builder
	sb.WriteString(path)
	ps := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 5)
	if len(ps) >= 5 || ps[0] == "read" {
		sb.WriteString("/read/" + strings.Join(ps[3:], "/"))
	}

	sb.WriteString(extra)
	return md5.Checksum(sb.String())
}
