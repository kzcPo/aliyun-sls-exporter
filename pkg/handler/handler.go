package handler

import (
	"aliyun-sls-exporter/pkg/config"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"sigs.k8s.io/yaml"
)

// Handler http metrics handler
type Handler struct {
	logger log.Logger
	server *http.Server
}

// New create http handler
func New(addr string, logger log.Logger, rate int, cfg *config.Config, c map[string]prometheus.Collector) (*Handler, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	h := &Handler{
		logger: logger,
		server: &http.Server{
			Addr: net.JoinHostPort(host, port),
		},
	}
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		handlerMetrics(w, r, c)
	})
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		c, err := yaml.Marshal(cfg)
		if err != nil {
			level.Error(logger).Log("msg", "Error marshaling configuration", "err", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(c)
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Service is UP"))
	})
	return h, nil
}

func handlerMetrics(w http.ResponseWriter, r *http.Request, c map[string]prometheus.Collector) {
	registry := prometheus.NewRegistry()
	for cloudId, _ := range c {
		registry.MustRegister(c[cloudId])
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

// Run start server
func (h *Handler) Run() error {
	level.Info(h.logger).Log("msg", "Starting metric handler", "port", h.server.Addr)
	fmt.Println("msg", "Starting metric handler", "port", h.server.Addr)
	return h.server.ListenAndServe()
}
