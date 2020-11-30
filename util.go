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

func writePath(load int, jaeger_sleep int, target_service string) {
	startTime := time.Now()
	SendRequest(RequestUrl, []byte(RequestBody), load, RequestBearer)
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
	f, err := os.OpenFile(filepath.Join("data", fmt.Sprintf("path.json")), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := f.Write(result); err != nil {
		fmt.Println(err)
	}
}

func predict() {
	exec.Command("python", "model.py").Run()
	time.Sleep(time.Duration(3) * time.Second)
	data, _ := ioutil.ReadFile("output.json")
	var m map[string]bool
	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(m)
}
