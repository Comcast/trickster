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

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Comcast/trickster/internal/config"
	"github.com/Comcast/trickster/internal/util/metrics"
)

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
	// LookupStatusPurge indicates a the cache key, if it existed, was purged as directed
	// in upstream response or down stream request http headers
	LookupStatusPurge
	// LookupStatusProxyError indicates a that a proxy error occurred retreiving a cacheable dataset
	// in upstream response or down stream request http headers
	LookupStatusProxyError
)

func (s LookupStatus) String() string {
	name := []string{"hit", "phit", "rmiss", "kmiss"}
	switch {
	case s <= LookupStatusKeyMiss:
		return name[int(s)]
	default:
		return strconv.Itoa(int(s))
	}
}

// Cache is the interface for the supported caching fabrics
// When making new cache types, Retrieve() must return an error on cache miss
type Cache interface {
	Connect() error
	Store(cacheKey string, data []byte, ttl time.Duration) error
	Retrieve(cacheKey string) ([]byte, error)
	Remove(cacheKey string)
	BulkRemove(cacheKeys []string, noLock bool)
	Close() error
	Configuration() *config.CachingConfig
}

// ObserveCacheMiss returns a standard Cache Miss response
func ObserveCacheMiss(cacheKey, cacheName, cacheType string) ([]byte, error) {
	ObserveCacheOperation(cacheName, cacheType, "get", "miss", 0)
	return nil, fmt.Errorf("value  for key [%s] not in cache", cacheKey)
}

// CacheError returns an empty cache object and the formatted error
func CacheError(cacheKey, cacheName, cacheType string, msg string) ([]byte, error) {
	ObserveCacheEvent(cacheName, cacheType, "error", msg)
	return nil, fmt.Errorf(msg, cacheKey)
}

// ObserveCacheOperation increments counters as cache operations occur
func ObserveCacheOperation(cache, cacheType, operation, status string, bytes float64) {
	metrics.CacheObjectOperations.WithLabelValues(cache, cacheType, operation, status).Inc()
	if bytes > 0 {
		metrics.CacheByteOperations.WithLabelValues(cache, cacheType, operation, status).Add(float64(bytes))
	}
}

// ObserveCacheEvent increments counters as cache events occur
func ObserveCacheEvent(cache, cacheType, event, reason string) {
	metrics.CacheEvents.WithLabelValues(cache, cacheType, event, reason).Inc()
}

// ObserveCacheSizeChange adjust counters and gauges as the cache size changes due to object operations
func ObserveCacheSizeChange(cache, cacheType string, byteCount, objectCount int64) {
	metrics.CacheObjects.WithLabelValues(cache, cacheType).Set(float64(objectCount))
	metrics.CacheBytes.WithLabelValues(cache, cacheType).Set(float64(byteCount))
}
