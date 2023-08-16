package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/common/version"
)

// NewRegistry returns new prometheus registry with registered Exporter and common exporters.
func NewRegistry(exp *Exporter) prometheus.Gatherer {
	reg := prometheus.NewRegistry()
	reg.MustRegister(version.NewCollector(Name))
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
		Namespace:    "",
		ReportErrors: false,
	}))
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(exp)
	return reg
}
