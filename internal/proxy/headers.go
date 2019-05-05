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

package proxy

import (
	"net/http"
	"strings"

	"github.com/Comcast/trickster/internal/config"
)

const (
	// Common HTTP Header Values
	hvNoCache         = "no-cache"
	hvApplicationJSON = "application/json"

	// Common HTTP Header Names
	hnCacheControl    = "Cache-Control"
	hnAllowOrigin     = "Access-Control-Allow-Origin"
	hnContentType     = "Content-Type"
	hnContentEncoding = "Content-Encoding"
	hnContentLength   = "Content-Length"
	hnAuthorization   = "Authorization"
	hnXAccelerator    = "X-Accelerator"
	hnXForwardedBy    = "X-Forwarded-By"
	hnXForwardedFor   = "X-Forwarded-For"
	hnAcceptEncoding  = "Accept-Encoding"
)

func copyHeaders(h http.Header) http.Header {
	headers := make(http.Header)
	for k, v := range h {
		headers[k] = make([]string, len(v))
		copy(headers[k], v)
	}
	return headers
}

func addProxyHeaders(remoteAddr string, headers http.Header) {
	if remoteAddr != "" {
		headers.Add(hnXForwardedFor, remoteAddr)
		headers.Add(hnXForwardedBy, config.ApplicationName+" "+config.ApplicationVersion)
	}
}

func addResponseHeaders(headers http.Header) {
	// We're read only and a harmless API, so allow all CORS
	headers.Set(hnAllowOrigin, "*")
	headers.Set(hnXAccelerator, config.ApplicationName+" "+config.ApplicationVersion)
}

func extractHeader(headers http.Header, header string) (string, bool) {
	if hv, ok := headers[header]; ok {
		return strings.Join(hv, "; "), true
	}
	return "", false
}

func removeClientHeaders(headers http.Header) {
	headers.Del(hnAcceptEncoding)
}
