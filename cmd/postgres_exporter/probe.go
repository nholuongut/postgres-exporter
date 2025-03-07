// Copyright 2022 Nho Luong DevOps
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

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/nholuongut/postgres_exporter/collector"
	"github.com/nholuongut/postgres_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func handleProbe(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conf := c.GetConfig()
		params := r.URL.Query()
		target := params.Get("target")
		if target == "" {
			http.Error(w, "target is required", http.StatusBadRequest)
			return
		}
		var authModule config.AuthModule
		authModuleName := params.Get("auth_module")
		if authModuleName == "" {
			level.Info(logger).Log("msg", "no auth_module specified, using default")
		} else {
			var ok bool
			authModule, ok = conf.AuthModules[authModuleName]
			if !ok {
				http.Error(w, fmt.Sprintf("auth_module %s not found", authModuleName), http.StatusBadRequest)
				return
			}
			if authModule.UserPass.Username == "" || authModule.UserPass.Password == "" {
				http.Error(w, fmt.Sprintf("auth_module %s has no username or password", authModuleName), http.StatusBadRequest)
				return
			}
		}

		dsn, err := authModule.ConfigureTarget(target)
		if err != nil {
			level.Error(logger).Log("msg", "failed to configure target", "err", err)
			http.Error(w, fmt.Sprintf("could not configure dsn for target: %v", err), http.StatusBadRequest)
			return
		}

		// TODO(@sysadmind): Timeout

		probeSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_success",
			Help: "Displays whether or not the probe was a success",
		})
		probeDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_duration_seconds",
			Help: "Returns how long the probe took to complete in seconds",
		})

		tl := log.With(logger, "target", target)

		start := time.Now()
		registry := prometheus.NewRegistry()
		registry.MustRegister(probeSuccessGauge)
		registry.MustRegister(probeDurationGauge)

		// Run the probe
		pc, err := collector.NewProbeCollector(tl, registry, dsn)
		if err != nil {
			probeSuccessGauge.Set(0)
			probeDurationGauge.Set(time.Since(start).Seconds())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO(@sysadmind): Remove the registry.MustRegister() call below and instead handle the collection here. That will allow
		// for the passing of context, handling of timeouts, and more control over the collection.
		// The current NewProbeCollector() implementation relies on the MustNewConstMetric() call to create the metrics which is not
		// ideal to use without the registry.MustRegister() call.
		_ = ctx

		registry.MustRegister(pc)

		duration := time.Since(start).Seconds()
		probeDurationGauge.Set(duration)

		// TODO check success, etc
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}
