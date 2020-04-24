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

package handlers

import (
	"net/http"
	"sync"

	"github.com/tricksterproxy/trickster/pkg/cache"
	"github.com/tricksterproxy/trickster/pkg/config"
	"github.com/tricksterproxy/trickster/pkg/config/reload"
	"github.com/tricksterproxy/trickster/pkg/proxy/headers"
	tl "github.com/tricksterproxy/trickster/pkg/util/log"
)

// ReloadHandleFunc will reload the running configuration if it has changed
func ReloadHandleFunc(f reload.ReloaderFunc, conf *config.Config, wg *sync.WaitGroup,
	log *tl.Logger, caches map[string]cache.Cache, args []string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if conf.IsStale() {
			log.Warn("configuration reload starting now", tl.Pairs{"source": "reloadEndpoint"})
			f(conf, wg, log, caches, args, false)
			w.Header().Set(headers.NameContentType, headers.ValueTextPlain)
			w.Header().Set(headers.NameCacheControl, headers.ValueNoCache)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("configuration reloaded"))
			return
		}
		w.Header().Set(headers.NameContentType, headers.ValueTextPlain)
		w.Header().Set(headers.NameCacheControl, headers.ValueNoCache)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("configuration NOT reloaded"))
	}
}
