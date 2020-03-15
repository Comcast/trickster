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

package config

import (
	"flag"
	"sync"

	d "github.com/Comcast/trickster/internal/config/defaults"
)

const (
	// Command-line flags
	cfConfig      = "config"
	cfVersion     = "version"
	cfLogLevel    = "log-level"
	cfInstanceID  = "instance-id"
	cfOrigin      = "origin-url"
	cfOriginType  = "origin-type"
	cfProxyPort   = "proxy-port"
	cfMetricsPort = "metrics-port"
	cfReloadPort  = "reload-port"
)

var parseOnce sync.Once

// Flags holds the values for whitelisted flags
type Flags struct {
	PrintVersion      bool
	ConfigPath        string
	customPath        bool
	Origin            string
	OriginType        string
	ProxyListenPort   int
	MetricsListenPort int
	ReloadListenPort  int
	LogLevel          string
	InstanceID        int

	providedOriginURL  string
	providedOriginType string
}

func parseFlags(applicationName string, arguments []string) (*Flags, error) {

	flags := &Flags{}
	flagSet := flag.NewFlagSet("trickster", flag.ContinueOnError)

	flagSet.BoolVar(&flags.PrintVersion, cfVersion, false, "Prints trickster version")
	flagSet.StringVar(&flags.ConfigPath, cfConfig, "", "Path to Trickster Config File")
	flagSet.StringVar(&flags.LogLevel, cfLogLevel, "", "Level of Logging to use (debug, info, warn, error)")
	flagSet.IntVar(&flags.InstanceID, cfInstanceID, 0, "Instance ID is for running multiple Trickster processes from the same config while logging to their own files.")
	flagSet.StringVar(&flags.Origin, cfOrigin, "", "URL to the Origin. Enter it like you would in grafana, e.g., http://prometheus:9090")
	flagSet.StringVar(&flags.OriginType, cfOriginType, "", "Type of origin (prometheus, influxdb)")
	flagSet.IntVar(&flags.ProxyListenPort, cfProxyPort, 0, "Port that the primary Proxy server will listen on.")
	flagSet.IntVar(&flags.MetricsListenPort, cfMetricsPort, 0, "Port that the /metrics endpoint will listen on.")
	//flagSet.IntVar(&flags.ReloadListenPort, cfReloadPort, 0, "Port that the /-/reload endpoint will listen on.")

	var err error
	err = flagSet.Parse(arguments)
	if err != nil {
		return nil, err
	}
	if flags.ConfigPath != "" {
		flags.customPath = true
	} else {
		flags.ConfigPath = d.DefaultConfigPath
	}
	return flags, nil
}

// loadFlags loads configuration from command line flags.
func (c *TricksterConfig) loadFlags(flags *Flags) {
	if len(flags.Origin) > 0 {
		c.providedOriginURL = flags.Origin
	}
	if len(flags.OriginType) > 0 {
		c.providedOriginType = flags.OriginType
	}
	if flags.ProxyListenPort > 0 {
		c.Frontend.ListenPort = flags.ProxyListenPort
	}
	if flags.MetricsListenPort > 0 {
		c.Metrics.ListenPort = flags.MetricsListenPort
	}
	// if flags.ReloadListenPort > 0 {
	// 	c.Main.Reload.ListenPort = flags.ReloadListenPort
	// }
	if flags.LogLevel != "" {
		c.Logging.LogLevel = flags.LogLevel
	}
	if flags.InstanceID > 0 {
		c.Main.InstanceID = flags.InstanceID
	}
}
