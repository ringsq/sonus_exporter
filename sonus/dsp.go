package sonus

import (
	"context"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ringsq/sonus_exporter/config"
)

/*
<collection xmlns:y="http://tail-f.com/ns/rest">
  <dspUsage xmlns="http://sonusnet.com/ns/mibs/SONUS-DRM-DSPSTATUS/1.0">
    <systemName>densbc01</systemName>
    <slot1ResourcesUtilized>68</slot1ResourcesUtilized>
    <slot2ResourcesUtilized>0</slot2ResourcesUtilized>
    ...
  </dspUsage>
</collection>
*/

type dspUsageCollection struct {
	DSPUsage *dspUsage `xml:"http://sonusnet.com/ns/mibs/SONUS-DRM-DSPSTATUS/1.0 dspUsage"`
}

type dspUsage struct {
	SystemName                         string  `xml:"systemName"`
	Slot1ResourcesUtilized             float64 `xml:"slot1ResourcesUtilized"`
	Slot2ResourcesUtilized             float64 `xml:"slot2ResourcesUtilized"`
	Slot3ResourcesUtilized             float64 `xml:"slot3ResourcesUtilized"`
	Slot4ResourcesUtilized             float64 `xml:"slot4ResourcesUtilized"`
	CompressionTotal                   float64 `xml:"compressionTotal"`
	CompressionAvailable               float64 `xml:"compressionAvailable"`
	CompressionUtilization             float64 `xml:"compressionUtilization"`
	CompressionHighPriorityUtilization float64 `xml:"compressionHighPriorityUtilization"`
	CompressionAllocFailures           float64 `xml:"compressionAllocFailures"`
	G711Total                          float64 `xml:"g711Total"`
	G711Utilization                    float64 `xml:"g711Utilization"`
	G711SsTotal                        float64 `xml:"g711SsTotal"`
	G711SsUtilization                  float64 `xml:"g711SsUtilization"`
	G726Total                          float64 `xml:"g726Total"`
	G726Utilization                    float64 `xml:"g726Utilization"`
	G7231Total                         float64 `xml:"g7231Total"`
	G7231Utilization                   float64 `xml:"g7231Utilization"`
	G722Total                          float64 `xml:"g722Total"`
	G722Utilization                    float64 `xml:"g722Utilization"`
	G7221Total                         float64 `xml:"g7221Total"`
	G7221Utilization                   float64 `xml:"g7221Utilization"`
	G729AbTotal                        float64 `xml:"g729AbTotal"`
	G729AbUtilization                  float64 `xml:"g729AbUtilization"`
	EcmTotal                           float64 `xml:"ecmTotal"`
	EcmUtilization                     float64 `xml:"ecmUtilization"`
	IlbcTotal                          float64 `xml:"ilbcTotal"`
	IlbcUtilization                    float64 `xml:"ilbcUtilization"`
	AmrNbTotal                         float64 `xml:"amrNbTotal"`
	AmrNbUtilization                   float64 `xml:"amrNbUtilization"`
	AmrNbT140Total                     float64 `xml:"amrNbT140Total"`
	AmrNbT140Utilization               float64 `xml:"amrNbT140Utilization"`
	AmrWbTotal                         float64 `xml:"amrWbTotal"`
	AmrWbUtilization                   float64 `xml:"amrWbUtilization"`
	AmrWbT140Total                     float64 `xml:"amrWbT140Total"`
	AmrWbT140Utilization               float64 `xml:"amrWbT140Utilization"`
	Evrcb0Total                        float64 `xml:"evrcb0Total"`
	Evrcb0Utilization                  float64 `xml:"evrcb0Utilization"`
	Evrc0Total                         float64 `xml:"evrc0Total"`
	Evrc0Utilization                   float64 `xml:"evrc0Utilization"`
	ToneTotal                          float64 `xml:"toneTotal"`
	ToneAvailable                      float64 `xml:"toneAvailable"`
	ToneUtilization                    float64 `xml:"toneUtilization"`
	ToneHighPriorityUtilization        float64 `xml:"toneHighPriorityUtilization"`
	ToneAllocFailures                  float64 `xml:"toneAllocFailures"`
	EfrTotal                           float64 `xml:"efrTotal"`
	EfrUtilization                     float64 `xml:"efrUtilization"`
	G711V8Total                        float64 `xml:"g711V8Total"`
	G711V8Utilization                  float64 `xml:"g711V8Utilization"`
	G711SsV8Total                      float64 `xml:"g711SsV8Total"`
	G711SsV8Utilization                float64 `xml:"g711SsV8Utilization"`
	G726V8Total                        float64 `xml:"g726V8Total"`
	G726V8Utilization                  float64 `xml:"g726V8Utilization"`
	G7231V8Total                       float64 `xml:"g7231V8Total"`
	G7231V8Utilization                 float64 `xml:"g7231V8Utilization"`
	G722V8Total                        float64 `xml:"g722V8Total"`
	G722V8Utilization                  float64 `xml:"g722V8Utilization"`
	G7221V8Total                       float64 `xml:"g7221V8Total"`
	G7221V8Utilization                 float64 `xml:"g7221V8Utilization"`
	G729AbV8Total                      float64 `xml:"g729AbV8Total"`
	G729AbV8Utilization                float64 `xml:"g729AbV8Utilization"`
	EcmV34Total                        float64 `xml:"ecmV34Total"`
	EcmV34Utilization                  float64 `xml:"ecmV34Utilization"`
	IlbcV8Total                        float64 `xml:"ilbcV8Total"`
	IlbcV8Utilization                  float64 `xml:"ilbcV8Utilization"`
	OpusTotal                          float64 `xml:"opusTotal"`
	OpusUtilization                    float64 `xml:"opusUtilization"`
	EvsTotal                           float64 `xml:"evsTotal"`
	EvsUtilization                     float64 `xml:"evsUtilization"`
	Silk8Total                         float64 `xml:"silk8Total"`
	Silk8Utilization                   float64 `xml:"silk8Utilization"`
	Silk16Total                        float64 `xml:"silk16Total"`
	Silk16Utilization                  float64 `xml:"silk16Utilization"`
}

func DSPMetrics(ctx context.Context, sbc *SBC, cfg *config.Config, registry *prometheus.Registry, logger log.Logger) error {
	var (
		DSP_Resources_Used = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_dsp_resources_used",
			Help: "Usage of DSP resources per slot",
		}, []string{"system", "slot"})

		DSP_Resources_Total = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_dsp_resources_total",
			Help: "Total compression resources",
		}, []string{"system"})
		DSP_Compression_Utilization = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_dsp_compression_utilization",
			Help: "Compression resource utilization, in percent",
		}, []string{"system"})

		DSP_Codec_Utilization = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "sonus_dsp_codec_utilization",
			Help: "Codec utilization, in percent",
		}, []string{"system", "codec"})

		dsp = new(dspUsageCollection)
	)

	registry.MustRegister(DSP_Resources_Total)
	registry.MustRegister(DSP_Resources_Used)
	registry.MustRegister(DSP_Codec_Utilization)
	registry.MustRegister(DSP_Compression_Utilization)

	err := sbc.GetAndParse(ctx, dsp, dspStatusPath)
	if err != nil {
		return err
	}

	d := dsp.DSPUsage

	DSP_Resources_Used.WithLabelValues(sbc.System, "1").Set(d.Slot1ResourcesUtilized)
	DSP_Resources_Used.WithLabelValues(sbc.System, "2").Set(d.Slot2ResourcesUtilized)
	DSP_Resources_Used.WithLabelValues(sbc.System, "3").Set(d.Slot3ResourcesUtilized)
	DSP_Resources_Used.WithLabelValues(sbc.System, "4").Set(d.Slot4ResourcesUtilized)

	DSP_Resources_Total.WithLabelValues(sbc.System).Set(d.CompressionTotal)
	DSP_Compression_Utilization.WithLabelValues(sbc.System).Set(d.CompressionUtilization)

	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.711").Set(d.G711Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.711 Silence Suppression").Set(d.G711SsUtilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.726").Set(d.G726Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.723.1").Set(d.G7231Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.722").Set(d.G722Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.722.1").Set(d.G7221Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.729").Set(d.G729AbUtilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "ECM").Set(d.EcmUtilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "iLBC").Set(d.IlbcUtilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "AMR-NB").Set(d.AmrNbUtilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "AMR-WB").Set(d.AmrWbUtilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "Tone").Set(d.ToneUtilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.711 V8").Set(d.G711V8Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.711 Silence Suppression V8").Set(d.G711SsV8Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.726 V8").Set(d.G726V8Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.723.1 V8").Set(d.G7231V8Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.722 V8").Set(d.G722V8Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.722.1 V8").Set(d.G7221V8Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "G.729 V8").Set(d.G729AbV8Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "ECM V.34").Set(d.EcmV34Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "iLBC V8").Set(d.IlbcV8Utilization)
	DSP_Codec_Utilization.WithLabelValues(sbc.System, "Opus").Set(d.OpusUtilization)

	return nil
}
