package sonus

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestBuildMetrics(t *testing.T) {
	type args struct {
		registry *prometheus.Registry
		t        reflect.Type
	}
	tests := []struct {
		name string
		args args
		want map[string]*prometheus.GaugeVec
	}{
		{
			name: "Test metrics",
			args: args{registry: prometheus.NewRegistry(), t: reflect.TypeOf(ZoneStats{})},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildMetrics(tt.args.registry, tt.args.t); !reflect.DeepEqual(got, tt.want) {

			}
		})
	}
}
