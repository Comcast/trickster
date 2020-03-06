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

package ruler

import (
	"net/http"

	"github.com/Comcast/trickster/internal/proxy/urls"
)

type inputType string
type extractionFunc func(*http.Request, string) string

var sourceExtractionFuncs = map[inputType]extractionFunc{
	"method":        extractMethodFromSource,
	"url":           extractURLFromSource,
	"url_no_params": extractURLNoParamsFromSource,
	"scheme":        extractSchemeFromSource,
	"host":          extracHostFromSource,
	"hostname":      extractHostnameFromSource,
	"port":          extractPortFromSource,
	"path":          extractPathFromSource,
	"params":        extractParamsFromSource,
	"param":         extractParamFromSource,
	"header":        extractHeaderFromSource,
}

// IsValidSourceName returns true ony if the provided source name is supported by the Rules engine
func IsValidSourceName(source string) bool {
	_, ok := sourceExtractionFuncs[inputType(source)]
	return ok
}

func extractMethodFromSource(r *http.Request, unused string) string {
	if r != nil {
		return r.Method
	}
	return ""
}

func extractURLFromSource(r *http.Request, unused string) string {
	if r != nil && r.URL != nil {
		return r.URL.String()
	}
	return ""
}

func extractURLNoParamsFromSource(r *http.Request, unused string) string {
	if r != nil && r.URL != nil {
		u := urls.Clone(r.URL)
		u.RawQuery = ""
		return u.String()
	}
	return ""
}

func extractSchemeFromSource(r *http.Request, unused string) string {
	if r != nil && r.URL != nil {
		return r.URL.Scheme
	}
	return ""
}

func extracHostFromSource(r *http.Request, unused string) string {
	if r != nil && r.URL != nil {
		return r.URL.Host
	}
	return ""
}

func extractHostnameFromSource(r *http.Request, unused string) string {
	if r != nil && r.URL != nil {
		return r.URL.Hostname()
	}
	return ""
}

func extractPortFromSource(r *http.Request, unused string) string {
	if r != nil && r.URL != nil {
		return r.URL.Port()
	}
	return ""
}

func extractPathFromSource(r *http.Request, unused string) string {
	if r != nil && r.URL != nil {
		return r.URL.Path
	}
	return ""
}

func extractParamsFromSource(r *http.Request, unused string) string {
	if r != nil && r.URL != nil {
		return r.URL.RawQuery
	}
	return ""
}

func extractParamFromSource(r *http.Request, paramName string) string {
	if r != nil && r.URL != nil && r.URL.Query() != nil {
		return r.URL.Query().Get(paramName)
	}
	return ""
}

func extractHeaderFromSource(r *http.Request, headerName string) string {
	if r != nil && r.Header != nil {
		return r.Header.Get(headerName)
	}
	return ""
}
