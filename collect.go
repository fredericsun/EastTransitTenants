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

type TestServiceMetric struct {
	ServiceName  string
	SelfPercent  float64
	ServiceDur   float64
	TotalPercent float64
	CountPercent float64
	TotalDur     int64
	TotalService int
	TotalSpan    int
	Load         int
}

func ConstructTrainingData(spans []MySpan, load int, criticalPath []string) []ServiceMetric {
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

	// request information
	totalSer := len(serExc)
	totalSpan := len(spans)
	totalDur := latestEndTime.UnixNano() - earliestStartTime.UnixNano()

	var serMetric []ServiceMetric
	candidates := make(map[string]bool)
	for _, serv := range criticalPath {
		candidates[serv] = true
	}
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

func ConstructTestingData(spans []MySpan, load int) []TestServiceMetric {
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

	// request information
	totalSer := len(serExc)
	totalSpan := len(spans)
	totalDur := latestEndTime.UnixNano() - earliestStartTime.UnixNano()

	var serMetric []TestServiceMetric
	for serv, _ := range serExc {
		serviceMetric := TestServiceMetric{
			ServiceName:  serv,
			SelfPercent:  float64(serExc[serv].Nanoseconds()) / float64(serDur[serv].Nanoseconds()),
			ServiceDur:   float64(serDur[serv].Nanoseconds()),
			TotalPercent: float64(serExc[serv].Nanoseconds()) / float64(totalDur),
			CountPercent: float64(serCount[serv]) / float64(totalSpan),
			TotalDur:     totalDur,
			TotalService: totalSer,
			TotalSpan:    totalSpan,
			Load:         load,
		}
		serMetric = append(serMetric, serviceMetric)
	}
	return serMetric
}
