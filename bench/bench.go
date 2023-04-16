package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"time"
)

var cores = [4]string{"askprogrammers", "askphotography", "youshouldknow", "fighters"}

var queries = [4][]string{
	{"c++", "java", "python", "javascript", "reddit", "stackoverflow", "programming", "html", "fork", "game", "network", "system", "connect", "pointer", "data"},
	{"fujifilm", "canon", "nikon", "film", "digital", "camera", "macro", "theme", "story", "macro", "blog", "portrait", "photographer"},
	{"cat", "dog", "flower", "boyfriend", "Windows", "free", "water", "coffee", "skill", "time", "youtube", "amazon", "desktop", "erase", "party", "school"},
	{"competitive", "fighting", "episode", "match", "community", "rules", "xbox", "mario", "game", "player", "PlayStation", "marvel", "dragon", "weapons", "Heroes"},
}

func printAllPossibleQueries() {
	for i := 0; i < len(cores); i++ {
		for j := 0; j < len(queries[i]); j++ {
			fmt.Printf("http://172.174.169.148:8983/solr/%s/select?q=text:%s&rows=1000\n", cores[i], queries[i][j])
		}
	}
}

var responseTimeTotal = [len(cores)]time.Duration{0, 0, 0}
var responseTimeCounter = [len(cores)]int{0, 0, 0}

var tr = &http.Transport{
	MaxIdleConns:        0,
	MaxIdleConnsPerHost: 10000,
	IdleConnTimeout:     30 * time.Second,
	DisableCompression:  true,
}

var client = &http.Client{Transport: tr}

func findArraySum(arr [len(cores)]int) int {
	res := 0
	for i := 0; i < len(arr); i++ {
		res += arr[i]
	}
	return res
}

func getURL() string {

	coreIndex := rand.Intn(len(cores))
	core := cores[coreIndex]
	query := queries[coreIndex][rand.Intn(len(queries[coreIndex]))]

	res := fmt.Sprintf("http://172.174.169.148:8983/solr/%s/select?q=text:%s&rows=1000", core, query)
	return res
}

func sendRequest(url string) {

	req, _ := http.NewRequest("GET", url, nil)

	var connect time.Time
	var duration time.Duration

	trace := &httptrace.ClientTrace{

		ConnectDone: func(network, addr string, err error) {
			connect = time.Now()
		},

		GotFirstResponseByte: func() {
			duration = time.Since(connect)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	response, err := client.Transport.RoundTrip(req)

	if err != nil {
		fmt.Println("Error while sending request:", err)
		return
	}

	defer response.Body.Close()

	fmt.Printf("Response status: %s\tduration: %v\n", response.Status, duration)
}

func main() {

	for {

		go sendRequest(getURL())

		// wait 1 second then repeat
		time.Sleep(10 * time.Millisecond)
	}
}
