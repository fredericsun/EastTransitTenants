package easttransittenants

import "time"

func constructTrainingData(spans []MySpan, load int) {

	serSelfPercent := make(map[string]float64)
	serTotalPercent := make(map[string]float64)
	serCountPercent := make(map[string]float64)
	serCand := make(map[string]bool)

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
		serCountPercent[serv] = float64(serCount[serv]) / float64(totalSpan)
		serSelfPercent[serv] = float64(serExc[serv].Nanoseconds()) / float64(serDur[serv].Nanoseconds())
		serTotalPercent[serv] = float64(serExc[serv].Nanoseconds()) / float64(totalDur)
		if res == serv {
			serCand[serv] = true
		} else {
			serCand[serv] = false
		}
	}
}
