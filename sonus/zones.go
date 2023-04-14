package sonus

import (
	"context"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
	"golang.org/x/sync/errgroup"
)

type zonesStats struct {
	Status []zoneStatus `xml:"zoneStatus"`
}
type zoneStatus struct {
	Name                 string  `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 name"`
	TotalCallsAvailable  float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 totalCallsAvailable"`
	InboundCallsUsage    float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 inboundCallsUsage"`
	OutboundCallsUsage   float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 outboundCallsUsage"`
	TotalCallsConfigured float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 totalCallsConfigured"`
	ActiveSipRegCount    float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-ZONE/1.0 activeSipRegCount"`
}

func ZoneMetrics(ctx context.Context, sbc *SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error {
	var (
		Zone_Total_Calls_Configured = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("sonus", "zone", "total_calls_configured"),
			Help: "Total call limit per zone",
		}, []string{"system", "addresscontext", "zone"})

		Zone_Usage_Total = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("sonus", "zone", "usage_total"),
			Help: "Total call limit per zone",
		}, []string{"system", "direction", "addresscontext", "zone"})
	)

	registry.MustRegister(Zone_Total_Calls_Configured)
	registry.MustRegister(Zone_Usage_Total)

	g := &errgroup.Group{}

	for _, aCtx := range sbc.AddressContexts.AddressContext {
		g.Go(func() error {
			aCtx := aCtx
			stats := &zonesStats{}
			err := sbc.GetAndParse(ctx, stats, zoneStatusPath, aCtx.Name)
			if err != nil {
				return err
			}
			for _, stat := range stats.Status {
				Zone_Total_Calls_Configured.WithLabelValues(sbc.System, aCtx.Name, stat.Name).Set(stat.TotalCallsConfigured)
				Zone_Usage_Total.WithLabelValues(sbc.System, "inbound", aCtx.Name, stat.Name).Set(stat.InboundCallsUsage)
				Zone_Usage_Total.WithLabelValues(sbc.System, "outbound", aCtx.Name, stat.Name).Set(stat.OutboundCallsUsage)
			}
			return nil
		})
	}
	return g.Wait()
}
