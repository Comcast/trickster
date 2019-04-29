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

package redis

import (
	"time"

	"github.com/Comcast/trickster/internal/cache"
	"github.com/Comcast/trickster/internal/util/log"
	"github.com/go-redis/redis"
)

// Connect connects to the configured Redis endpoint
func (c *Cache) clientConnect() error {

	log.Info("connecting to redis", log.Pairs{"protocol": c.Config.Redis.Protocol, "Endpoint": c.Config.Redis.Endpoint})
	opts, err := c.clientOpts()
	if err != nil {
		return err
	}
	c.client = redis.NewClient(opts)
	return c.client.Ping().Err()
}

// Store places the the data into the Redis Cache using the provided Key and TTL
func (c *Cache) clientStore(cacheKey string, data []byte, ttl time.Duration) error {
	cache.ObserveCacheOperation(c.Name, c.Config.Type, "set", "none", float64(len(data)))
	log.Debug("redis cache store", log.Pairs{"key": cacheKey})
	return c.client.Set(cacheKey, data, ttl).Err()
}

// Retrieve gets data from the Redis Cache using the provided Key
func (c *Cache) clientRetrieve(cacheKey string) ([]byte, error) {
	res, err := c.client.Get(cacheKey).Result()
	if err != nil {
		log.Debug("redis cache miss", log.Pairs{"key": cacheKey})
		cache.ObserveCacheMiss(cacheKey, c.Name, c.Config.Type)
		return []byte{}, err
	}
	data := []byte(res)
	log.Debug("redis cache retrieve", log.Pairs{"key": cacheKey})
	cache.ObserveCacheOperation(c.Name, c.Config.Type, "get", "hit", float64(len(data)))
	return data, nil
}

// Remove removes an object in cache, if present
func (c *Cache) clientRemove(cacheKey string) {
	log.Debug("redis cache remove", log.Pairs{"key": cacheKey})
	c.client.Del(cacheKey)
}

// BulkRemove removes a list of objects from the cache. noLock is not used for Redis
func (c *Cache) clientBulkRemove(cacheKeys []string, noLock bool) {
	log.Debug("redis cache bulk remove", log.Pairs{})
	c.client.Del(cacheKeys...)
}

// Close disconnects from the Redis Cache
func (c *Cache) clientClose() error {
	log.Info("closing redis connection", log.Pairs{})
	c.client.Close()
	return nil
}

func (c *Cache) clientOpts() (*redis.Options, error) {

	o := &redis.Options{
		Addr: c.Config.Redis.Endpoint,
	}

	if c.Config.Redis.Protocol != "" {
		o.Network = c.Config.Redis.Protocol
	}

	if c.Config.Redis.Password != "" {
		o.Password = c.Config.Redis.Password
	}

	if c.Config.Redis.DB != 0 {
		o.DB = c.Config.Redis.DB
	}

	if c.Config.Redis.MaxRetries != 0 {
		o.MaxRetries = c.Config.Redis.MaxRetries
	}

	if c.Config.Redis.MinRetryBackoffMS != 0 {
		o.MinRetryBackoff = durationFromMS(c.Config.Redis.MinRetryBackoffMS)
	}

	if c.Config.Redis.MaxRetryBackoffMS != 0 {
		o.MaxRetryBackoff = durationFromMS(c.Config.Redis.MaxRetryBackoffMS)
	}

	if c.Config.Redis.DialTimeoutMS != 0 {
		o.DialTimeout = durationFromMS(c.Config.Redis.DialTimeoutMS)
	}

	if c.Config.Redis.ReadTimeoutMS != 0 {
		o.ReadTimeout = durationFromMS(c.Config.Redis.ReadTimeoutMS)
	}

	if c.Config.Redis.WriteTimeoutMS != 0 {
		o.WriteTimeout = durationFromMS(c.Config.Redis.WriteTimeoutMS)
	}

	if c.Config.Redis.PoolSize != 0 {
		o.PoolSize = c.Config.Redis.PoolSize
	}

	if c.Config.Redis.MinIdleConns != 0 {
		o.MinIdleConns = c.Config.Redis.MinIdleConns
	}

	if c.Config.Redis.MaxConnAgeMS != 0 {
		o.MaxConnAge = durationFromMS(c.Config.Redis.MaxConnAgeMS)
	}

	if c.Config.Redis.PoolTimeoutMS != 0 {
		o.PoolTimeout = durationFromMS(c.Config.Redis.PoolTimeoutMS)
	}

	if c.Config.Redis.IdleTimeoutMS != 0 {
		o.IdleTimeout = durationFromMS(c.Config.Redis.IdleTimeoutMS)
	}

	if c.Config.Redis.IdleCheckFrequencyMS != 0 {
		o.IdleCheckFrequency = durationFromMS(c.Config.Redis.IdleCheckFrequencyMS)
	}

	return o, nil
}
