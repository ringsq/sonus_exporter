// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	configReloadSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "blackbox_exporter",
		Name:      "config_last_reload_successful",
		Help:      "Blackbox exporter config loaded successfully.",
	})

	configReloadSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "blackbox_exporter",
		Name:      "config_last_reload_success_timestamp_seconds",
		Help:      "Timestamp of the last successful configuration reload.",
	})
)

func init() {
	prometheus.MustRegister(configReloadSuccess)
	prometheus.MustRegister(configReloadSeconds)
}

type Config struct {
}

type SafeConfig struct {
	sync.RWMutex
	C *Config
}

func (sc *SafeConfig) ReloadConfig(confFile string, logger log.Logger) (err error) {
	var c = &Config{}
	defer func() {
		if err != nil {
			configReloadSuccess.Set(0)
		} else {
			configReloadSuccess.Set(1)
			configReloadSeconds.SetToCurrentTime()
		}
	}()

	// load the config here

	sc.Lock()
	sc.C = c
	sc.Unlock()

	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}
	return nil
}

// isCompressionAcceptEncodingValid validates the compression +
// Accept-Encoding combination.
//
// If there's a compression setting, and there's also an accept-encoding
// header, they MUST match, otherwise we end up requesting something
// that doesn't include the specified compression, and that's likely to
// fail, depending on how the server is configured. Testing that the
// server _ignores_ Accept-Encoding, e.g. by not including a particular
// compression in the header but expecting it in the response falls out
// of the scope of the tests we perform.
//
// With that logic, this function validates that if a compression
// algorithm is specified, it's covered by the specified accept encoding
// header. It doesn't need to be the most prefered encoding, but it MUST
// be included in the prefered encodings.
func isCompressionAcceptEncodingValid(encoding, acceptEncoding string) bool {
	// unspecified compression + any encoding value is valid
	// any compression + no accept encoding is valid
	if encoding == "" || acceptEncoding == "" {
		return true
	}

	type encodingQuality struct {
		encoding string
		quality  float32
	}

	var encodings []encodingQuality

	for _, parts := range strings.Split(acceptEncoding, ",") {
		var e encodingQuality

		if idx := strings.LastIndexByte(parts, ';'); idx == -1 {
			e.encoding = strings.TrimSpace(parts)
			e.quality = 1.0
		} else {
			parseQuality := func(str string) float32 {
				q, err := strconv.ParseFloat(str, 32)
				if err != nil {
					return 0
				}
				return float32(math.Round(q*1000) / 1000)
			}

			e.encoding = strings.TrimSpace(parts[:idx])

			q := strings.TrimSpace(parts[idx+1:])
			q = strings.TrimPrefix(q, "q=")
			e.quality = parseQuality(q)
		}

		encodings = append(encodings, e)
	}

	sort.SliceStable(encodings, func(i, j int) bool {
		return encodings[j].quality < encodings[i].quality
	})

	for _, e := range encodings {
		if encoding == e.encoding || e.encoding == "*" {
			return e.quality > 0
		}
	}

	return false
}
