package client

import (
	"aliyun-sls-exporter/pkg/config"
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	sls20201230 "github.com/alibabacloud-go/sls-20201230/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"sort"
	"strconv"
	"time"
)

var ignores = map[string]struct{}{
	"timestamp": {},
	"Maximum":   {},
	"Minimum":   {},
	"Average":   {},
}

// Datapoint datapoint
type Datapoint map[string]interface{}

// Get return value for measure
func (d Datapoint) Get(pv string) float64 {
	v, ok := d[pv]
	if !ok {
		return 0
	}

	distFloat, err := strconv.ParseFloat(v.(string), 64)
	if err != nil {
		return 0
	}
	return distFloat
}

// Labels return labels that not in ignores
func (d Datapoint) Labels() []string {
	labels := make([]string, 0)
	for k := range d {
		if _, ok := ignores[k]; !ok {
			labels = append(labels, k)
		}
	}
	sort.Strings(labels)
	return labels
}

// Values return values for lables
func (d Datapoint) Values(labels ...string) []string {
	values := make([]string, 0, len(labels))
	for i := range labels {
		values = append(values, fmt.Sprintf("%s", d[labels[i]]))
	}
	return values
}

// MetricClient wrap cms.client
type MetricClient struct {
	cloudID string
	sls     *sls20201230.Client
	logger  log.Logger
}

// CreateClient create sls client
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *sls20201230.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("cn-zhangjiakou.log.aliyuncs.com")
	_result = &sls20201230.Client{}
	_result, _err = sls20201230.NewClient(config)
	return _result, _err
}

// NewMetricClient create metric Client
func NewMetricClient(cloudID, ak, secret string, logger log.Logger) (*MetricClient, error) {
	sls, err := CreateClient(tea.String(ak), tea.String(secret))

	if err != nil {
		return nil, err
	}
	//cmsClient.SetTransport(rt)
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &MetricClient{cloudID, sls, logger}, nil
}

func (c *MetricClient) SetTransport(rate int) {
	//rt := ratelimit.New(rate)
	//c.sls.SetTransport(rt)
}

// retrive get datapoints for metric
func (c *MetricClient) retrive(project, logstore, query string) ([]Datapoint, error) {
	timeAgo := c.getOneMinuteTime()
	var datapoints []Datapoint
	pageNumber := 1
	for {

		getLogsRequest := &sls20201230.GetLogsRequest{
			Type:  tea.String("log"),
			From:  tea.Int64(timeAgo[0]),
			To:    tea.Int64(timeAgo[1]),
			Query: tea.String(query + getPageLimit(pageNumber)),
		}
		runtime := &util.RuntimeOptions{}
		headers := make(map[string]*string)
		resp, err := c.sls.GetLogsWithOptions(tea.String(project), tea.String(logstore), getLogsRequest, headers, runtime)
		if err != nil {
			return nil, err
		}

		b, _ := json.Marshal(&resp.Body)
		var dp []Datapoint
		if err = json.Unmarshal(b, &dp); err != nil {
			// some execpected error
			return nil, err
		}
		datapoints = append(datapoints, dp...)
		if len(dp) < 100 {
			break
		}
		pageNumber++

	}

	return datapoints, nil
}

// Collect do collect metrics into channel
func (c *MetricClient) Collect(namespace string, sub string, l *config.LogsMetric, ch chan<- prometheus.Metric) {
	if l.Project == "" {
		level.Warn(c.logger).Log("msg", "metric name must been set")
		return
	}

	datapoints, err := c.retrive(l.Project, l.Logstore, l.Query)
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to retrive datapoints", "cloudID", c.cloudID, "namespace", sub, "name", l.Project, "error", err)
		return
	}

	for _, dp := range datapoints {
		val := dp.Get(l.Measure)
		ch <- prometheus.MustNewConstMetric(
			l.Desc(namespace, sub, dp.Labels()...),
			prometheus.GaugeValue,
			val,
			append(dp.Values(l.Dimensions...), c.cloudID)...,
		)
	}
}

// getOneMinuteTime Get the time before one minute, rounded
func (c *MetricClient) getOneMinuteTime() []int64 {

	var d []int64
	t1 := time.Now().Add(-time.Minute * 1)
	u1 := t1.Truncate(time.Minute * 1).Unix()
	d = append(d, u1)

	t2 := time.Now()
	u2 := t2.Truncate(time.Minute * 1).Unix()
	d = append(d, u2)
	level.Info(c.logger).Log("from:", u1, "to:", u2)

	return d
}

// getPageLimit Get Request Page Limit Value
func getPageLimit(pageNumber int) string {
	start, end := 0, 100
	step := (pageNumber - 1) * 100
	start = step + start
	end = step + end
	return " limit " + tea.ToString(start) + "," + tea.ToString(end)
}
