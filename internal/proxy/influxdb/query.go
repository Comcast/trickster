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

package influxdb

import (
	"net/http"

	"github.com/Comcast/trickster/internal/proxy"
)

func (c Client) QueryHandler(w http.ResponseWriter, r *http.Request) {
	u := c.BuildUpstreamURL(r)
	proxy.ObjectProxyCacheRequest(proxy.NewRequest(c.Name, proxy.OtInfluxDb, "QueryHandler", r.Method, u, r.Header, c.Config.Timeout, r), w, &c, c.Cache, 30, false, false)
}
