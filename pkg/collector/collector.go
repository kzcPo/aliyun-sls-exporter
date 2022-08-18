package collector

import (
	"aliyun-sls-exporter/pkg/client"
	"aliyun-sls-exporter/pkg/config"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

const AppName = "slsmonitor"

// cloudMonitor ..
type cloudMonitor struct {
	namespace string
	cfg       *config.Config
	logger    log.Logger
	// sdk client
	client *client.MetricClient
	rate   int
	lock   sync.Mutex
}

// NewCloudMonitorCollector create a new collector for cloud monitor
func NewCloudMonitorCollector(appName string, cfg *config.Config, rate int, logger log.Logger) (map[string]prometheus.Collector, error) {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	cloudMonitors := make(map[string]prometheus.Collector)
	for cloudID, credential := range cfg.Credentials {
		cli, err := client.NewMetricClient(cloudID, credential.AccessKey, credential.AccessKeySecret, logger)
		if err != nil {
			continue
		}
		cloudMonitors[cloudID] = &cloudMonitor{
			namespace: appName,
			cfg:       cfg,
			logger:    logger,
			client:    cli,
			rate:      rate,
		}
	}
	return cloudMonitors, nil
}

// NewCloudMonitorCollectorFromURL create a new collector from HTTP Request URL for cloud monitor
func NewCloudMonitorCollectorFromURL(cli *client.MetricClient, cloudID string, cfg *config.Config, rate int, logger log.Logger) map[string]prometheus.Collector {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	collectors := make(map[string]prometheus.Collector)
	collectors[cloudID] = &cloudMonitor{
		namespace: AppName,
		cfg:       cfg,
		logger:    logger,
		client:    cli,
		rate:      rate,
	}
	return collectors
}

func (m *cloudMonitor) Describe(ch chan<- *prometheus.Desc) {
}

func (m *cloudMonitor) Collect(ch chan<- prometheus.Metric) {
	m.lock.Lock()
	defer m.lock.Unlock()

	wg := &sync.WaitGroup{}
	// do collect
	m.client.SetTransport(m.rate)
	for sub, req := range m.cfg.LogsMetric {
		for i := range req {
			wg.Add(1)
			go func(namespace string, metric *config.LogsMetric) {
				defer wg.Done()
				m.client.Collect(m.namespace, namespace, metric, ch)
			}(sub, req[i])
		}
	}
	wg.Wait()
}
