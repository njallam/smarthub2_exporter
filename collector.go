package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "smarthub2"

var volumeRegex = regexp.MustCompile(`<wan_conn_volume_list type="array" value="\[\['\d+%3B(\d+)%3B(\d+)'`)
var deviceVolumeRegex = regexp.MustCompile(`{timestamp:'\d+',app:'\d+',mac:'([\da-f%3A]+)',tx:'(\d+)',rx:'(\d+)'`)

type SmartHub2Collector struct {
	baseURL string

	client *http.Client
}

var metrics = map[string]*prometheus.Desc{
	"totalTraffic":  prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "total_bytes"), "Total traffic in bytes", []string{"direction"}, nil),
	"deviceTraffic": prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "device_bytes"), "Device traffic in bytes", []string{"mac", "direction"}, nil),
}

func newSmartHub2Collector(ip string) *SmartHub2Collector {
	return &SmartHub2Collector{
		baseURL: "http://" + ip,
		client:  &http.Client{},
	}
}

func (collector *SmartHub2Collector) Fetch(path string) string {
	response, err := collector.client.Get(collector.baseURL + "/" + path)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

func (collector *SmartHub2Collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		body := collector.Fetch("nonAuth/wan_conn.xml")
		volumeMatches := volumeRegex.FindStringSubmatch(body)
		if len(volumeMatches) != 3 {
			log.Fatal("Volume regex did not match")
		}
		rx, _ := strconv.Atoi(volumeMatches[1])
		tx, _ := strconv.Atoi(volumeMatches[2])
		ch <- prometheus.MustNewConstMetric(metrics["totalTraffic"], prometheus.CounterValue, float64(tx), "upload")
		ch <- prometheus.MustNewConstMetric(metrics["totalTraffic"], prometheus.CounterValue, float64(rx), "download")
		wg.Done()
	}()

	go func() {
		body := collector.Fetch("cgi/cgi_basicMyDevice.js")
		matches := deviceVolumeRegex.FindAllStringSubmatch(body, -1)
		for _, match := range matches {
			mac := strings.ToLower(strings.Replace(match[1], "%3A", ":", -1))
			tx, _ := strconv.Atoi(match[2])
			rx, _ := strconv.Atoi(match[3])
			ch <- prometheus.MustNewConstMetric(metrics["deviceTraffic"], prometheus.CounterValue, float64(tx), mac, "upload")
			ch <- prometheus.MustNewConstMetric(metrics["deviceTraffic"], prometheus.CounterValue, float64(rx), mac, "download")
		}
		wg.Done()
	}()

	wg.Wait()
}

func (collector *SmartHub2Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range metrics {
		ch <- desc
	}
}
