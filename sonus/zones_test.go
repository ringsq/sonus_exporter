package sonus

import (
	"context"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promlog"
	"github.com/ringsq/sonus_exporter/config"
)

func TestZoneProbe(t *testing.T) {
	type args struct {
		ctx      context.Context
		sbc      *SBC
		cfg      *config.Config
		registry *prometheus.Registry
		logger   log.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test run",
			args: args{
				ctx:      context.Background(),
				sbc:      testSBC,
				cfg:      &config.Config{},
				registry: prometheus.NewRegistry(),
				logger:   promlog.New(&promlog.Config{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ZoneProbe(tt.args.ctx, tt.args.sbc, tt.args.cfg, tt.args.registry, tt.args.logger); (err != nil) != tt.wantErr {
				t.Errorf("ZoneProbe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
