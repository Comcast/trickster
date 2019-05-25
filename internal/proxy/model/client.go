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

package model

import (
	"net/url"

	"github.com/Comcast/trickster/internal/config"
	"github.com/Comcast/trickster/internal/timeseries"
)

// Client is the primary interface for interoperating with Trickster and upstream TSDB's
type Client interface {

	/*
		unused by engines
			RegisterRoutes(string, *config.OriginConfig) // *
			BaseURL() *url.URL // *
			Cache() cache.Cache
			BuildUpstreamURL(*http.Request) *url.URL
			HealthHandler(http.ResponseWriter, *http.Request)

		// Used by engines
			client.Configuration()
			client.ParseTimeRangeQuery(r)
			client.DeriveCacheKey(r, r.Headers.Get(headers.NameAuthorization))
			client.UnmarshalTimeseries(doc.Body)
			client.Name()
			client.FastForwardURL(r)
			client.SetExtent(rq, e)
			client.UnmarshalTimeseries(body)
			client.UnmarshalInstantaneous(body)
			client.MarshalTimeseries(rts)
	*/

	// RegisterRoutes provides a method to register upstream routes to HTTP Handlers
	RegisterRoutes(string, *config.OriginConfig)
	// ParseTimeRangeQuery returns a timeseries.TimeRangeQuery based on the provided HTTP Request
	ParseTimeRangeQuery(*Request) (*timeseries.TimeRangeQuery, error)
	// Configuration returns the configuration for the Proxy Client
	Configuration() *config.OriginConfig
	// Name returns the name of the origin the Proxy Client is handling
	Name() string
	// // BaseURL returns the base URL (schema://host:port/path_prefix) for accessing an upstream origin
	// BaseURL() *url.URL // *
	// // Cache returns the Cache object the Client uses in for Proxy Caching
	// Cache() cache.Cache
	// DeriveCacheKey returns a hashed key for the request, used for request synchronization and cache deconfliction
	DeriveCacheKey(*Request, string) string
	// // BuildUpstreamURL returns an URL for an upstream origin request based on the request URL
	// BuildUpstreamURL(*http.Request) *url.URL
	// FastForwardURL returns the URL to the origin to collect Fast Forward data points based on the provided HTTP Request
	FastForwardURL(*Request) (*url.URL, error)
	// SetExtent will update an upstream request's timerange parameters based on the provided timeseries.Extent
	SetExtent(*Request, *timeseries.Extent)
	// // HealthHandler is an HTTP Handler that checks the health of the upstream origin
	// HealthHandler(http.ResponseWriter, *http.Request)
	// UnmarshalTimeseries will return a Timeseries from the provided byte slice
	UnmarshalTimeseries([]byte) (timeseries.Timeseries, error)
	// MarshalTimeseries will return a byte slice from  the provided Timeseries
	MarshalTimeseries(timeseries.Timeseries) ([]byte, error)
	// UnmarshalInstantaneous will return an Instantaneous Timeseries (only one value instead of a series) from the provided byte slice
	UnmarshalInstantaneous([]byte) (timeseries.Timeseries, error)
}
