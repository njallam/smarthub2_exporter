package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "smarthub2"

var volumeRegex = regexp.MustCompile(`<wan_conn_volume_list type="array" value="\[\['\d+%3B(\d+)%3B(\d+)'`)

type SmartHub2Collector struct {
	baseURL string

	client *http.Client
}

var metrics = map[string]*prometheus.Desc{
	"volume": prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "volume_bytes"), "Traffic volume in bytes", []string{"direction"}, nil),
}

func newSmartHub2Collector(ip string) *SmartHub2Collector {
	return &SmartHub2Collector{
		baseURL: "http://" + ip,
		client:  &http.Client{},
	}
}

func (collector *SmartHub2Collector) Collect(ch chan<- prometheus.Metric) {
	response, err := collector.client.Get(collector.baseURL + "/nonAuth/wan_conn.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	body := string(data)

	downstream, err := strconv.Atoi(volumeRegex.FindStringSubmatch(body)[1])
	if err != nil {
		log.Fatal(err)
	}
	upstream, err := strconv.Atoi(volumeRegex.FindStringSubmatch(body)[2])
	if err != nil {
		log.Fatal(err)
	}

	ch <- prometheus.MustNewConstMetric(metrics["volume"], prometheus.CounterValue, float64(downstream), "downstream")
	ch <- prometheus.MustNewConstMetric(metrics["volume"], prometheus.CounterValue, float64(upstream), "upstream")
}

func (collector *SmartHub2Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range metrics {
		ch <- desc
	}
}
