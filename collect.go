package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

type ServiceMetric struct {
	SelfPercent  float64
	TotalPercent float64
	CountPercent float64
	BottleNeck   bool
	TotalDur     int64
	TotalService int
	TotalSpan    int
	Load         int
}

func ConstructTrainingData(spans []MySpan, load int) {

	serMetric := make(map[string]ServiceMetric)

	serExc := make(map[string]time.Duration)
	serDur := make(map[string]time.Duration)
	serCount := make(map[string]int)
	earliestStartTime := spans[0].StartTime
	latestEndTime := earliestStartTime.Add(spans[0].Duration)

	spanToServ := make(map[string]string)
	for _, span := range spans {
		spanToServ[span.SpanID] = span.Process.GetServiceName()
		serCount[span.Process.GetServiceName()] += 1
		serDur[span.Process.GetServiceName()] += span.Duration
		if earliestStartTime.After(span.StartTime) {
			earliestStartTime = span.StartTime
		}
		if span.StartTime.Add(span.Duration).After(latestEndTime) {
			latestEndTime = span.StartTime.Add(span.Duration)
		}
	}

	for _, span := range spans {
		serExc[span.Process.GetServiceName()] += span.Duration
		for _, ref := range span.References {
			if ref["RefType"] == "CHILD_OF" {
				serExc[spanToServ[ref["SpanID"]]] -= span.Duration
			}
		}
	}
	max := time.Second * 0
	res := ""

	for serv, dur := range serExc {
		if dur > max {
			max = dur
			res = serv
		}
	}

	// request information
	totalSer := len(serExc)
	totalSpan := len(spans)
	totalDur := latestEndTime.UnixNano() - earliestStartTime.UnixNano()

	for serv, _ := range serExc {
		isBottleNeck := false
		if res == serv {
			isBottleNeck = true
		}
		serviceMetric := ServiceMetric{
			SelfPercent:  float64(serExc[serv].Nanoseconds()) / float64(serDur[serv].Nanoseconds()),
			TotalPercent: float64(serExc[serv].Nanoseconds()) / float64(totalDur),
			CountPercent: float64(serCount[serv]) / float64(totalSpan),
			BottleNeck:   isBottleNeck,
			TotalDur:     totalDur,
			TotalService: totalSer,
			TotalSpan:    totalSpan,
			Load:         load,
		}
		serMetric[serv] = serviceMetric
	}
	result, err := json.Marshal(serMetric)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := ioutil.WriteFile(filepath.Join("data", "training"), result, 0644); err != nil {
		fmt.Println(err)
		return
	}
}
