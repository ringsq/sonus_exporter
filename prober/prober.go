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
	"context"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
	"github.com/ringsq/sonus_exporter/sonus"
)

// A ProbeFn calls the SBC and adds its metrics to the registry
type ProbeFn func(ctx context.Context, sbc *sonus.SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error
