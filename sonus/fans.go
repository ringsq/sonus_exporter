package sonus

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
)

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <fanStatus xmlns="http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0">
    <serverName>densbc01a</serverName>
    <fanId>FAN1/BOT</fanId>
    <speed>5632 RPM</speed>
  </fanStatus>
...
</collection>
*/

type fanCollection struct {
	FanStatus []*fanStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 fanStatus,omitempty"`
}

type fanStatus struct {
	ServerName string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 serverName"`
	FanID      string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 fanId"`
	Speed      string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 speed"`
}

func (f fanStatus) speedToRPM() (float64, error) {
	var rpm = strings.TrimSuffix(f.Speed, " RPM")
	return strconv.ParseFloat(rpm, 64)
}

func FanMetrics(ctx context.Context, sbc *SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error {
	var (
		Fan_Speed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_fan_speed",
			Help: "Current speed of fans, in RPM",
		}, []string{"system", "server", "fanID"})
		fans = new(fanCollection)
	)

	registry.MustRegister(Fan_Speed)

	err := sbc.GetAndParse(ctx, fans, fanStatusPath)
	if err != nil {
		return err
	}
	for _, fan := range fans.FanStatus {
		var fanRpm, err = fanStatus.speedToRPM(*fan)
		if err != nil {
			level.Error(logger).Log("msg", fmt.Sprintf("Failed to convert fan speed (%q) to rpm", fan.Speed), "err", err)
			continue
		}
		Fan_Speed.WithLabelValues(sbc.System, fan.ServerName, fan.FanID).Set(fanRpm)
	}

	return nil
}
