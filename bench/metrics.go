package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Core struct {
	name     string
	shards   int
	replicas int
}

type SolrMetricsResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	Metrics map[string]struct {
		QueryRequestTimes struct {
			Count          int     `json:"count"`
			MeanRate       float64 `json:"meanRate"`
			OneMinRate     float64 `json:"1minRate"`
			FiveMinRate    float64 `json:"5minRate"`
			FifteenMinRate float64 `json:"15minRate"`
			MinMs          float64 `json:"min_ms"`
			MaxMs          float64 `json:"max_ms"`
			MeanMs         float64 `json:"mean_ms"`
			MedianMs       float64 `json:"median_ms"`
			StdDevMs       float64 `json:"stddev_ms"`
			P75Ms          float64 `json:"p75_ms"`
			P95Ms          float64 `json:"p95_ms"`
			P99Ms          float64 `json:"p99_ms"`
			P999Ms         float64 `json:"p999_ms"`
		} `json:"QUERY./select.requestTimes"`
	} `json:"metrics"`
}

var cores = []Core{
	{name: "askprogrammers", shards: 1, replicas: 1},
	{name: "askphotography", shards: 1, replicas: 1},
	{name: "youshouldknow", shards: 1, replicas: 1},
	{name: "fighters", shards: 1, replicas: 1},
	{name: "fighters", shards: 1, replicas: 2},
	{name: "fighters", shards: 2, replicas: 4},
	{name: "fighters", shards: 2, replicas: 6},
}

func getMetrics() {
	response, err := http.Get("http://172.174.169.148:8983/solr/admin/metrics?group=core&prefix=QUERY./select.requestTimes")

	if err != nil {
		fmt.Println("Error while sending request:", err)
		return
	}

	b, _ := ioutil.ReadAll(response.Body)

	var metricResponse SolrMetricsResponse
	err = json.Unmarshal(b, &metricResponse)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(cores); i++ {
		fmt.Printf("Name: %s\tshard: %d\treplica: %d\n", cores[i].name, cores[i].shards, cores[i].replicas)
		fmt.Printf("1minRate: %f\tp95: %f\tp99: %f\tp999: %f\n", metricResponse.Metrics[fmt.Sprintf("solr.core.%s.shard%d.replica_n%d", cores[i].name, cores[i].shards, cores[i].replicas)].QueryRequestTimes.OneMinRate, metricResponse.Metrics[fmt.Sprintf("solr.core.%s.shard%d.replica_n%d", cores[i].name, cores[i].shards, cores[i].replicas)].QueryRequestTimes.P95Ms, metricResponse.Metrics[fmt.Sprintf("solr.core.%s.shard%d.replica_n%d", cores[i].name, cores[i].shards, cores[i].replicas)].QueryRequestTimes.P99Ms, metricResponse.Metrics[fmt.Sprintf("solr.core.%s.shard%d.replica_n%d", cores[i].name, cores[i].shards, cores[i].replicas)].QueryRequestTimes.P999Ms)
	}

	fmt.Println()

	defer response.Body.Close()
}

func main() {
	for {
		getMetrics()

		time.Sleep(2 * time.Second)
	}
}
