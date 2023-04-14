package sonus

import (
	"context"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
)

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <powerSupplyStatus xmlns="http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0">
    <serverName>densbc01a</serverName>
    <powerSupplyId>PSA</powerSupplyId>
    <present>true</present>
    <productName>TECTROL  TC92S-1525R</productName>
    <serialNum>00000000</serialNum>
    <partNum>TC92S-1525R</partNum>
    <powerFault>false</powerFault>
    <voltageFault>false</voltageFault>
  </powerSupplyStatus>
...
</collection>
*/

type powerSupplyCollection struct {
	PowerSupplyStatus []*powerSupplyStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 powerSupplyStatus,omitempty"`
}

type powerSupplyStatus struct {
	ServerName    string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 serverName"`
	PowerSupplyID string `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 powerSupplyId"`
	Present       bool   `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 present"`
	PowerFault    bool   `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 powerFault"`
	VoltageFault  bool   `xml:"http://sonusnet.com/ns/mibs/SONUS-SYSTEM-MIB/1.0 voltageFault"`
}

func boolToMetric(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func PowerMetrics(ctx context.Context, sbc *SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error {
	var (
		PowerSupply_Power_Fault = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_powersupply_powerfault",
			Help: "Is there a power fault, per supply",
		}, []string{"system", "server", "powerSupplyID"})

		PowerSupply_Voltage_Fault = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_powersupply_voltagefault",
			Help: "Is there a voltage fault, per supply",
		}, []string{"system", "server", "powerSupplyID"})
		PowerSupply_Present = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_powersupply_present",
			Help: "Indicates if the powersupply is installed",
		}, []string{"system", "server", "powerSupplyID"})

		powerSupplies = new(powerSupplyCollection)
	)

	err := sbc.GetAndParse(ctx, powerSupplies, powerSupplyPath)
	if err != nil {
		return err
	}
	registry.MustRegister(PowerSupply_Power_Fault)
	registry.MustRegister(PowerSupply_Voltage_Fault)
	registry.MustRegister(PowerSupply_Present)

	for _, psu := range powerSupplies.PowerSupplyStatus {
		PowerSupply_Power_Fault.WithLabelValues(sbc.System, psu.ServerName, psu.PowerSupplyID).Set(boolToMetric(psu.PowerFault))
		PowerSupply_Voltage_Fault.WithLabelValues(sbc.System, psu.ServerName, psu.PowerSupplyID).Set(boolToMetric(psu.VoltageFault))
		PowerSupply_Present.WithLabelValues(sbc.System, psu.ServerName, psu.PowerSupplyID).Set(boolToMetric(psu.Present))
	}
	return nil
}
