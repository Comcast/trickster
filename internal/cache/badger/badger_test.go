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

package badger

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	bo "github.com/Comcast/trickster/internal/cache/badger/options"
	co "github.com/Comcast/trickster/internal/cache/options"
	"github.com/Comcast/trickster/internal/cache/status"
	tl "github.com/Comcast/trickster/internal/util/log"
)

const cacheType = "badger"
const cacheKey = "cacheKey"

func newCacheConfig(t *testing.T) *co.Options {
	dir, err := ioutil.TempDir("/tmp", cacheType)
	if err != nil {
		t.Fatalf("could not create temp directory (%s): %s", dir, err)
	}
	return &co.Options{CacheType: cacheType, Badger: &bo.Options{Directory: dir, ValueDirectory: dir}}
}

func TestConfiguration(t *testing.T) {
	cacheConfig := newCacheConfig(t)
	defer os.RemoveAll(cacheConfig.Badger.Directory)
	bc := Cache{Config: cacheConfig, Logger: tl.ConsoleLogger("error")}

	cfg := bc.Configuration()
	if cfg.CacheType != cacheType {
		t.Fatalf("expected %s got %s", cacheType, cfg.CacheType)
	}
}

func TestBadgerCache_Connect(t *testing.T) {
	cacheConfig := newCacheConfig(t)
	defer os.RemoveAll(cacheConfig.Badger.Directory)
	bc := Cache{Config: cacheConfig, Logger: tl.ConsoleLogger("error")}

	// it should connect
	if err := bc.Connect(); err != nil {
		t.Error(err)
	}
	bc.Close()
}

func TestBadgerCache_ConnectFailed(t *testing.T) {
	cacheConfig := newCacheConfig(t)
	cacheConfig.Badger.Directory = "/root/trickster-test-noaccess"
	os.RemoveAll(cacheConfig.Badger.Directory)
	bc := Cache{Config: cacheConfig, Logger: tl.ConsoleLogger("error")}

	// it should connect
	err := bc.Connect()
	if err == nil {
		t.Errorf("expected file access error for %s", cacheConfig.Badger.Directory)
		bc.Close()
	}
}

func TestBadgerCache_Store(t *testing.T) {
	cacheConfig := newCacheConfig(t)
	defer os.RemoveAll(cacheConfig.Badger.Directory)
	bc := Cache{Config: cacheConfig, Logger: tl.ConsoleLogger("error")}

	if err := bc.Connect(); err != nil {
		t.Error(err)
	}
	defer bc.Close()

	// it should store a value
	err := bc.Store(cacheKey, []byte("data"), time.Duration(60)*time.Second)
	if err != nil {
		t.Error(err)
	}
}

func TestBadgerCache_Remove(t *testing.T) {
	cacheConfig := newCacheConfig(t)
	defer os.RemoveAll(cacheConfig.Badger.Directory)
	bc := Cache{Config: cacheConfig, Logger: tl.ConsoleLogger("error")}

	if err := bc.Connect(); err != nil {
		t.Error(err)
	}
	defer bc.Close()

	// it should store a value
	err := bc.Store(cacheKey, []byte("data"), time.Duration(60)*time.Second)
	if err != nil {
		t.Error(err)
	}

	// it should retrieve a value
	data, ls, err := bc.Retrieve(cacheKey, false)
	if err != nil {
		t.Error(err)
	}
	if ls != status.LookupStatusHit {
		t.Errorf("expected %s got %s", status.LookupStatusHit, ls)
	}
	if string(data) != "data" {
		t.Errorf("wanted \"%s\". got \"%s\".", "data", data)
	}

	bc.Remove(cacheKey)

	// it should be a cache miss
	_, ls, err = bc.Retrieve(cacheKey, false)
	if err == nil {
		t.Errorf("expected key not found error for %s", cacheKey)
	}
	if ls != status.LookupStatusKeyMiss {
		t.Errorf("expected %s got %s", status.LookupStatusKeyMiss, ls)
	}
}

func TestBadgerCache_BulkRemove(t *testing.T) {
	cacheConfig := newCacheConfig(t)
	defer os.RemoveAll(cacheConfig.Badger.Directory)
	bc := Cache{Config: cacheConfig, Logger: tl.ConsoleLogger("error")}

	if err := bc.Connect(); err != nil {
		t.Error(err)
	}
	defer bc.Close()

	// it should store a value
	err := bc.Store(cacheKey, []byte("data"), time.Duration(60)*time.Second)
	if err != nil {
		t.Error(err)
	}

	// it should retrieve a value
	data, ls, err := bc.Retrieve(cacheKey, false)
	if err != nil {
		t.Error(err)
	}
	if string(data) != "data" {
		t.Errorf("wanted \"%s\". got \"%s\".", "data", data)
	}
	if ls != status.LookupStatusHit {
		t.Errorf("expected %s got %s", status.LookupStatusHit, ls)
	}

	bc.BulkRemove([]string{""}, true)
	bc.BulkRemove([]string{cacheKey}, true)

	// it should be a cache miss
	_, ls, err = bc.Retrieve(cacheKey, false)
	if err == nil {
		t.Errorf("expected key not found error for %s", cacheKey)
	}

	if ls != status.LookupStatusKeyMiss {
		t.Errorf("expected %s got %s", status.LookupStatusKeyMiss, ls)
	}

}

func TestBadgerCache_Retrieve(t *testing.T) {
	cacheConfig := newCacheConfig(t)
	defer os.RemoveAll(cacheConfig.Badger.Directory)
	bc := Cache{Config: cacheConfig, Logger: tl.ConsoleLogger("error")}

	if err := bc.Connect(); err != nil {
		t.Error(err)
	}
	defer bc.Close()

	// it should be a cache miss
	_, ls, err := bc.Retrieve(cacheKey, false)
	if err == nil {
		t.Errorf("expected key not found error for %s", cacheKey)
	}
	if ls != status.LookupStatusKeyMiss {
		t.Errorf("expected %s got %s", status.LookupStatusKeyMiss, ls)
	}

	err = bc.Store(cacheKey, []byte("data"), time.Duration(5)*time.Second)
	if err != nil {
		t.Error(err)
	}

	// it should retrieve a value
	data, ls, err := bc.Retrieve(cacheKey, false)
	if err != nil {
		t.Error(err)
	}
	if ls != status.LookupStatusHit {
		t.Errorf("expected %s got %s", status.LookupStatusHit, ls)
	}
	if string(data) != "data" {
		t.Errorf("wanted \"%s\". got \"%s\".", "data", data)
	}

	exp1, err := bc.getExpires(cacheKey)
	if err != nil {
		t.Error(err)
	}

	// the TTL is currently 1 second. update it to 1 hour then wait more than
	// 1 second, to ensure it remained in cache with the correct value
	bc.SetTTL(cacheKey, time.Duration(3600)*time.Second)

	exp2, err := bc.getExpires(cacheKey)
	if err != nil {
		t.Error(err)
	}

	// should be around 3595
	diff := exp2 - exp1
	const expected = 3590

	if diff < 3590 {
		t.Errorf("expected diff >= %d, got %d from: %d - %d", expected, diff, exp2, exp1)
	}

	// try a non-existent cacheKey
	ck2 := cacheKey + "xxxx"
	bc.SetTTL(ck2, time.Duration(3600)*time.Second)

	// it should be a cache miss
	_, ls, err = bc.Retrieve(ck2, false)
	if err == nil {
		t.Errorf("expected key not found error for %s", ck2)
	}
	if ls != status.LookupStatusKeyMiss {
		t.Errorf("expected %s got %s", status.LookupStatusKeyMiss, ls)
	}

	// it should also not have an expires
	_, err = bc.getExpires(ck2)
	if err == nil {
		t.Errorf("expected key not found error for %s", ck2)
	}
}

func TestBadgerCache_Close(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", cacheType)
	if err != nil {
		t.Fatalf("could not create temp directory (%s): %s", dir, err)
	}
	defer os.RemoveAll(dir)

	cacheConfig := &co.Options{CacheType: cacheType, Badger: &bo.Options{Directory: dir, ValueDirectory: dir}}
	bc := Cache{Config: cacheConfig, Logger: tl.ConsoleLogger("error")}

	if err := bc.Connect(); err != nil {
		t.Error(err)
	}

	// it should close
	if err := bc.Close(); err != nil {
		t.Error(err)
	}
}
