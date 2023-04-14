package sonus

import (
	"context"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
	"golang.org/x/sync/errgroup"
)

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <sipCurrentStatistics xmlns="http://sonusnet.com/ns/mibs/SONUS-SIP-PEER-PERF-STATS/1.0">
    <name>TEST_Logan</name>
    <rcvInvite>0</rcvInvite>
    <sndInvite>0</sndInvite>
    <rcvAck>0</rcvAck>
    <sndAck>0</sndAck>
    <rcvPrack>0</rcvPrack>
...
</collection>
*/

type sipStatisticCollection struct {
	SipStatistics []*sipStatistics `xml:"http://sonusnet.com/ns/mibs/SONUS-SIP-PEER-PERF-STATS/1.0 sipCurrentStatistics,omitempty"`
}

type sipStatistics struct {
	Name                          string  `xml:"name"`
	RcvInvite                     float64 `xml:"rcvInvite"`
	SndInvite                     float64 `xml:"sndInvite"`
	RcvAck                        float64 `xml:"rcvAck"`
	SndAck                        float64 `xml:"sndAck"`
	RcvPrack                      float64 `xml:"rcvPrack"`
	SndPrack                      float64 `xml:"sndPrack"`
	RcvInfo                       float64 `xml:"rcvInfo"`
	SndInfo                       float64 `xml:"sndInfo"`
	RcvRefer                      float64 `xml:"rcvRefer"`
	SndRefer                      float64 `xml:"sndRefer"`
	RcvBye                        float64 `xml:"rcvBye"`
	SndBye                        float64 `xml:"sndBye"`
	RcvCancel                     float64 `xml:"rcvCancel"`
	SndCancel                     float64 `xml:"sndCancel"`
	RcvRegister                   float64 `xml:"rcvRegister"`
	SndRegister                   float64 `xml:"sndRegister"`
	RcvUpdate                     float64 `xml:"rcvUpdate"`
	SndUpdate                     float64 `xml:"sndUpdate"`
	Rcv18x                        float64 `xml:"rcv18x"`
	Snd18x                        float64 `xml:"snd18x"`
	Rcv1xx                        float64 `xml:"rcv1xx"`
	Snd1xx                        float64 `xml:"snd1xx"`
	Rcv2xx                        float64 `xml:"rcv2xx"`
	Snd2xx                        float64 `xml:"snd2xx"`
	RcvNonInv2xx                  float64 `xml:"rcvNonInv2xx"`
	SndNonInv2xx                  float64 `xml:"sndNonInv2xx"`
	Rcv3xx                        float64 `xml:"rcv3xx"`
	Snd3xx                        float64 `xml:"snd3xx"`
	Rcv4xx                        float64 `xml:"rcv4xx"`
	Snd4xx                        float64 `xml:"snd4xx"`
	Rcv5xx                        float64 `xml:"rcv5xx"`
	Snd5xx                        float64 `xml:"snd5xx"`
	Rcv6xx                        float64 `xml:"rcv6xx"`
	Snd6xx                        float64 `xml:"snd6xx"`
	RcvNonInvErr                  float64 `xml:"rcvNonInvErr"`
	SndNonInvErr                  float64 `xml:"sndNonInvErr"`
	RcvUnknownMsg                 float64 `xml:"rcvUnknownMsg"`
	RcvSubscriber                 float64 `xml:"rcvSubscriber"`
	SndSubscriber                 float64 `xml:"sndSubscriber"`
	RcvNotify                     float64 `xml:"rcvNotify"`
	SndNotify                     float64 `xml:"sndNotify"`
	RcvOption                     float64 `xml:"rcvOption"`
	SndOption                     float64 `xml:"sndOption"`
	InvReTransmit                 float64 `xml:"invReTransmit"`
	RegReTransmit                 float64 `xml:"regReTransmit"`
	ByeReTransmit                 float64 `xml:"byeRetransmit"`
	CancelReTransmit              float64 `xml:"cancelReTransmit"`
	OtherReTransmit               float64 `xml:"otherReTransmit"`
	RcvMessage                    float64 `xml:"rcvMessage"`
	SndMessage                    float64 `xml:"sndMessage"`
	RcvPublish                    float64 `xml:"rcvPublish"`
	SndPublish                    float64 `xml:"sndPublish"`
	EmergencyAccept               float64 `xml:"emergencyAccept"`
	EmergencyRejectBWCall         float64 `xml:"emergencyRejectBWCall"`
	EmergencyRejectPolicer        float64 `xml:"emergencyRejectPolicer"`
	HpcAccept                     float64 `xml:"hpcAccept"`
	NumberOfCallsSendingAARs      float64 `xml:"numberOfCallsSendingAARs"`
	NumberOfReceivedAAAFailures   float64 `xml:"numberOfReceivedAAAFailures"`
	NumberOfTotalAARSent          float64 `xml:"numberOfTotalAARSent"`
	NumberOfTimeoutOrErrorAAR     float64 `xml:"numberOfTimeoutOrErrorAAR"`
	EmergencyRegAccept            float64 `xml:"emergencyRegAccept"`
	EmergencyRegRejectLimit       float64 `xml:"emergencyRegRejectLimit"`
	EmergencyRegRejectPolicer     float64 `xml:"emergencyRegRejectPolicer"`
	NumberOfReceivedAAASuccesses  float64 `xml:"numberOfReceivedAAASuccesses"`
	NumberOfReceivedRARs          float64 `xml:"numberOfReceivedRARs"`
	NumberOfReceivedASRs          float64 `xml:"numberOfReceivedASRs"`
	NumberOfSentSTRs              float64 `xml:"numberOfSentSTRs"`
	EmergencyOODAccept            float64 `xml:"emergencyOODAccept"`
	EmergencyOODRejectPolicer     float64 `xml:"emergencyOODRejectPolicer"`
	EmergencySubsAccept           float64 `xml:"emergencySubsAccept"`
	EmergencySubsRejectLimit      float64 `xml:"emergencySubsRejectLimit"`
	EmergencySubsRejectPolicer    float64 `xml:"emergencySubsRejectPolicer"`
	ParseError                    float64 `xml:"parseError"`
	NumberOfTotalUDRSent          float64 `xml:"numberOfTotalUDRSent"`
	NumberOfTimeoutOrErrorUDR     float64 `xml:"numberOfTimeoutOrErrorUDR"`
	NumberOfReceivedUDASuccesses  float64 `xml:"numberOfReceivedUDASuccesses"`
	NumberOfReceivedUDAFailures   float64 `xml:"numberOfReceivedUDAFailures"`
	TotNumOfS8hrOutbndReg         float64 `xml:"totNumOfS8hrOutbndReg"`
	NumOfS8hrOutbndRegSuc         float64 `xml:"numOfS8hrOutbndRegSuc"`
	NumOfS8hrOutbndRegFail        float64 `xml:"numOfS8hrOutbndRegFail"`
	TotNumOfS8hrOutbndNormalCall  float64 `xml:"totNumOfS8hrOutbndNormalCall"`
	NumOfS8hrOutbndNormalCallSuc  float64 `xml:"numOfS8hrOutbndNormalCallSuc"`
	NumOfS8hrOutbndNormalCallFail float64 `xml:"numOfS8hrOutbndNormalCallFail"`
	NumOfS8hrOutbndEmgCallRej     float64 `xml:"numOfS8hrOutbndEmgCallRej"`
	NumOfS8hrInboundRegSuc        float64 `xml:"numOfS8hrInboundRegSuc"`
	NumOfS8hrInboundRegFail       float64 `xml:"numOfS8hrInboundRegFail"`
	NumOfS8hrInboundEmgCallSuc    float64 `xml:"numOfS8hrInboundEmgCallSuc"`
	NumOfS8hrInboundEmgCallFail   float64 `xml:"numOfS8hrInboundEmgCallFail"`
	InHpcAccept                   float64 `xml:"inHpcAccept"`
	OutHpcAccept                  float64 `xml:"outHpcAccept"`
	Hpc403Out                     float64 `xml:"hpc403Out"`
	HpcOverloadExempt             float64 `xml:"hpcOverloadExempt"`
}

func SIPMetrics(ctx context.Context, sbc *SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error {
	var (
		TG_SIP_Req_Sent = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "sonus_TG_sip_req_sent",
			Help: "Number of SIP requests sent",
		}, []string{"system", "zone", "name", "method"})

		TG_SIP_Req_Received = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "sonus_TG_sip_req_recv",
			Help: "Number of SIP requests received",
		}, []string{"system", "zone", "name", "method"})

		TG_SIP_Resp_Sent = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "sonus_TG_sip_resp_sent",
			Help: "Number of SIP responses sent",
		}, []string{"system", "zone", "name", "code"})

		TG_SIP_Resp_Received = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "sonus_TG_sip_resp_recv",
			Help: "Number of SIP responses received",
		}, []string{"system", "zone", "name", "code"})

		// sipStats = new(sipStatisticCollection)
	)

	registry.MustRegister(TG_SIP_Req_Received)
	registry.MustRegister(TG_SIP_Req_Sent)
	registry.MustRegister(TG_SIP_Resp_Received)
	registry.MustRegister(TG_SIP_Resp_Sent)

	g := &errgroup.Group{}
	g.SetLimit(3)

	for _, aCtx := range sbc.AddressContexts.AddressContext {
		for _, zone := range aCtx.Zone {
			g.Go(func() error {
				aCtx := aCtx
				zone := zone

				stats := &sipStatisticCollection{}
				err := sbc.GetAndParse(ctx, stats, sipStatsPath, aCtx.Name, zone.Name)
				if err != nil {
					return err
				}
				for _, sipStat := range stats.SipStatistics {
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "INVITE").Add(sipStat.SndInvite)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "PRACK").Add(sipStat.SndPrack)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "INFO").Add(sipStat.SndInfo)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "REFER").Add(sipStat.SndRefer)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "BYE").Add(sipStat.SndBye)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "CANCEL").Add(sipStat.SndCancel)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "REGISTER").Add(sipStat.SndRegister)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "UPDATE").Add(sipStat.SndUpdate)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "SUBSCRIBE").Add(sipStat.SndSubscriber)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "NOTIFY").Add(sipStat.SndNotify)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "OPTIONS").Add(sipStat.SndOption)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "MESSAGE").Add(sipStat.SndMessage)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "PUBLISH").Add(sipStat.SndPublish)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "INVITE (retrans)").Add(sipStat.InvReTransmit)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "REGISTER (retrans)").Add(sipStat.RegReTransmit)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "BYE (retrans)").Add(sipStat.ByeReTransmit)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "CANCEL (retrans)").Add(sipStat.CancelReTransmit)
					TG_SIP_Req_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "Other (retrans)").Add(sipStat.OtherReTransmit)

					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "INVITE").Add(sipStat.RcvInvite)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "PRACK").Add(sipStat.RcvPrack)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "INFO").Add(sipStat.RcvInfo)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "REFER").Add(sipStat.RcvRefer)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "BYE").Add(sipStat.RcvBye)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "CANCEL").Add(sipStat.RcvCancel)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "REGISTER").Add(sipStat.RcvRegister)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "UPDATE").Add(sipStat.RcvUpdate)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "SUBSCRIBE").Add(sipStat.RcvSubscriber)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "NOTIFY").Add(sipStat.RcvNotify)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "OPTIONS").Add(sipStat.RcvOption)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "MESSAGE").Add(sipStat.RcvMessage)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "PUBLISH").Add(sipStat.RcvPublish)
					TG_SIP_Req_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "Unknown").Add(sipStat.RcvUnknownMsg)

					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "ACK").Add(sipStat.SndAck)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "18x").Add(sipStat.Snd18x)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "1xx").Add(sipStat.Snd1xx)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "2xx").Add(sipStat.Snd2xx)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "Non-INVITE 2xx").Add(sipStat.SndNonInv2xx)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "3xx").Add(sipStat.Snd3xx)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "4xx").Add(sipStat.Snd4xx)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "5xx").Add(sipStat.Snd5xx)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "6xx").Add(sipStat.Snd6xx)
					TG_SIP_Resp_Sent.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "Non-INVITE error").Add(sipStat.SndNonInvErr)

					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "ACK").Add(sipStat.RcvAck)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "18x").Add(sipStat.Rcv18x)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "1xx").Add(sipStat.Rcv1xx)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "2xx").Add(sipStat.Rcv2xx)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "Non-INVITE 2xx").Add(sipStat.RcvNonInv2xx)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "3xx").Add(sipStat.Rcv3xx)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "4xx").Add(sipStat.Rcv4xx)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "5xx").Add(sipStat.Rcv5xx)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "6xx").Add(sipStat.Rcv6xx)
					TG_SIP_Resp_Received.WithLabelValues(sbc.System, zone.Name, sipStat.Name, "Non-INVITE error").Add(sipStat.RcvNonInvErr)
				}
				return nil
			})
		}
	}
	return g.Wait()
}
