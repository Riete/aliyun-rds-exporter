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
	ConnectionUsage    string = "ConnectionUsage"
	ReadOnly           string = "Readonly"
	DataDelay          string = "DataDelay"
)

var NoDataMetricName = map[string]bool{
	"DockerCpuUsage":             true,
	"AvgLogSize":                 true,
	"Flow":                       true,
	"GroupCPUUtilization":        true,
	"GroupConnectionUtilization": true,
	"GroupDiskUtilization":       true,
	"GroupIOPSUtilization":       true,
	"MaxLogSize":                 true,
	"PG_RO_ReadLag":              true,
	"PG_RO_StreamingStatus":      true,
	"Rt":                         true,
	"SQLServer_NetworkInNew":     true,
	"SQLServer_NetworkOutNew":    true,
	"ServiceCurrentConnections":  true,
	"ServiceQueries":             true,
	"ServiceTotalConnections":    true,
	"TPS":                        true,
	"active_connections_per_cpu": true,
	"conn_usgae":                 true,
	"cpu_usage":                  true,
	"iops_usage":                 true,
	"local_fs_inode_usage":       true,
	"local_fs_size_usage":        true,
	"mem_usage":                  true,
}

type RdsExporter struct {
	client     *cms.Client
	DataPoints []struct {
		InstanceId string  `json:"instanceId"`
		Average    float64 `json:"Average"`
	}
	metrics        map[string]*prometheus.GaugeVec
	instances      map[string]string
	instancesType  map[string]string
	maxConnections map[string]float64
	metricMeta     []string
}
