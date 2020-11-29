package main

import (
	"time"
)

func getSortedService(spans []MySpan) []string {
	servicesCount := make(map[string]int)
	for _, span := range spans {
		servicesCount[span.Process.GetServiceName()] += 1
	}
	servicesSorted := rankByWordCount(servicesCount)
	services := make([]string, 0)
	for _, pair := range servicesSorted {
		services = append(services, pair.Service)
	}
	return services
}

func findBottleNeck(spans []MySpan) []string {
	spanToServ := make(map[string]string)
	for _, span := range spans {
		spanToServ[span.SpanID] = span.Process.GetServiceName()
	}
	// fmt.Println(spanToServ)

	serExc := make(map[string]time.Duration)
	for _, span := range spans {
		serExc[span.Process.GetServiceName()] += span.Duration
		for _, ref := range span.References {
			if ref["RefType"] == "CHILD_OF" {
				if _, ok := spanToServ[ref["SpanID"]]; !ok {
					// fmt.Printf("nonexist: %s\n", ref["SpanID"])
				}
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
	candidates := make(map[string]bool)
	candidates[res] = true
	for serv, dur := range serExc {
		if float64(dur.Nanoseconds()) >= float64(max.Nanoseconds())*0.7 {
			candidates[serv] = true
		}
	}
	// fmt.Println(serExc)
	result := make([]string, 0)
	for serv := range candidates {
		result = append(result, serv)
	}
	return result
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
		stack = stack[:n]
		service, from := servParent.service, servParent.parent
		if _, ok := done[service]; ok {
			continue
		}
		done[service] = true
		parent[service] = from
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

func exists(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CleanupTraces(traces map[string][]MySpan) []string {
	bottlenecks := make(map[string]bool)
	for _, spans := range traces {
		found := findBottleNeck(spans)
		// fmt.Printf("%s: %s\n", id, bottleneck)
		for _, bottleneck := range found {
			bottlenecks[bottleneck] = true
		}
	}
	// fmt.Println(bottlenecks)

	callGraph := make(map[string][]string)
	root := ""
	for _, spans := range traces {
		spanToServ := make(map[string]string)
		for _, span := range spans {
			spanToServ[span.SpanID] = span.Process.GetServiceName()
		}

		for _, span := range spans {
			serv := span.Process.GetServiceName()
			if len(span.References) == 0 {
				root = serv
			}
			for _, ref := range span.References {
				parentServ := spanToServ[ref["SpanID"]]
				if parentServ != serv && !exists(serv, callGraph[parentServ]) {
					callGraph[parentServ] = append(callGraph[parentServ], serv)
				}
			}
		}
		if _, ok := callGraph[""]; !ok {
			break
		}
		callGraph = make(map[string][]string)
	}
	// fmt.Println(root)
	// fmt.Println(callGraph)

	return findCriticalPaths(root, callGraph, bottlenecks)
}
