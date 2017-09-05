package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jbub/pgbouncer_exporter/collector"
	"github.com/jbub/pgbouncer_exporter/config"
	"github.com/jbub/pgbouncer_exporter/domain"
	"github.com/jbub/pgbouncer_exporter/store"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

func init() {
	prometheus.MustRegister(version.NewCollector(collector.Name))
}

func getLandingPage(telemetryPath string) []byte {
	return []byte(`
	<html>
	<head>
	<title>` + collector.Name + `</title>
	</head>
	<body>
	<h1>` + collector.Name + `</h1>
	<p><a href="` + telemetryPath + `">Metrics</a></p>
	</body>
	</html>`)
}

func New(cfg config.Config) (*HTTPServer, error) {
	st, err := store.NewSQLStore(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize store: %v", err)
	}

	return &HTTPServer{
		cfg: cfg,
		st:  st,
	}, nil
}

// HTTPServer represents prometheus exporter http server.
type HTTPServer struct {
	cfg config.Config
	st  domain.Store
}

// Run runs http server.
func (s *HTTPServer) Run() error {
	// cleanup on program interrupt/termination
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(st domain.Store) {
		<-c
		st.Close()
		os.Exit(1)
	}(s.st)

	exp := collector.New(s.cfg, s.st)
	prometheus.MustRegister(exp)

	log.Infoln("Starting ", collector.Name, version.Info())
	log.Infoln("Server listening on", s.cfg.ListenAddress)
	log.Infoln("Metrics available at", s.cfg.TelemetryPath)

	http.Handle(s.cfg.TelemetryPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write(getLandingPage(s.cfg.TelemetryPath))
	})
	return http.ListenAndServe(s.cfg.ListenAddress, nil)
}
