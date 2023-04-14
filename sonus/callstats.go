package sonus

import (
	"context"

	"github.com/fatih/structs"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
	"golang.org/x/sync/errgroup"
)

type Calls struct {
	CallCurrentStatistics []callStats `xml:"callCurrentStatistics"`
}

type callStats struct {
	Name                       string  `xml:"name"`
	InUsage                    float64 `xml:"inUsage"`
	OutUsage                   float64 `xml:"outUsage"`
	InCalls                    float64 `xml:"inCalls"`
	OutCalls                   float64 `xml:"outCalls"`
	InCallAttempts             float64 `xml:"inCallAttempts"`
	OutCallAttempts            float64 `xml:"outCallAttempts"`
	MaxCompletedCalls          float64 `xml:"maxCompletedCalls"`
	CallSetupTime              float64 `xml:"callSetupTime"`
	CallSetups                 float64 `xml:"callSetups"`
	RoutingAttempts            float64 `xml:"routingAttempts"`
	InBwUsage                  float64 `xml:"inBwUsage"`
	OutBwUsage                 float64 `xml:"outBwUsage"`
	MaxActiveBwUsage           float64 `xml:"maxActiveBwUsage"`
	CallsWithPktOutage         float64 `xml:"callsWithPktOutage"`
	CallsWithPktOutageAtEnd    float64 `xml:"callsWithPktOutageAtEnd"`
	TotalPktOutage             float64 `xml:"totalPktOutage"`
	MaxPktOutage               float64 `xml:"maxPktOutage"`
	PodEvents                  float64 `xml:"podEvents"`
	PlayoutBufferGood          float64 `xml:"playoutBufferGood"`
	PlayoutBufferAcceptable    float64 `xml:"playoutBufferAcceptable"`
	PlayoutBufferPoor          float64 `xml:"playoutBufferPoor"`
	PlayoutBufferUnacceptable  float64 `xml:"playoutBufferUnacceptable"`
	SipRegAttempts             float64 `xml:"sipRegAttempts"`
	SipRegCompletions          float64 `xml:"sipRegCompletions"`
	CallsWithPsxDips           float64 `xml:"callsWithPsxDips"`
	TotalPsxDips               float64 `xml:"totalPsxDips"`
	ActiveRegs                 float64 `xml:"activeRegs"`
	MaxActiveRegs              float64 `xml:"maxActiveRegs"`
	ActiveSubs                 float64 `xml:"activeSubs"`
	MaxActiveSubs              float64 `xml:"maxActiveSubs"`
	PeakCallRate               float64 `xml:"peakCallRate"`
	TotalOnGoingCalls          float64 `xml:"totalOnGoingCalls"`
	TotalStableCalls           float64 `xml:"totalStableCalls"`
	TotalCallUpdates           float64 `xml:"totalCallUpdates"`
	TotalEmergencyStableCalls  float64 `xml:"totalEmergencyStableCalls"`
	TotalEmergencyOnGoingCalls float64 `xml:"totalEmergencyOnGoingCalls"`
	InRetargetCalls            float64 `xml:"inRetargetCalls"`
	InRetargetRegs             float64 `xml:"inRetargetRegs"`
	OutRetargetCalls           float64 `xml:"outRetargetCalls"`
	OutRetargetRegs            float64 `xml:"outRetargetRegs"`
}

// metricHelp contains a table that describes the values in the Sonus response.
// It is used for building the HELP line for the corresponding metric.
var metricHelp = map[string]string{
	"InUsage":                    "The current usage in the inbound direction for this trunk group in seconds. Usage is defined as the time media bandwidth is activated to the time it is deactivated.",                                                                                                                                                                                                        // "inUsage"`
	"OutUsage":                   "The current usage in the outbound direction for this trunk group in seconds. Usage is defined as the time media bandwidth is activated to the time it is deactivated.",                                                                                                                                                                                                       // "outUsage"`
	"InCalls":                    "The current number of completed inbound calls on this trunk group.",                                                                                                                                                                                                                                                                                                          // "inCalls"`
	"OutCalls":                   "The current number of completed outbound calls on this trunk group.",                                                                                                                                                                                                                                                                                                         // "outCalls"`
	"InCallAttempts":             "The current number of inbound call attempts on this trunk group.",                                                                                                                                                                                                                                                                                                            // "inCallAttempts"`
	"OutCallAttempts":            "The current number of outbound call attempts on this trunk group.",                                                                                                                                                                                                                                                                                                           // "outCallAttempts"`
	"MaxCompletedCalls":          "Displayed as maxActiveCalls. The current high water mark of total number of active calls in both the inbound and outbound directions on the trunk group. This statistic accounts for calls that are setting up, stable, or tearing down.",                                                                                                                                    // "maxCompletedCalls"`
	"CallSetupTime":              "The cumulative duration (in 100th's of a seconds) from an INVITE sent to receiving the first backward 18x response on the egress leg. This value is nearly identical on the ingress counter with any latency due to the time spent for the SBC to send out the received 18x. If no 18x response is present, the callSetupTime is the final 200 response (cumulative count).", // "callSetupTime"`
	"CallSetups":                 "The current total number of calls setup but not necessarily completed in the inbound and outbound directions for this trunk group. This object can be used as the denominator for calculating average call setup time.",                                                                                                                                                      // "callSetups"`
	"RoutingAttempts":            "The current number of routing attempts for this trunk group.",                                                                                                                                                                                                                                                                                                                // "routingAttempts"`
	"InBwUsage":                  "The sum of BW usage (expected data rate in Kbits per second multiplied by call duration in seconds) for every inbound call associated with this trunk group.",                                                                                                                                                                                                                // "inBwUsage"`
	"OutBwUsage":                 "The sum of BW usage (expected data rate in Kbits per second multiplied by call duration in seconds) for every outbound call associated with this trunk group.",                                                                                                                                                                                                               // "outBwUsage"`
	"MaxActiveBwUsage":           "The high water mark of BW usage in either direction associated with this trunk group.",                                                                                                                                                                                                                                                                                       // "maxActiveBwUsage"`
	"CallsWithPktOutage":         "The number of calls with a maximum packet outage whose duration exceeds the configured minimum for this trunk group.",                                                                                                                                                                                                                                                        // "callsWithPktOutage"`
	"CallsWithPktOutageAtEnd":    "The number of calls whose maximum packet outage occurs at the end of the call for this trunk group. This is an indication that the call may have been terminated the because of poor quality.",                                                                                                                                                                               // "callsWithPktOutageAtEnd"`
	"TotalPktOutage":             "The summation of all packet outage durations (in milliseconds) whose duration exceeds the configured minimum, which is experienced during the current performance interval for this trunk group. The average packet outage duration can be calculated by dividing this field by the number of calls reporting packet outages.",                                               // "totalPktOutage"`
	"MaxPktOutage":               "The single longest maximum reported packet outage duration (in milliseconds) experienced during the current performance interval for this trunk group.",                                                                                                                                                                                                                      // "maxPktOutage"`
	"PodEvents":                  "The number of Packet Outage Detection (POD) Events detected for this trunk group. A POD event occurs when a configurable number of calls experience a packet outage with duration exceeding a programmable threshold.",                                                                                                                                                       // "podEvents"`
	"PlayoutBufferGood":          "Number of calls with all sub-intervals reporting GOOD playout buffer quality for this trunk group.",                                                                                                                                                                                                                                                                          // "playoutBufferGood"`
	"PlayoutBufferAcceptable":    "Number of calls with all sub-intervals reporting ACCEPTABLE or better playout buffer quality for this trunk group.",                                                                                                                                                                                                                                                          // "playoutBufferAcceptable"`
	"PlayoutBufferPoor":          "Number of calls with all sub-intervals reporting POOR or better playout buffer quality for this trunk group.",                                                                                                                                                                                                                                                                // "playoutBufferPoor"`
	"PlayoutBufferUnacceptable":  "Number of calls with at least one sub-interval reporting UNACCEPTABLE playout buffer quality for this trunk group.",                                                                                                                                                                                                                                                          // "playoutBufferUnacceptable"`
	"SipRegAttempts":             "The current number of SIP registration attempts on a trunk group.",                                                                                                                                                                                                                                                                                                           // "sipRegAttempts"`
	"SipRegCompletions":          "The current number of SIP registrations that have successfully completed on a trunk group.",                                                                                                                                                                                                                                                                                  // "sipRegCompletions"`
	"CallsWithPsxDips":           "The current number of calls that made a PSX Dip",                                                                                                                                                                                                                                                                                                                             // "callsWithPsxDips"`
	"TotalPsxDips":               "The current number of PSX Dips made.",                                                                                                                                                                                                                                                                                                                                        // "totalPsxDips"`
	"ActiveRegs":                 "The current number of active registrations on this trunk group.",                                                                                                                                                                                                                                                                                                             // "activeRegs"`
	"MaxActiveRegs":              "The current number of maximum active registrations on this trunk group (this is the high-watermark achieved on this TG).",                                                                                                                                                                                                                                                    // "maxActiveRegs"`
	"ActiveSubs":                 "The current number of active subscriptions on this trunk group.",                                                                                                                                                                                                                                                                                                             // "activeSubs"`
	"MaxActiveSubs":              "The current number of maximum active subscriptions on this trunk group (this is the high-watermark achieved on this TG).",                                                                                                                                                                                                                                                    // "maxActiveSubs"`
	"PeakCallRate":               "Peak call arrival rate for the current interval on this trunk group",                                                                                                                                                                                                                                                                                                         // "peakCallRate"`
	"TotalOnGoingCalls":          "Total Calls (Non-Stable + Stable) on this trunk group",                                                                                                                                                                                                                                                                                                                       // "totalOnGoingCalls"`
	"TotalStableCalls":           "Total Stable Calls on this trunk group",                                                                                                                                                                                                                                                                                                                                      // "totalStableCalls"`
	"TotalCallUpdates":           "Total Call Updates on this trunk group",                                                                                                                                                                                                                                                                                                                                      // "totalCallUpdates"`
	"TotalEmergencyStableCalls":  "Total Emergency Stable Calls on this trunk group",                                                                                                                                                                                                                                                                                                                            // "totalEmergencyStableCalls"`
	"TotalEmergencyOnGoingCalls": "Total Emergency Calls in establishing state on this trunk group",                                                                                                                                                                                                                                                                                                             // "totalEmergencyOnGoingCalls"`
	"InRetargetCalls":            "The current number of incoming calls that are retargeted by Load Balancing Service",                                                                                                                                                                                                                                                                                          // "inRetargetCalls"`
	"InRetargetRegs":             "The current number of incoming registrations that are retargeted by Load Balancing Service",                                                                                                                                                                                                                                                                                  // "inRetargetRegs"`
	"OutRetargetCalls":           "The current number of outgoing calls that are retargeted by Load Balancing Service",                                                                                                                                                                                                                                                                                          // "outRetargetCalls"`
	"OutRetargetRegs":            "The current number of outgoing registrations that are retargeted by Load Balancing Service",                                                                                                                                                                                                                                                                                  // "outRetargetRegs"`
}

func getHelp(field string) string {
	help, ok := metricHelp[field]
	if !ok || len(help) == 0 {
		return field
	}
	return help
}

func CallMetrics(ctx context.Context, sbc *SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error {
	var (
		callStats = new(callStats)
	)

	stats := map[string]*prometheus.GaugeVec{}
	for _, field := range structs.Names(callStats) {
		if field == "Name" {
			continue
		}
		metric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sonus",
			Subsystem: "TG_calls",
			Name:      field,
			Help:      getHelp(field),
		}, []string{"system", "addresscontext", "zone", "name"})
		registry.MustRegister(metric)
		stats[field] = metric
	}

	g := &errgroup.Group{}
	g.SetLimit(3)

	for _, aCtx := range sbc.AddressContexts.AddressContext {
		for _, zone := range aCtx.Zone {
			g.Go(func() error {
				aCtx := aCtx
				zone := zone

				calls := &Calls{}
				err := sbc.GetAndParse(ctx, calls, callStatusPath, aCtx.Name, zone.Name)
				if err != nil {
					return err
				}
				for _, sipStat := range calls.CallCurrentStatistics {
					for _, field := range structs.Fields(sipStat) {
						if field.Kind().String() != "float64" {
							continue
						}
						metric, ok := stats[field.Name()]
						if !ok {
							continue
						}
						metric.WithLabelValues(sbc.System, aCtx.Name, zone.Name, sipStat.Name).Set(field.Value().(float64))
					}
				}
				return nil
			})
		}
	}
	return g.Wait()
}
