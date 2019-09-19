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

	"github.com/Comcast/trickster/internal/config"
	"github.com/Comcast/trickster/internal/proxy/headers"
)

func (c *Client) registerHandlers() {
	c.handlersRegistered = true
	// This is the registry of handlers that Trickster supports for Prometheus,
	// and are able to be referenced by name (map key) in Config Files
	c.handlers["health"] = http.HandlerFunc(c.HealthHandler)
	c.handlers[mnQuery] = http.HandlerFunc(c.QueryHandler)
	c.handlers["proxy"] = http.HandlerFunc(c.ProxyHandler)
}

// Handlers returns a map of the HTTP Handlers the client has registered
func (c *Client) Handlers() map[string]http.Handler {
	if !c.handlersRegistered {
		c.registerHandlers()
	}
	return c.handlers
}

// DefaultPathConfigs returns the default PathConfigs for the given OriginType
func (c *Client) DefaultPathConfigs() (map[string]*config.PathConfig, []string) {

	paths := map[string]*config.PathConfig{

		"/" + mnQuery: &config.PathConfig{
			Path:            "/" + mnQuery,
			HandlerName:     mnQuery,
			Methods:         []string{http.MethodGet, http.MethodPost},
			CacheKeyParams:  []string{upDB, upQuery, "u", "p"},
			CacheKeyHeaders: []string{headers.NameAuthorization},
			DefaultTTLSecs:  c.cache.Configuration().ObjectTTLSecs,
			DefaultTTL:      c.cache.Configuration().ObjectTTL,
		},

		"/": &config.PathConfig{
			Path:        "/",
			HandlerName: "proxy",
			Methods:     []string{http.MethodGet, http.MethodPost},
		},
	}

	orderedPaths := []string{"/" + mnQuery, "/"}

	return paths, orderedPaths

}
