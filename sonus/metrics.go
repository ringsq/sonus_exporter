package sonus

import (
	"fmt"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
)

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

// BuildMetrics takes a registry and a structure and creates Gauge metrics for any float64 items.
// The metric names are based on the fields in the structure and any substructures.
// eg. structname_structname_fieldName
//
// The returned map is keyed  by the metric name
func BuildMetrics(registry *prometheus.Registry, t reflect.Type) map[string]*prometheus.GaugeVec {
	ch := make(chan *Metric)
	stats := map[string]*prometheus.GaugeVec{}
	go func() {
		examiner(ch, "sonus", t)
		close(ch)
	}()
	for x := range ch {
		registry.MustRegister(x.Metric)
		stats[x.Name] = x.Metric
	}

	return stats
}

func examiner(ch chan *Metric, sub string, t reflect.Type) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		examiner(ch, fmt.Sprintf("%s_%s", sub, t.Name()), t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			switch f.Type.Kind() {
			case reflect.Slice:
				examiner(ch, fmt.Sprintf("%s_%s", sub, f.Name), f.Type.Elem())
			case reflect.Float64:
				// ch <- fmt.Sprintf("%s_%s (%s)\n", sub, f.Name, f.Type.Kind().String())
				metric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: sub,
					Name:      f.Name,
					Help:      getHelp(f.Name),
				}, []string{"system", "addresscontext", "zone", "trunkgroup"})
				m := &Metric{Name: fmt.Sprintf("%s_%s", sub, f.Name), Metric: metric}
				ch <- m
			case reflect.Struct:
				examiner(ch, fmt.Sprintf("%s_%s", sub, f.Name), f.Type)
			}
		}
	}
}
