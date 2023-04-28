package sonus

import (
	"fmt"
	"reflect"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	Name   string
	Metric *prometheus.GaugeVec
}

type ZoneStats struct {
	Zone []struct {
		ZONE                   string `xml:"ZONE,attr"`
		ADDRESSCONTEXT         string `xml:"ADDRESS_CONTEXT,attr"`
		Name                   string `xml:"name"`
		ID                     string `xml:"id"`
		FilterSipSrcStatistics struct {
			Name                string `xml:"name"`
			FilteredSipMsgCount string `xml:"filteredSipMsgCount"`
		} `xml:"filterSipSrcStatistics"`
		SipSigPort struct {
			Index                     string `xml:"index"`
			IpInterfaceGroupName      string `xml:"ipInterfaceGroupName"`
			IpAddressV4               string `xml:"ipAddressV4"`
			Mode                      string `xml:"mode"`
			State                     string `xml:"state"`
			TransportProtocolsAllowed string `xml:"transportProtocolsAllowed"`
			MonitoringProfileName     string `xml:"monitoringProfileName"`
			PortNumber                string `xml:"portNumber"`
			TlsProfileName            string `xml:"tlsProfileName"`
		} `xml:"sipSigPort"`
		SipSigPortStatus struct {
			Index        string `xml:"index"`
			State        string `xml:"state"`
			LocalIpType  string `xml:"localIpType"`
			FixedIpV4    string `xml:"fixedIpV4"`
			FixedIpV6    string `xml:"fixedIpV6"`
			FloatingIpV4 string `xml:"floatingIpV4"`
			FloatingIpV6 string `xml:"floatingIpV6"`
			PortNumber   string `xml:"portNumber"`
			VnfId        string `xml:"vnfId"`
		} `xml:"sipSigPortStatus"`
		SipSigPortStatistics struct {
			Index     string  `xml:"index"`
			CallRate  float64 `xml:"callRate"`
			OrigCalls float64 `xml:"origCalls"`
			TermCalls float64 `xml:"termCalls"`
			TxPdus    float64 `xml:"txPdus"`
			RxPdus    float64 `xml:"rxPdus"`
			TxBytes   float64 `xml:"txBytes"`
			RxBytes   float64 `xml:"rxBytes"`
			InRegs    float64 `xml:"inRegs"`
			OutRegs   float64 `xml:"outRegs"`
			Tx500s    float64 `xml:"tx500s"`
			Tx503s    float64 `xml:"tx503s"`
		} `xml:"sipSigPortStatistics"`
		SipSigConnStatus []struct {
			ConnectionId  string  `xml:"connectionId"`
			Index         string  `xml:"index"`
			PeerIpAddress string  `xml:"peerIpAddress"`
			PeerPortNum   string  `xml:"peerPortNum"`
			Socket        string  `xml:"socket"`
			Transport     string  `xml:"transport"`
			State         string  `xml:"state"`
			Role          string  `xml:"role"`
			Aging         string  `xml:"aging"`
			IdleTime      string  `xml:"idleTime"`
			BytesSent     float64 `xml:"bytesSent"`
			BytesRcvd     float64 `xml:"bytesRcvd"`
			PduSendQueued float64 `xml:"pduSendQueued"`
			PduRecvQueued float64 `xml:"pduRecvQueued"`
		} `xml:"sipSigConnStatus"`
		SipSigConnStatistics struct {
			Index                  string  `xml:"index"`
			TcpConnection          string  `xml:"tcpConnection"`
			TotalTcpConnection     float64 `xml:"totalTcpConnection"`
			ActiveTlsTcpConnection float64 `xml:"activeTlsTcpConnection"`
			TotalTlsTcpConnection  float64 `xml:"totalTlsTcpConnection"`
		} `xml:"sipSigConnStatistics"`
		SipSigPortTlsStatistics struct {
			Index                    string  `xml:"index"`
			CurrentServerSessions    float64 `xml:"currentServerSessions"`
			TotalServerSessions      float64 `xml:"totalServerSessions"`
			CurrentClientHandshakes  float64 `xml:"currentClientHandshakes"`
			CurrentServerHandshakes  float64 `xml:"currentServerHandshakes"`
			SessionResumptions       float64 `xml:"sessionResumptions"`
			NoCipherSuite            float64 `xml:"noCipherSuite"`
			HandshakeTimeouts        float64 `xml:"handshakeTimeouts"`
			HigherAuthTimeout        float64 `xml:"higherAuthTimeout"`
			ClientAuthFailures       float64 `xml:"clientAuthFailures"`
			ServerAuthFailures       float64 `xml:"serverAuthFailures"`
			FatelAlertsReceived      float64 `xml:"fatelAlertsReceived"`
			WarningAlertsReceived    float64 `xml:"warningAlertsReceived"`
			HandshakeFailures        float64 `xml:"handshakeFailures"`
			ReceiveFailures          float64 `xml:"receiveFailures"`
			SendFailures             float64 `xml:"sendFailures"`
			NoAuthDrops              float64 `xml:"noAuthDrops"`
			NoAuth488                float64 `xml:"noAuth488"`
			MidConnectionHello       float64 `xml:"midConnectionHello"`
			NoClientCert             float64 `xml:"noClientCert"`
			ValidationFailures       float64 `xml:"validationFailures"`
			CurrentClientConnections float64 `xml:"currentClientConnections"`
			TotalClientConnections   float64 `xml:"totalClientConnections"`
			CurrentServerConnections float64 `xml:"currentServerConnections"`
			TotalServerConnections   float64 `xml:"totalServerConnections"`
		} `xml:"sipSigPortTlsStatistics"`
		MessageManipulation struct {
			OutputAdapterProfile string `xml:"outputAdapterProfile"`
		} `xml:"messageManipulation"`
		SipCurrentStatistics []struct {
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
			ByeRetransmit                 float64 `xml:"byeRetransmit"`
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
		} `xml:"sipCurrentStatistics"`
		SipIntervalStatistics []struct {
			Number                        string  `xml:"number"`
			Name                          string  `xml:"name"`
			IntervalValid                 string  `xml:"intervalValid"`
			Time                          string  `xml:"time"`
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
			ByeRetransmit                 float64 `xml:"byeRetransmit"`
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
		} `xml:"sipIntervalStatistics"`
		SipOcsCallCurrentStatistics []struct {
			Name             string `xml:"name"`
			AttemptedCalls   string `xml:"attemptedCalls"`
			RelayedCalls     string `xml:"relayedCalls"`
			EstablishedCalls string `xml:"establishedCalls"`
			SuccessfulCalls  string `xml:"successfulCalls"`
			FailedCalls      string `xml:"failedCalls"`
			PendingCalls     string `xml:"pendingCalls"`
			RejectedCalls    string `xml:"rejectedCalls"`
		} `xml:"sipOcsCallCurrentStatistics"`
		SipOcsCallIntervalStatistics []struct {
			Number           string `xml:"number"`
			Name             string `xml:"name"`
			IntervalValid    string `xml:"intervalValid"`
			Time             string `xml:"time"`
			AttemptedCalls   string `xml:"attemptedCalls"`
			RelayedCalls     string `xml:"relayedCalls"`
			EstablishedCalls string `xml:"establishedCalls"`
			SuccessfulCalls  string `xml:"successfulCalls"`
			FailedCalls      string `xml:"failedCalls"`
			PendingCalls     string `xml:"pendingCalls"`
			RejectedCalls    string `xml:"rejectedCalls"`
		} `xml:"sipOcsCallIntervalStatistics"`
		SipInviteResponseCurrentStatistics []struct {
			Name        string  `xml:"name"`
			Response401 float64 `xml:"response401"`
			Response403 float64 `xml:"response403"`
			Response407 float64 `xml:"response407"`
			Response481 float64 `xml:"response481"`
		} `xml:"sipInviteResponseCurrentStatistics"`
		SipRegisterResponseCurrentStatistics []struct {
			Name        string `xml:"name"`
			Response401 string `xml:"response401"`
			Response403 string `xml:"response403"`
			Response407 string `xml:"response407"`
			Response481 string `xml:"response481"`
		} `xml:"sipRegisterResponseCurrentStatistics"`
		SipByeResponseCurrentStatistics []struct {
			Name        string  `xml:"name"`
			Response401 float64 `xml:"response401"`
			Response403 float64 `xml:"response403"`
			Response407 float64 `xml:"response407"`
			Response481 float64 `xml:"response481"`
		} `xml:"sipByeResponseCurrentStatistics"`
		SipOptionResponseCurrentStatistics []struct {
			Name        string  `xml:"name"`
			Response401 float64 `xml:"response401"`
			Response403 float64 `xml:"response403"`
			Response407 float64 `xml:"response407"`
			Response481 float64 `xml:"response481"`
		} `xml:"sipOptionResponseCurrentStatistics"`
		SipInviteResponseIntervalStatistics []struct {
			Number        float64 `xml:"number"`
			Name          string  `xml:"name"`
			IntervalValid string  `xml:"intervalValid"`
			Time          float64 `xml:"time"`
			Response401   float64 `xml:"response401"`
			Response403   float64 `xml:"response403"`
			Response407   float64 `xml:"response407"`
			Response481   float64 `xml:"response481"`
		} `xml:"sipInviteResponseIntervalStatistics"`
		SipRegisterResponseIntervalStatistics []struct {
			Number        string  `xml:"number"`
			Name          string  `xml:"name"`
			IntervalValid string  `xml:"intervalValid"`
			Time          float64 `xml:"time"`
			Response401   float64 `xml:"response401"`
			Response403   float64 `xml:"response403"`
			Response407   float64 `xml:"response407"`
			Response481   float64 `xml:"response481"`
		} `xml:"sipRegisterResponseIntervalStatistics"`
		SipByeResponseIntervalStatistics []struct {
			Number        string  `xml:"number"`
			Name          string  `xml:"name"`
			IntervalValid string  `xml:"intervalValid"`
			Time          float64 `xml:"time"`
			Response401   float64 `xml:"response401"`
			Response403   float64 `xml:"response403"`
			Response407   float64 `xml:"response407"`
			Response481   float64 `xml:"response481"`
		} `xml:"sipByeResponseIntervalStatistics"`
		SipOptionResponseIntervalStatistics []struct {
			Number        string  `xml:"number"`
			Name          string  `xml:"name"`
			IntervalValid string  `xml:"intervalValid"`
			Time          float64 `xml:"time"`
			Response401   float64 `xml:"response401"`
			Response403   float64 `xml:"response403"`
			Response407   float64 `xml:"response407"`
			Response481   float64 `xml:"response481"`
		} `xml:"sipOptionResponseIntervalStatistics"`
		CallCurrentStatistics []struct {
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
		} `xml:"callCurrentStatistics"`
		CallIntervalStatistics []struct {
			Number                     string  `xml:"number"`
			Name                       string  `xml:"name"`
			IntervalValid              string  `xml:"intervalValid"`
			Time                       float64 `xml:"time"`
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
		} `xml:"callIntervalStatistics"`
		CallFailureCurrentStatistics []struct {
			Name                      string  `xml:"name"`
			InCallFailNoRoutes        float64 `xml:"inCallFailNoRoutes"`
			InCallFailNoResources     float64 `xml:"inCallFailNoResources"`
			InCallFailNoService       float64 `xml:"inCallFailNoService"`
			InCallFailInvalidCall     float64 `xml:"inCallFailInvalidCall"`
			InCallFailNetworkFailure  float64 `xml:"inCallFailNetworkFailure"`
			InCallFailProtocolError   float64 `xml:"inCallFailProtocolError"`
			InCallFailUnspecified     float64 `xml:"inCallFailUnspecified"`
			OutCallFailNoRoutes       float64 `xml:"outCallFailNoRoutes"`
			OutCallFailNoResources    float64 `xml:"outCallFailNoResources"`
			OutCallFailNoService      float64 `xml:"outCallFailNoService"`
			OutCallFailInvalidCall    float64 `xml:"outCallFailInvalidCall"`
			OutCallFailNetworkFailure float64 `xml:"outCallFailNetworkFailure"`
			OutCallFailProtocolError  float64 `xml:"outCallFailProtocolError"`
			OutCallFailUnspecified    float64 `xml:"outCallFailUnspecified"`
			RoutingFailuresResv       float64 `xml:"routingFailuresResv"`
			AllocFailBwLimit          float64 `xml:"allocFailBwLimit"`
			AllocFailCallLimit        float64 `xml:"allocFailCallLimit"`
			NoPsxRoute                float64 `xml:"noPsxRoute"`
			CallAbandoned             float64 `xml:"callAbandoned"`
			CallFailPolicing          float64 `xml:"callFailPolicing"`
			SipRegFailPolicing        float64 `xml:"sipRegFailPolicing"`
			SipRegFailInternal        float64 `xml:"sipRegFailInternal"`
			SipRegFailOther           float64 `xml:"sipRegFailOther"`
			SecurityFail              float64 `xml:"securityFail"`
			RegCallsFailed            float64 `xml:"regCallsFailed"`
			NonMatchSrcIpCallsFail    float64 `xml:"nonMatchSrcIpCallsFail"`
			InvalidSPCallsFailed      float64 `xml:"invalidSPCallsFailed"`
			AllocFailParentConstraint float64 `xml:"allocFailParentConstraint"`
			SipSubsFailPolicing       float64 `xml:"sipSubsFailPolicing"`
			SipOtherReqFailPolicing   float64 `xml:"sipOtherReqFailPolicing"`
			VideoThresholdLimit       float64 `xml:"videoThresholdLimit"`
			SipOtherReqFailInternal   float64 `xml:"sipOtherReqFailInternal"`
			SipOtherReqFailOther      float64 `xml:"sipOtherReqFailOther"`
		} `xml:"callFailureCurrentStatistics"`
		CallFailureIntervalStatistics []struct {
			Number                    string  `xml:"number"`
			Name                      string  `xml:"name"`
			IntervalValid             string  `xml:"intervalValid"`
			Time                      float64 `xml:"time"`
			InCallFailNoRoutes        float64 `xml:"inCallFailNoRoutes"`
			InCallFailNoResources     float64 `xml:"inCallFailNoResources"`
			InCallFailNoService       float64 `xml:"inCallFailNoService"`
			InCallFailInvalidCall     float64 `xml:"inCallFailInvalidCall"`
			InCallFailNetworkFailure  float64 `xml:"inCallFailNetworkFailure"`
			InCallFailProtocolError   float64 `xml:"inCallFailProtocolError"`
			InCallFailUnspecified     float64 `xml:"inCallFailUnspecified"`
			OutCallFailNoRoutes       float64 `xml:"outCallFailNoRoutes"`
			OutCallFailNoResources    float64 `xml:"outCallFailNoResources"`
			OutCallFailNoService      float64 `xml:"outCallFailNoService"`
			OutCallFailInvalidCall    float64 `xml:"outCallFailInvalidCall"`
			OutCallFailNetworkFailure float64 `xml:"outCallFailNetworkFailure"`
			OutCallFailProtocolError  float64 `xml:"outCallFailProtocolError"`
			OutCallFailUnspecified    float64 `xml:"outCallFailUnspecified"`
			RoutingFailuresResv       float64 `xml:"routingFailuresResv"`
			AllocFailBwLimit          float64 `xml:"allocFailBwLimit"`
			AllocFailCallLimit        float64 `xml:"allocFailCallLimit"`
			NoPsxRoute                float64 `xml:"noPsxRoute"`
			CallAbandoned             float64 `xml:"callAbandoned"`
			CallFailPolicing          float64 `xml:"callFailPolicing"`
			SipRegFailPolicing        float64 `xml:"sipRegFailPolicing"`
			SipRegFailInternal        float64 `xml:"sipRegFailInternal"`
			SipRegFailOther           float64 `xml:"sipRegFailOther"`
			SecurityFail              float64 `xml:"securityFail"`
			NonMatchSrcIpCallsFail    float64 `xml:"nonMatchSrcIpCallsFail"`
			InvalidSPCallsFailed      float64 `xml:"invalidSPCallsFailed"`
			AllocFailParentConstraint float64 `xml:"allocFailParentConstraint"`
			SipSubsFailPolicing       float64 `xml:"sipSubsFailPolicing"`
			SipOtherReqFailPolicing   float64 `xml:"sipOtherReqFailPolicing"`
			RegCallsFailed            float64 `xml:"regCallsFailed"`
			VideoThresholdLimit       float64 `xml:"videoThresholdLimit"`
			SipOtherReqFailInternal   float64 `xml:"sipOtherReqFailInternal"`
			SipOtherReqFailOther      float64 `xml:"sipOtherReqFailOther"`
		} `xml:"callFailureIntervalStatistics"`
		TrafficControlCurrentStatistics []struct {
			Name              string  `xml:"name"`
			Silc              float64 `xml:"silc"`
			StrCant           float64 `xml:"strCant"`
			StrSkip           float64 `xml:"strSkip"`
			Skip              float64 `xml:"skip"`
			Cant              float64 `xml:"cant"`
			Canf              float64 `xml:"canf"`
			AccCant           float64 `xml:"accCant"`
			AccSkip           float64 `xml:"accSkip"`
			RouteAttemptsIRR  float64 `xml:"routeAttemptsIRR"`
			RouteAttemptsSIRR float64 `xml:"routeAttemptsSIRR"`
			RouteAttemptsORR  float64 `xml:"routeAttemptsORR"`
			RouteAttemptsSORR float64 `xml:"routeAttemptsSORR"`
			SuccessfulIRR     float64 `xml:"successfulIRR"`
			SuccessfulSIRR    float64 `xml:"successfulSIRR"`
			SuccessfulORR     float64 `xml:"successfulORR"`
			SuccessfulSORR    float64 `xml:"successfulSORR"`
		} `xml:"trafficControlCurrentStatistics"`
		TrafficControlIntervalStatistics []struct {
			Number            string  `xml:"number"`
			Name              string  `xml:"name"`
			IntervalValid     string  `xml:"intervalValid"`
			Time              float64 `xml:"time"`
			Silc              float64 `xml:"silc"`
			StrCant           float64 `xml:"strCant"`
			StrSkip           float64 `xml:"strSkip"`
			Skip              float64 `xml:"skip"`
			Cant              float64 `xml:"cant"`
			Canf              float64 `xml:"canf"`
			AccCant           float64 `xml:"accCant"`
			AccSkip           float64 `xml:"accSkip"`
			RouteAttemptsIRR  float64 `xml:"routeAttemptsIRR"`
			RouteAttemptsSIRR float64 `xml:"routeAttemptsSIRR"`
			RouteAttemptsORR  float64 `xml:"routeAttemptsORR"`
			RouteAttemptsSORR float64 `xml:"routeAttemptsSORR"`
			SuccessfulIRR     float64 `xml:"successfulIRR"`
			SuccessfulSIRR    float64 `xml:"successfulSIRR"`
			SuccessfulORR     float64 `xml:"successfulORR"`
			SuccessfulSORR    float64 `xml:"successfulSORR"`
		} `xml:"trafficControlIntervalStatistics"`
		SipTrunkGroup []struct {
			Name   string `xml:"name"`
			State  string `xml:"state"`
			Mode   string `xml:"mode"`
			Policy struct {
				Carrier                string `xml:"carrier"`
				Country                string `xml:"country"`
				LocalizationVariant    string `xml:"localizationVariant"`
				TgIPVersionPreference  string `xml:"tgIPVersionPreference"`
				PreferredIdentity      string `xml:"preferredIdentity"`
				DigitParameterHandling struct {
					NumberingPlan string `xml:"numberingPlan"`
				} `xml:"digitParameterHandling"`
				CallRouting struct {
					ElementRoutingPriority string `xml:"elementRoutingPriority"`
				} `xml:"callRouting"`
				Media struct {
					PacketServiceProfile string `xml:"packetServiceProfile"`
				} `xml:"media"`
				Services struct {
					ClassOfService string `xml:"classOfService"`
				} `xml:"services"`
				Signaling struct {
					IpSignalingProfile string `xml:"ipSignalingProfile"`
					SignalingProfile   string `xml:"signalingProfile"`
				} `xml:"signaling"`
				FeatureControlProfile string `xml:"featureControlProfile"`
				IpSignalingPeerGroup  string `xml:"ipSignalingPeerGroup"`
				Ingress               struct {
					Flags struct {
						NonZeroVideoBandwidthBasedRoutingForSip  string `xml:"nonZeroVideoBandwidthBasedRoutingForSip"`
						NonZeroVideoBandwidthBasedRoutingForH323 string `xml:"nonZeroVideoBandwidthBasedRoutingForH323"`
						HdPreferredRouting                       string `xml:"hdPreferredRouting"`
						HdSupportedRouting                       string `xml:"hdSupportedRouting"`
					} `xml:"flags"`
				} `xml:"ingress"`
			} `xml:"policy"`
			Cac struct {
				CallLimit float64 `xml:"callLimit"`
				Ingress   struct {
					CallRateMax  float64 `xml:"callRateMax"`
					CallBurstMax float64 `xml:"callBurstMax"`
				} `xml:"ingress"`
			} `xml:"cac"`
			Signaling struct {
				MessageManipulation struct {
					OutputAdapterProfile string `xml:"outputAdapterProfile"`
					InputAdapterProfile  string `xml:"inputAdapterProfile"`
					IncludeAppHdrs       string `xml:"includeAppHdrs"`
					SmmProfileExecution  string `xml:"smmProfileExecution"`
				} `xml:"messageManipulation"`
				RetryCounters struct {
					Invite  float64 `xml:"invite"`
					General float64 `xml:"general"`
				} `xml:"retryCounters"`
			} `xml:"signaling"`
			Services struct {
				SipArsProfile       string `xml:"sipArsProfile"`
				SipJipProfile       string `xml:"sipJipProfile"`
				TransparencyProfile string `xml:"transparencyProfile"`
				NatTraversal        struct {
					MediaNat string `xml:"mediaNat"`
				} `xml:"natTraversal"`
				NoRDIUpdateOn3XX           string `xml:"noRDIUpdateOn3XX"`
				BlockProgressOn3XXResponse string `xml:"blockProgressOn3XXResponse"`
			} `xml:"services"`
			Media struct {
				MediaIpInterfaceGroupName string `xml:"mediaIpInterfaceGroupName"`
				DirectMediaAllowed        string `xml:"directMediaAllowed"`
			} `xml:"media"`
			IngressIpPrefix []struct {
				IpAddress    string `xml:"ipAddress"`
				PrefixLength string `xml:"prefixLength"`
			} `xml:"ingressIpPrefix"`
			SipResponseCodeStats string `xml:"sipResponseCodeStats"`
		} `xml:"sipTrunkGroup"`
		TrunkGroupStatus []struct {
			Name                       string  `xml:"name"`
			State                      string  `xml:"state"`
			TotalCallsAvailable        float64 `xml:"totalCallsAvailable"`
			TotalCallsInboundReserved  float64 `xml:"totalCallsInboundReserved"`
			InboundCallsUsage          float64 `xml:"inboundCallsUsage"`
			OutboundCallsUsage         float64 `xml:"outboundCallsUsage"`
			TotalCallsConfigured       float64 `xml:"totalCallsConfigured"`
			PriorityCallUsage          float64 `xml:"priorityCallUsage"`
			TotalOutboundCallsReserved float64 `xml:"totalOutboundCallsReserved"`
			BwCurrentLimit             float64 `xml:"bwCurrentLimit"`
			BwAvailable                float64 `xml:"bwAvailable"`
			BwInboundUsage             float64 `xml:"bwInboundUsage"`
			BwOutboundUsage            float64 `xml:"bwOutboundUsage"`
			PacketOutDetectState       string  `xml:"packetOutDetectState"`
			PriorityBwUsage            float64 `xml:"priorityBwUsage"`
		} `xml:"trunkGroupStatus"`
		TrunkGroupQoeStatus []struct {
			Name                                        string  `xml:"name"`
			InboundRFactor                              float64 `xml:"inboundRFactor"`
			InboundRFactorFromSBXBOOT                   float64 `xml:"inboundRFactorFromSBXBOOT"`
			InboundRFactorNumCriticalThresholdBreached  float64 `xml:"inboundRFactorNumCriticalThresholdBreached"`
			InboundRFactorNumMajorThresholdBreached     float64 `xml:"inboundRFactorNumMajorThresholdBreached"`
			OutboundRFactor                             float64 `xml:"outboundRFactor"`
			OutboundRFactorFromSBXBOOT                  float64 `xml:"outboundRFactorFromSBXBOOT"`
			OutboundRFactorNumCriticalThresholdBreached float64 `xml:"outboundRFactorNumCriticalThresholdBreached"`
			OutboundRFactorNumMajorThresholdBreached    float64 `xml:"outboundRFactorNumMajorThresholdBreached"`
			CurrentASR                                  float64 `xml:"currentASR"`
			AsrFromSBXBOOT                              float64 `xml:"asrFromSBXBOOT"`
			AsrCriticalThresholdExceeded                float64 `xml:"asrCriticalThresholdExceeded"`
			AsrMajorThresholdExceeded                   float64 `xml:"asrMajorThresholdExceeded"`
			EgressSustainedCallRate                     float64 `xml:"egressSustainedCallRate"`
			EgressActiveCalls                           float64 `xml:"egressActiveCalls"`
			CurrentPgrd                                 float64 `xml:"currentPgrd"`
			QosDropCount                                float64 `xml:"qosDropCount"`
		} `xml:"trunkGroupQoeStatus"`
		SipRegAdaptiveNaptLearningStatistics []struct {
			Name                            string  `xml:"name"`
			SessionsInitiated               float64 `xml:"sessionsInitiated"`
			SessionsCompleted               float64 `xml:"sessionsCompleted"`
			SessionsCompletedDueToTimeout   float64 `xml:"sessionsCompletedDueToTimeout"`
			SessionsAbortedDueToTraffic     float64 `xml:"sessionsAbortedDueToTraffic"`
			SessionsReachedRelearnThreshold float64 `xml:"sessionsReachedRelearnThreshold"`
			SessionsInProgress              float64 `xml:"sessionsInProgress"`
			OptionsPolicerReject            float64 `xml:"optionsPolicerReject"`
			SessionAdmissionReject          float64 `xml:"sessionAdmissionReject"`
		} `xml:"sipRegAdaptiveNaptLearningStatistics"`
		SipTrunkgroupPortRangeStatistics []struct {
			Name                          string  `xml:"name"`
			PortRangeActivePorts          float64 `xml:"portRangeActivePorts"`
			PortRangeRegistrationFailures float64 `xml:"portRangeRegistrationFailures"`
		} `xml:"sipTrunkgroupPortRangeStatistics"`
		IpPeer []struct {
			Name                 string `xml:"name"`
			IpAddress            string `xml:"ipAddress"`
			SipResponseCodeStats string `xml:"sipResponseCodeStats"`
			Policy               struct {
				Description string `xml:"description"`
				Sip         struct {
					Fqdn     string `xml:"fqdn"`
					FqdnPort string `xml:"fqdnPort"`
				} `xml:"sip"`
			} `xml:"policy"`
			IpPort string `xml:"ipPort"`
		} `xml:"ipPeer"`
		PeerQosStatus []struct {
			Name                    string  `xml:"name"`
			EgressActiveCalls       float64 `xml:"egressActiveCalls"`
			EgressSustainedCallRate float64 `xml:"egressSustainedCallRate"`
			CurrentPGRD             float64 `xml:"currentPGRD"`
			CurrentASR              float64 `xml:"currentASR"`
			QosDropCount            float64 `xml:"qosDropCount"`
		} `xml:"peerQosStatus"`
		IpPeerCurrentStatistics []struct {
			Name                string  `xml:"name"`
			InboundSessions     float64 `xml:"inboundSessions"`
			InboundCPS          float64 `xml:"inboundCPS"`
			InboundMaxSessions  float64 `xml:"inboundMaxSessions"`
			OutboundSessions    float64 `xml:"outboundSessions"`
			OutboundCPS         float64 `xml:"outboundCPS"`
			OutboundMaxSessions float64 `xml:"outboundMaxSessions"`
		} `xml:"ipPeerCurrentStatistics"`
		IpPeerIntervalStatistics []struct {
			Number              string  `xml:"number"`
			Name                string  `xml:"name"`
			IntervalValid       string  `xml:"intervalValid"`
			Time                float64 `xml:"time"`
			InboundSessions     float64 `xml:"inboundSessions"`
			InboundCPS          float64 `xml:"inboundCPS"`
			InboundMaxSessions  float64 `xml:"inboundMaxSessions"`
			OutboundSessions    float64 `xml:"outboundSessions"`
			OutboundCPS         float64 `xml:"outboundCPS"`
			OutboundMaxSessions float64 `xml:"outboundMaxSessions"`
		} `xml:"ipPeerIntervalStatistics"`
		SipIpPeerResponseCurrentStatistics []struct {
			Name          string `xml:"name"`
			Direction     string `xml:"direction"`
			ResponseCode  string `xml:"responseCode"`
			ResponseCount string `xml:"responseCount"`
		} `xml:"sipIpPeerResponseCurrentStatistics"`
		SipIpPeerResponseIntervalStatistics []struct {
			Number        string  `xml:"number"`
			Name          string  `xml:"name"`
			Direction     string  `xml:"direction"`
			ResponseCode  string  `xml:"responseCode"`
			IntervalValid string  `xml:"intervalValid"`
			Time          string  `xml:"time"`
			ResponseCount float64 `xml:"responseCount"`
		} `xml:"sipIpPeerResponseIntervalStatistics"`
		SipArsStatus []struct {
			SigZoneId                   string `xml:"sigZoneId"`
			RecordIndex                 string `xml:"recordIndex"`
			SigPortNum                  string `xml:"sigPortNum"`
			EndpointDomainName          string `xml:"endpointDomainName"`
			EndpointIpAddress           string `xml:"endpointIpAddress"`
			EndpointIpPortNum           string `xml:"endpointIpPortNum"`
			EndpointArsState            string `xml:"endpointArsState"`
			EndpointStateTransitionTime string `xml:"endpointStateTransitionTime"`
		} `xml:"sipArsStatus"`
		SipTrunkGroupResponseCurrentStatistics []struct {
			Name          string  `xml:"name"`
			Direction     string  `xml:"direction"`
			ResponseCode  string  `xml:"responseCode"`
			ResponseCount float64 `xml:"responseCount"`
		} `xml:"sipTrunkGroupResponseCurrentStatistics"`
		SipTrunkGroupResponseIntervalStatistics []struct {
			Number        string  `xml:"number"`
			Name          string  `xml:"name"`
			Direction     string  `xml:"direction"`
			ResponseCode  string  `xml:"responseCode"`
			IntervalValid string  `xml:"intervalValid"`
			Time          string  `xml:"time"`
			ResponseCount float64 `xml:"responseCount"`
		} `xml:"sipTrunkGroupResponseIntervalStatistics"`
		TracerouteSigPort struct {
			State string `xml:"state"`
		} `xml:"tracerouteSigPort"`
		SipSigTlsSessionStatus []struct {
			Socket        string `xml:"socket"`
			Index         string `xml:"index"`
			PeerIpAddress string `xml:"peerIpAddress"`
			State         string `xml:"state"`
			Role          string `xml:"role"`
			Resumptions   string `xml:"resumptions"`
			TlsSessionId  string `xml:"tlsSessionId"`
		} `xml:"sipSigTlsSessionStatus"`
	} `xml:"zone"`
}

type processStructParams struct {
	Metrics    map[string]*prometheus.GaugeVec
	MetricName string
	Zone       string
	Context    string
	System     string
	logger     log.Logger
	tgName     string
}

// processStruct iterates through the fields in the structure `s`, if the field
// is a float64 it will set the value for the metric with the corresponding name
func processStruct(params processStructParams, s interface{}) {
	metrics := params.Metrics
	name := params.MetricName
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	sVal := reflect.ValueOf(s)
	if typ.Kind() != reflect.Struct {
		level.Warn(params.logger).Log("msg", fmt.Sprintf("%v: %v\n", typ.Kind().String(), name))
		return
	}

	tgName := sVal.FieldByName("Name")
	if tgName.IsValid() {
		if params.tgName == "" {
			params.tgName = tgName.String()
		} else {
			params.tgName = fmt.Sprintf("%s:%s", params.tgName, tgName.String())
		}
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		switch f.Type.Kind() {
		case reflect.Float64:
			metric, ok := metrics[fmt.Sprintf("%s_%s", name, f.Name)]
			if ok {
				metric.WithLabelValues(params.System, params.Context, params.Zone, params.tgName).Set(sVal.Field(i).Interface().(float64))
			} else {
				level.Warn(params.logger).Log("msg", fmt.Sprintf("Could not find metric %s_%s\n", name, f.Name))
			}
		case reflect.Struct:
			params.MetricName = fmt.Sprintf("%s_%s", name, f.Name)
			processStruct(params, sVal.Field(i).Interface())
		}
	}
}
