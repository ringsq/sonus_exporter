package sonus

import (
	"context"

	"github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
)

type ServerInfo struct {
	ServerStatus []struct {
		Name                 string `xml:"name"`
		HwType               string `xml:"hwType"`
		SerialNum            string `xml:"serialNum"`
		PartNum              string `xml:"partNum"`
		PlatformVersion      string `xml:"platformVersion"`
		ApplicationVersion   string `xml:"applicationVersion"`
		MgmtRedundancyRole   string `xml:"mgmtRedundancyRole"`
		UpTime               string `xml:"upTime"`
		ApplicationUpTime    string `xml:"applicationUpTime"`
		LastRestartReason    string `xml:"lastRestartReason"`
		SyncStatus           string `xml:"syncStatus"`
		DaughterBoardPresent string `xml:"daughterBoardPresent"`
		CurrentTime          string `xml:"currentTime"`
		PktPortSpeed         string `xml:"pktPortSpeed"`
		ActualCeName         string `xml:"actualCeName"`
		HwSubType            string `xml:"hwSubType"`
		Fingerprint          string `xml:"fingerprint"`
	} `xml:"serverStatus"`
}

func ServerInfoMetrics(ctx context.Context, sbc *SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error {
	var (
		serverInfoVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_info",
			Help: "System Information",
		}, []string{"hwType", "serial", "server", "system", "version"})
	)
	serverInfo := &ServerInfo{}
	err := sbc.GetAndParse(ctx, serverInfo, serverInfoPath)
	if err != nil {
		return err
	}

	registry.MustRegister(serverInfoVec)

	for _, server := range serverInfo.ServerStatus {
		serverInfoVec.WithLabelValues(server.HwType, server.SerialNum, server.Name, sbc.System, server.ApplicationVersion).Set(1)
	}

	return nil
}
