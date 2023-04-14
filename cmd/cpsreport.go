// cpsreport runs periodically via cron to query Ribbon SBC APIs for CPS configuration,
// and pushes the results to the Prometheus pushgateway for monitoring and analysis.
//
// All servers being queried are expected to use the same username and password,
// specified by SONUS_USER and SONUS_PASSWORD environment variables

package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/ringsq/go-logger"
	"github.com/ringsq/sonus_exporter/sonus"
)

const (
	tgConfigPath = "/config/addressContext/%s/zone/%s/sipTrunkGroup/"
)

var servers = []string{"10.0.204.60", "10.5.204.60"}

const pushgateway = "http://pgw.m5.run:9091"

var (
	scrapeDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sonus",
		Subsystem: "TG",
		Name:      "cps_scrape_duration",
		Help:      "The amount of time it took to scrape the Sonus API",
	}, []string{"instance", "system"})

	callRate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sonus",
		Subsystem: "TG",
		Name:      "cps",
		Help:      "Maximum calls-per-second configured for this trunkgroup",
	}, []string{"instance", "system", "addresscontext", "zone", "name"})
	burstRate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sonus",
		Subsystem: "TG",
		Name:      "cps_burst",
		Help:      "The burst calls-per-second rate configured for this trunkgroup",
	}, []string{"instance", "system", "addresscontext", "zone", "name"})
)

func main() {
	user := os.Getenv("SONUS_USER")
	password := os.Getenv("SONUS_PASSWORD")
	registry := prometheus.NewRegistry()
	registry.MustRegister(scrapeDuration, callRate, burstRate)
	pusher := push.New(pushgateway, "sonus").Gatherer(registry)

	for _, server := range servers {
		start := time.Now()
		sbc := sonus.NewSBC(server, user, password)
		if sbc == nil {
			log.Errorf("Could not establish connection to %s", server)
			continue
		}
		for _, aCtx := range sbc.AddressContexts.AddressContext {
			for _, zone := range aCtx.Zone {
				tglist := &sonus.Trunkgroups{}
				err := sbc.GetAndParse(context.Background(), tglist, tgConfigPath, aCtx.Name, zone.Name)
				if err != nil {
					log.Errorf("Error getting TG configs: %v", err)
				}
				for _, tg := range tglist.SipTrunkGroup {
					if len(tg.IngressIpPrefix) > 0 {
						if strings.HasPrefix(tg.Name, "C_") {
							callRate.WithLabelValues(server, sbc.System, aCtx.Name, zone.Name, tg.Name).Set(tg.Cac.Ingress.CallRateMax)
							burstRate.WithLabelValues(server, sbc.System, aCtx.Name, zone.Name, tg.Name).Set(tg.Cac.Ingress.CallBurstMax)
						}
					}
				}
			}
		}
		duration := float64(time.Since(start).Seconds() * 1e3)
		scrapeDuration.WithLabelValues(server, sbc.System).Set(duration)
	}
	if err := pusher.Push(); err != nil {
		log.Errorf("Failure to push metrics: %v", err)
	}
}
