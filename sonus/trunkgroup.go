package sonus

import (
	"context"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
)

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <globalTrunkGroupStatus xmlns="http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0">
    <name>TEST</name>
    <state>inService</state>
    <totalCallsAvailable>100</totalCallsAvailable>
    <totalCallsInboundReserved>0</totalCallsInboundReserved>
    <inboundCallsUsage>0</inboundCallsUsage>
    <outboundCallsUsage>0</outboundCallsUsage>
    <totalCallsConfigured>100</totalCallsConfigured>
    <priorityCallUsage>0</priorityCallUsage>
    <totalOutboundCallsReserved>0</totalOutboundCallsReserved>
    <bwCurrentLimit>-1</bwCurrentLimit>
    <bwAvailable>-1</bwAvailable>
    <bwInboundUsage>0</bwInboundUsage>
    <bwOutboundUsage>0</bwOutboundUsage>
    <packetOutDetectState>normal</packetOutDetectState>
    <addressContext>default</addressContext>
    <zone>zone_23</zone>
    <priorityBwUsage>0</priorityBwUsage>
  </globalTrunkGroupStatus>
  ...
</collection>
*/

type trunkGroupCollection struct {
	TrunkGroupStatus []*trunkGroupStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 globalTrunkGroupStatus,omitempty"`
}

type trunkGroupStatus struct {
	AddressContext             string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 addressContext"`
	BandwidthAvailable         float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwAvailable"`
	BandwidthCurrentLimit      float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwCurrentLimit"`
	BandwidthInboundUsage      float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwInboundUsage"`
	BandwidthOutboundUsage     float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwOutboundUsage"`
	InboundCallsUsage          float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 inboundCallsUsage"`
	Name                       string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 name"`
	OutboundCallsUsage         float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 outboundCallsUsage"`
	PacketOutDetectState       string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 packetOutDetectState"`
	PriorityBwUsage            float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 priorityBwUsage"`
	PriorityCallUsage          float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 priorityCallUsage"`
	State                      string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 state"`
	TotalCallsAvailable        float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsAvailable"`
	TotalCallsConfigured       float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsConfigured"`
	TotalCallsInboundReserved  float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsInboundReserved"`
	TotalOutboundCallsReserved float64 `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalOutboundCallsReserved"`
	Zone                       string  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 zone"`
}

type Trunkgroups struct {
	SipTrunkGroup []struct {
		Name  string `xml:"name"`
		State string `xml:"state"`
		Mode  string `xml:"mode"`
		Cac   struct {
			CallLimit string `xml:"callLimit"`
			Ingress   struct {
				CallRateMax  float64 `xml:"callRateMax"`
				CallBurstMax float64 `xml:"callBurstMax"`
			} `xml:"ingress"`
		} `xml:"cac"`
		IngressIpPrefix []struct {
			IpAddress    string `xml:"ipAddress"`
			PrefixLength string `xml:"prefixLength"`
		} `xml:"ingressIpPrefix"`
	} `xml:"sipTrunkGroup"`
}

func (t trunkGroupStatus) stateToMetric() float64 {
	switch t.State {
	case "inService":
		return 1
	default:
		return 0
	}
}

func (t trunkGroupStatus) outStateToMetric() float64 {
	switch t.PacketOutDetectState {
	case "normal":
		return 1
	default:
		return 0
	}
}

func TGMetrics(ctx context.Context, sbc *SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error {
	var (
		TG_Bandwidth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("sonus", "TG", "bytes"),
			Help: "Bandwidth in use by current calls",
		}, []string{"system", "zone", "name", "direction"})

		TG_OBState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("sonus", "TG", "outbound_state"),
			Help: "State of outbound calls on the trunkgroup",
		}, []string{"system", "zone", "name"})

		TG_State = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("sonus", "TG", "state"),
			Help: "State of the trunkgroup",
		}, []string{"system", "zone", "name"})

		TG_TotalChans = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("sonus", "TG", "total_channels"),
			Help: "Number of configured channels",
		}, []string{"system", "zone", "name"})

		TG_Usage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("sonus", "TG", "usage_total"),
			Help: "Number of active calls",
		}, []string{"system", "zone", "name", "direction"})

		tgs = new(trunkGroupCollection)
	)

	registry.MustRegister(TG_Bandwidth)
	registry.MustRegister(TG_OBState)
	registry.MustRegister(TG_State)
	registry.MustRegister(TG_TotalChans)
	registry.MustRegister(TG_Usage)

	err := sbc.GetAndParse(ctx, tgs, tgStatusPath)
	if err != nil {
		return err
	}
	for _, tg := range tgs.TrunkGroupStatus {
		TG_Usage.WithLabelValues(sbc.System, tg.Zone, tg.Name, "inbound").Set(tg.InboundCallsUsage)
		TG_Usage.WithLabelValues(sbc.System, tg.Zone, tg.Name, "outbound").Set(tg.OutboundCallsUsage)
		TG_Bandwidth.WithLabelValues(sbc.System, tg.Zone, tg.Name, "inbound").Set(tg.BandwidthInboundUsage)
		TG_Bandwidth.WithLabelValues(sbc.System, tg.Zone, tg.Name, "outbound").Set(tg.BandwidthOutboundUsage)
		TG_TotalChans.WithLabelValues(sbc.System, tg.Zone, tg.Name).Set(tg.TotalCallsConfigured)
		TG_State.WithLabelValues(sbc.System, tg.Zone, tg.Name).Set(trunkGroupStatus.stateToMetric(*tg))
		TG_OBState.WithLabelValues(sbc.System, tg.Zone, tg.Name).Set(trunkGroupStatus.outStateToMetric(*tg))
	}

	return nil
}
