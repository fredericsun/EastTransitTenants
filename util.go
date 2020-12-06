package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

type Pair struct {
	Service string
	Value   int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func writePath(requestUrl string, requestBody string, requestBearer string, load int, jaeger_sleep int, target_service string) {
	startTime := time.Now()
	SendRequest(requestUrl, []byte(requestBody), load, requestBearer)
	time.Sleep(time.Duration(jaeger_sleep) * time.Second)

	jc := NewJaegerClient(JaegerIP)
	traces, err := jc.QueryTraces(target_service, "", startTime, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data []TestServiceMetric
	for _, trace := range traces {
		serMetric := ConstructTestingData(trace, load)
		data = append(data, serMetric...)
	}

	result, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.OpenFile(filepath.Join("data", "test", fmt.Sprintf("path.json")), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := f.Write(result); err != nil {
		fmt.Println(err)
	}
}

func train() {
	exec.Command("python", "train.py").Run()
	time.Sleep(time.Duration(3) * time.Second)
}

func predict() map[string]bool {
	exec.Command("python", "predict.py").Run()
	data, _ := ioutil.ReadFile("output.json")
	var m map[string]bool

	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Println("Error: ", err)
	}
	return m
}

func bottleneckChanged(request RequestData, service string, latency_ms int, jc *JaegerClient, batch int) bool {
	jeagerSleep := 3 * time.Second
	channel := make(chan bool)

	virtualSpeedUp(service, latency_ms)
	startTime := time.Now()
	SendRequest(request.url, request.body, batch, request.bearer)
	go removeVirtualSpeedUp(channel)
	time.Sleep(jeagerSleep)
	traces, err := jc.QueryTraces("ts-travel2-service", "", startTime, 0)
	if err != nil {
		fmt.Println(err)
		<-channel
		return false
	}
	count := 0
	for _, trace := range traces {
		bottleneck := findCurrentBottleneck(trace, service, latency_ms)
		fmt.Printf("found bottleneck: %s\n", bottleneck)
		if bottleneck == service {
			count++
		}
	}
	<-channel
	return count <= batch/2
}

func lookForBoundary(request RequestData, service string, precision_ms int, load int) int {
	jc := NewJaegerClient(JaegerIP)
	lo := 0
	hi := precision_ms
	for hi < 30000 {
		if bottleneckChanged(request, service, hi, jc, load) {
			break
		}
		lo = hi
		hi *= 2
	}
	if hi >= 30000 {
		panic(fmt.Errorf("Testing Hi of %d to change bottleneck, something is wrong", hi))
	}
	for hi-lo > precision_ms {
		mid := (hi + lo) / 2
		if !bottleneckChanged(request, service, mid, jc, load) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return (hi + lo) / 2
}

func generate(url string, body []byte, bearer string, filename string, iteration int, target_service string, jaeger_sleep int) {
	fmt.Printf("%s: Sending %d request at one time for %d iterations\n", filename, 100/iteration, iteration)

	var trainData []ServiceMetric
	for i := 0; i < 100/iteration; i++ {
		/*************************
		Generate requests
		*************************/
		startTime := time.Now()
		SendRequest(url, []byte(body), iteration, bearer)
		fmt.Println("Request generation completed")
		time.Sleep(time.Duration(jaeger_sleep) * time.Second)

		/*************************
		Get tracing information
		*************************/
		// Create jaeger service client
		jc := NewJaegerClient(JaegerIP)
		// Query traces
		traces, err := jc.QueryTraces(target_service, "", startTime, 0)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Construct training data
		for _, trace := range traces {
			critcalPath := CleanupTrace(trace)
			serMetric := ConstructTrainingData(trace, iteration, critcalPath)
			trainData = append(trainData, serMetric...)
		}
	}
	result, err := json.Marshal(trainData)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.OpenFile(filepath.Join("data", fmt.Sprintf("%s_%d.json", filename, iteration)), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := f.Write(result); err != nil {
		fmt.Println(err)
	}
}

func generateTrainingData(requests []TrainRequestData, bearer string, workload []int, target_service string) {
	for _, requestData := range requests {
		for _, iteration := range workload {
			generate(requestData.url, []byte(requestData.body), bearer, requestData.name, iteration, target_service, jaeger_sleep)
		}
	}
}
