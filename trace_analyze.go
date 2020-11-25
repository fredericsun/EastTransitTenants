package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func findBottleNeck(spans []MySpan, spanToServ map[string]string) string {
	serExc := make(map[string]time.Duration)
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

	return res
}

type pair struct {
	service string
	parent  string
}

func findCriticalPaths(root string, callGraph map[string][]string, bottlenecks map[string]bool) []string {
	stack := make([]pair, 0)
	done := make(map[string]bool)
	potential := make(map[string]bool)
	critical := make(map[string]bool)
	parent := make(map[string]string)

	stack = append(stack, pair{root, ""})
	for len(stack) > 0 {
		n := len(stack) - 1
		servParent := stack[n]
		service, from := servParent.service, servParent.parent
		if _, ok := done[service]; ok {
			continue
		}
		done[service] = true
		parent[service] = from
		stack = stack[:n]
		if _, ok := bottlenecks[service]; ok {
			critical[service] = true

			curService := service
			for {
				if _, ok := potential[curService]; !ok {
					break
				}
				critical[curService] = true
				delete(potential, curService)
				curService = parent[curService]
			}

			for _, child := range callGraph[service] {
				stack = append(stack, pair{child, service})
				potential[child] = true
			}

		} else {
			for _, child := range callGraph[service] {
				stack = append(stack, pair{child, service})
				if _, ok := potential[service]; ok {
					potential[child] = true
				}
			}
		}
	}

	criticalPath := make([]string, 0)
	for service := range critical {
		criticalPath = append(criticalPath, service)
	}
	return criticalPath
}

func CleanupTraces(traces map[string][]MySpan) []string {
	spanToServ := make(map[string]string)
	for _, spans := range traces {
		for _, span := range spans {
			spanToServ[span.SpanID] = span.Process.GetServiceName()
		}
		break
	}

	bottlenecks := make(map[string]bool)
	for _, spans := range traces {
		bottleneck := findBottleNeck(spans, spanToServ)
		bottlenecks[bottleneck] = true
	}
	bin, _ := json.Marshal(bottlenecks)
	fmt.Println(bin)

	callGraph := make(map[string][]string)
	root := ""
	for _, spans := range traces {
		for _, span := range spans {
			serv := span.Process.GetServiceName()
			if len(span.References) > 0 {
				root = serv
			}
			for _, ref := range span.References {
				parentServ := spanToServ[ref["SpanID"]]
				if parentServ != serv {
					callGraph[parentServ] = append(callGraph[parentServ], serv)
				}
			}
		}
		break
	}

	return findCriticalPaths(root, callGraph, bottlenecks)
}
