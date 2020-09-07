package exporter

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/prometheus/client_golang/prometheus"
)

var sleep = false

func wakeup() {
	lock := sync.RWMutex{}
	lock.Lock()
	time.Sleep(time.Minute * time.Duration(1))
	sleep = false
	lock.Unlock()
}

func QueryMaxConnection(instanceId string) float64 {
	client, err := rds.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	request := rds.CreateDescribeDBInstanceAttributeRequest()
	request.DBInstanceId = instanceId
	response, err := client.DescribeDBInstanceAttribute(request)
	if err != nil {
		panic(err)
	}
	maxConnections := response.Items.DBInstanceAttribute[0].MaxConnections
	return float64(maxConnections)
}

func (r *RdsExporter) NewClient() {
	client, err := cms.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	r.client = client
}

func (r *RdsExporter) GetInstance() {
	client, err := rds.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	request := rds.CreateDescribeDBInstancesRequest()
	request.PageSize = requests.NewInteger(100)
	response, err := client.DescribeDBInstances(request)
	if err != nil {
		panic(err)
	}
	instances := make(map[string]string)
	instancesType := make(map[string]string)
	maxConnections := make(map[string]float64)
	for _, v := range response.Items.DBInstance {
		maxConnection := QueryMaxConnection(v.DBInstanceId)
		maxConnections[v.DBInstanceId] = maxConnection
		if v.DBInstanceDescription != "" {
			instances[v.DBInstanceId] = v.DBInstanceDescription
		} else {
			instances[v.DBInstanceId] = v.DBInstanceId
		}
		instancesType[v.DBInstanceId] = v.DBInstanceType

	}
	r.instances = instances
	r.instancesType = instancesType
	r.maxConnections = maxConnections
}

func (r *RdsExporter) GetMetricMeta() {
	request := cms.CreateDescribeMetricMetaListRequest()
	request.Namespace = PROJECT
	request.PageSize = requests.NewInteger(100)
	response, err := r.client.DescribeMetricMetaList(request)
	if err != nil {
		panic(err)
	}
	for _, v := range response.Resources.Resource {
		if !NoDataMetricName[v.MetricName] {
			r.metricMeta = append(r.metricMeta, v.MetricName)
		}
	}
	r.metricMeta = append(r.metricMeta, MysqlTotalSessions)
}

func (r *RdsExporter) GetMetric(metricName string) {
	var dimensions []map[string]string
	for k := range r.instances {
		d := map[string]string{"instanceId": k}
		dimensions = append(dimensions, d)
	}
	dimension, err := json.Marshal(dimensions)
	if err != nil {
		log.Println(err)
	}
	request := cms.CreateDescribeMetricLastRequest()
	request.Namespace = PROJECT
	request.MetricName = metricName
	request.Dimensions = string(dimension)
	request.Period = "180"
	response, err := r.client.DescribeMetricLast(request)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal([]byte(response.Datapoints), &r.DataPoints)
	if err != nil {
		log.Println(err)
	}
}

func (r *RdsExporter) InitGauge() {
	r.NewClient()
	r.GetInstance()
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			r.GetInstance()
		}
	}()
	r.GetMetricMeta()
	r.metrics = map[string]*prometheus.GaugeVec{}
	for _, v := range r.metricMeta {
		r.metrics[v] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "aliyun_rds",
			Name:      strings.ToLower(v),
		}, []string{"instance_id", "instance_name"})
	}
}

func (r *RdsExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range r.metrics {
		v.Describe(ch)
	}
}

func (r *RdsExporter) Collect(ch chan<- prometheus.Metric) {
	if !sleep {
		sleep = true
		go wakeup()
		for _, v := range r.metricMeta {
			if v == MysqlTotalSessions {
				continue
			}
			r.GetMetric(v)
			for _, d := range r.DataPoints {
				if v == DataDelay {
					if r.instancesType[d.InstanceId] != ReadOnly {
						continue
					}
				}
				r.metrics[v].With(prometheus.Labels{
					"instance_id":   d.InstanceId,
					"instance_name": r.instances[d.InstanceId],
				}).Set(d.Average)
			}
			if v == ConnectionUsage {
				for _, d := range r.DataPoints {
					r.metrics[MysqlTotalSessions].With(prometheus.Labels{
						"instance_id":   d.InstanceId,
						"instance_name": r.instances[d.InstanceId],
					}).Set(d.Average * r.maxConnections[d.InstanceId] / 100)
				}
			}
			time.Sleep(34 * time.Millisecond)
		}
	}
	for _, m := range r.metrics {
		m.Collect(ch)
	}
}
