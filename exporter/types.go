package exporter

import (
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	accessKeyId     = os.Getenv("ACCESS_KEY_ID")
	accessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
	regionId        = os.Getenv("REGION_ID")
)

const (
	PROJECT            string = "acs_rds_dashboard"
	MysqlTotalSessions string = "MySQL_TotalSessions"
	ConnectionUsage           = "ConnectionUsage"
)

type RdsExporter struct {
	client     *cms.Client
	DataPoints []struct {
		InstanceId string  `json:"instanceId"`
		Average    float64 `json:"Average"`
	}
	metrics        map[string]*prometheus.GaugeVec
	instances      map[string]string
	maxConnections map[string]float64
	metricMeta     []string
}
