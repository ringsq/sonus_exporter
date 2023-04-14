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

package prober

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	golog "github.com/ringsq/go-logger"
	"github.com/ringsq/sonus_exporter/config"
	"github.com/ringsq/sonus_exporter/sonus"
	"golang.org/x/sync/errgroup"
)

var (
	user     = os.Getenv("SONUS_USER")
	password = os.Getenv("SONUS_PASSWORD")
)

var (
	Probers = []ProbeFn{
		sonus.ServerInfoMetrics,
		sonus.SIPMetrics,
		sonus.CallMetrics,
		sonus.FanMetrics,
		sonus.PowerMetrics,
		sonus.DSPMetrics,
		sonus.TGMetrics,
		sonus.ZoneMetrics,
	}
)

func Handler(w http.ResponseWriter, r *http.Request, c *config.Config, logger log.Logger,
	rh *ResultHistory, timeoutOffset float64,
	params url.Values) {

	if params == nil {
		params = r.URL.Query()
	}
	// moduleName := params.Get("module")
	// if moduleName == "" {
	// 	moduleName = "http_2xx"
	// }

	success := false
	timeoutSeconds, err := getTimeout(r, timeoutOffset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse timeout from Prometheus header: %s", err), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutSeconds*float64(time.Second)))
	defer cancel()
	r = r.WithContext(ctx)

	probeSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "probe_success",
		Help: "Displays whether or not the probe was a success",
	})
	probeDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "probe_duration_seconds",
		Help: "Returns how long the probe took to complete in seconds",
	})

	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	start := time.Now()
	sbc := sonus.NewSBC(target, user, password)

	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)
	registry.MustRegister(probeDurationGauge)

	sl := newScrapeLogger(logger, "test", target)
	if sbc != nil {
		golog.Infof("Starting probe of %s", target)
		g := &errgroup.Group{}
		level.Info(sl).Log("msg", "Beginning probe", "probe", "test", "timeout_seconds", timeoutSeconds)
		for _, probe := range Probers {
			probe := probe
			g.Go(func() error {
				return probe(ctx, sbc, c, registry, sl)
			})
		}
		if err := g.Wait(); err != nil {
			level.Error(sl).Log("msg", "Probe failed", "err", err)
		} else {
			probeSuccessGauge.Set(1)
			level.Info(sl).Log("msg", "Probe succeeded")
			success = true
		}
		golog.Infof("Probe of %s complete", target)
	}

	duration := time.Since(start).Seconds()
	probeDurationGauge.Set(duration)
	level.Info(sl).Log("duration_seconds", duration)

	debugOutput := DebugOutput(&sl.buffer, registry)
	rh.Add("test", target, debugOutput, success)

	if r.URL.Query().Get("debug") == "true" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(debugOutput))
		return
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

type scrapeLogger struct {
	next         log.Logger
	buffer       bytes.Buffer
	bufferLogger log.Logger
}

func newScrapeLogger(logger log.Logger, module string, target string) *scrapeLogger {
	logger = log.With(logger, "module", module, "target", target)
	sl := &scrapeLogger{
		next:   logger,
		buffer: bytes.Buffer{},
	}
	bl := log.NewLogfmtLogger(&sl.buffer)
	sl.bufferLogger = log.With(bl, "ts", log.DefaultTimestampUTC, "caller", log.Caller(6), "module", module, "target", target)
	return sl
}

func (sl scrapeLogger) Log(keyvals ...interface{}) error {
	sl.bufferLogger.Log(keyvals...)
	kvs := make([]interface{}, len(keyvals))
	copy(kvs, keyvals)
	// Switch level to debug for application output.
	for i := 0; i < len(kvs); i += 2 {
		if kvs[i] == level.Key() {
			kvs[i+1] = level.DebugValue()
		}
	}
	return sl.next.Log(kvs...)
}

// DebugOutput returns plaintext debug output for a probe.
func DebugOutput(logBuffer *bytes.Buffer, registry *prometheus.Registry) string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "Logs for the probe:\n")
	logBuffer.WriteTo(buf)
	fmt.Fprintf(buf, "\n\n\nMetrics that would have been returned:\n")
	mfs, err := registry.Gather()
	if err != nil {
		fmt.Fprintf(buf, "Error gathering metrics: %s\n", err)
	}
	for _, mf := range mfs {
		expfmt.MetricFamilyToText(buf, mf)
	}
	// fmt.Fprintf(buf, "\n\n\nModule configuration:\n")
	// c, err := yaml.Marshal(module)
	// if err != nil {
	// 	fmt.Fprintf(buf, "Error marshalling config: %s\n", err)
	// }
	// buf.Write(c)

	return buf.String()
}

func getTimeout(r *http.Request, offset float64) (timeoutSeconds float64, err error) {
	// If a timeout is configured via the Prometheus header, add it to the request.
	if v := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds"); v != "" {
		var err error
		timeoutSeconds, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
	}
	if timeoutSeconds == 0 {
		timeoutSeconds = 120
	}

	var maxTimeoutSeconds = timeoutSeconds - offset
	timeoutSeconds = maxTimeoutSeconds

	return timeoutSeconds, nil
}
