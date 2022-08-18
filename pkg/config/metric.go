package config

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// LogsMetric meta
type LogsMetric struct {
	Project     string   `json:"project"`
	Logstore    string   `json:"logstore"`
	Name        string   `json:"name"`
	Description string   `json:"desc,omitempty"`
	Dimensions  []string `json:"dimensions,omitempty"`
	Query       string   `json:"query"`
	Measure     string   `json:"measure,omitempty"`
	desc        *prometheus.Desc
}

// Desc to prometheus desc
func (l *LogsMetric) Desc(ns, sub string, dimensions ...string) *prometheus.Desc {
	if len(l.Dimensions) == 0 {
		l.Dimensions = dimensions
	}
	if l.desc == nil {
		l.desc = prometheus.NewDesc(
			strings.Join([]string{ns, sub, l.Name}, "_"),
			l.Description,
			append(l.Dimensions, "cloudID"),
			nil,
		)
	}
	return l.desc
}

// setDefaults options
func (l *LogsMetric) setDefaults() {
	if l.Description == "" {
		l.Description = l.Name
	}

	l.Description = fmt.Sprintf("%s measure:%s", l.Description, l.Measure)
}

var durationUnitMapping = map[string]string{
	"s": "second",
	"m": "minute",
	"h": "hour",
	"d": "day",
}

func formatUnit(s string) string {
	s = strings.TrimSpace(s)
	if s == "%" {
		return "percent"
	}

	if strings.IndexAny(s, "/") > 0 {
		fields := strings.Split(s, "/")
		if len(fields) == 2 {
			if v, ok := durationUnitMapping[fields[1]]; ok {
				return strings.ToLower(strings.Join([]string{fields[0], "per", v}, "_"))
			}
		}
	}
	return strings.ToLower(s)
}
