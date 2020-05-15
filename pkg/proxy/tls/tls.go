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

// Package tls handles options for TLS (https) requests
package tls

import "github.com/tricksterproxy/trickster/pkg/config"

// OptionsChanged will return true if the TLS options for any origin
// is different between configs
func OptionsChanged(conf, oldConf *config.Config) bool {

	if conf == nil {
		return false
	}
	if oldConf == nil {
		return true
	}

	for k, v := range oldConf.Origins {
		if v.TLS != nil && v.TLS.ServeTLS {
			if o, ok := conf.Origins[k]; !ok ||
				o.TLS == nil || !o.TLS.ServeTLS ||
				!o.TLS.Equal(v.TLS) {
				return true
			}
		}
	}

	for k, v := range conf.Origins {
		if v.TLS != nil && v.TLS.ServeTLS {
			if o, ok := oldConf.Origins[k]; !ok ||
				o.TLS == nil || !o.TLS.ServeTLS ||
				!o.TLS.Equal(v.TLS) {
				return true
			}
		}
	}

	return false
}
