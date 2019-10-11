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

package cache

import "strconv"

// LookupStatus defines the possible status of a cache lookup
type LookupStatus int

const (
	// LookupStatusHit indicates a full cache hit on lookup
	LookupStatusHit = LookupStatus(iota)
	// LookupStatusPartialHit indicates a partial cache hit (key exists and has some data
	// for requested time range, but not all) on lookup
	LookupStatusPartialHit
	// LookupStatusRangeMiss indicates a range miss (key exists but no data for requested time range) on lookup
	LookupStatusRangeMiss
	// LookupStatusKeyMiss indicates a full key miss (cache key does not exist) on lookup
	LookupStatusKeyMiss
	// LookupStatusPurge indicates the cache key, if it existed, was purged as directed
	// in upstream response or down stream request http headers
	LookupStatusPurge
	// LookupStatusProxyError indicates that a proxy error occurred retrieving a cacheable dataset
	// in upstream response or down stream request http headers
	LookupStatusProxyError
	// LookupStatusProxyOnly indicates that the request was fully proxied to the origin without using the cache
	LookupStatusProxyOnly
)

var cacheLookupStatusNames = map[string]LookupStatus{
	"hit":         LookupStatusHit,
	"phit":        LookupStatusPartialHit,
	"rmiss":       LookupStatusRangeMiss,
	"kmiss":       LookupStatusKeyMiss,
	"purge":       LookupStatusPurge,
	"proxy-error": LookupStatusProxyError,
	"proxy-only":  LookupStatusProxyOnly,
}

var cacheLookupStatusValues = map[LookupStatus]string{
	LookupStatusHit:        "hit",
	LookupStatusPartialHit: "phit",
	LookupStatusRangeMiss:  "rmiss",
	LookupStatusKeyMiss:    "kmiss",
	LookupStatusPurge:      "purge",
	LookupStatusProxyError: "proxy-error",
	LookupStatusProxyOnly:  "proxy-only",
}

func (s LookupStatus) String() string {
	if v, ok := cacheLookupStatusValues[s]; ok {
		return v
	}
	return strconv.Itoa(int(s))
}