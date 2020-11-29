package main

import (
	"time"
)

type ServiceMetric struct {
	SelfPercent  float64
	ServiceDur   float64
	TotalPercent float64
	CountPercent float64
	TotalDur     int64
	TotalService int
	TotalSpan    int
	Load         int
	BottleNeck   bool
}

func ConstructTrainingData(spans []MySpan, load int) []ServiceMetric {

	var serMetric []ServiceMetric

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
	candidates := make(map[string]bool)

	for serv, dur := range serExc {
		if dur > max {
			max = dur
			res = serv
		}
	}
	candidates[res] = true
	for serv, dur := range serExc {
		if float64(dur.Nanoseconds()) >= float64(max.Nanoseconds())*0.7 {
			candidates[serv] = true
		}
	}

	// request information
	totalSer := len(serExc)
	totalSpan := len(spans)
	totalDur := latestEndTime.UnixNano() - earliestStartTime.UnixNano()

	for serv, _ := range serExc {
		isBottleNeck := false
		if _, ok := candidates[serv]; ok {
			isBottleNeck = true
		}
		serviceMetric := ServiceMetric{
			SelfPercent:  float64(serExc[serv].Nanoseconds()) / float64(serDur[serv].Nanoseconds()),
			ServiceDur:   float64(serDur[serv].Nanoseconds()),
			TotalPercent: float64(serExc[serv].Nanoseconds()) / float64(totalDur),
			CountPercent: float64(serCount[serv]) / float64(totalSpan),
			TotalDur:     totalDur,
			TotalService: totalSer,
			TotalSpan:    totalSpan,
			Load:         load,
			BottleNeck:   isBottleNeck,
		}
		serMetric = append(serMetric, serviceMetric)
	}
	return serMetric
}
