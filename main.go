package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

func main() {
	/*************************
	Generate some requests
	*************************/
	// search tickets advance
	// url := "http://34.123.147.40:32677/api/v1/travelplanservice/travelPlan/cheapest"
	// search tickets
	url := "http://34.123.147.40:32677/api/v1/travel2service/trips/left"
	body := []byte(`{"startingPlace": "Shang Hai", "endPlace": "Tai Yuan", "departureTime": "2020-11-25"}`)
	iteration := 10
	startTime := time.Now()
	SendRequest(url, body, iteration)
	fmt.Println("Request generation completed")

	/*************************
	Get tracing information
	*************************/
	// Create jaeger service client
	jc := NewJaegerClient("localhost:16686")
	// Query traces
	traces, err := jc.QueryTraces("ts-ui-dashboard.default", "", startTime, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := json.Marshal(traces)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := ioutil.WriteFile(filepath.Join("data", "traces"), result, 0644); err != nil {
		fmt.Println(err)
		return
	}
	// critcalPath := CleanupTraces(traces)
	// fmt.Println(critcalPath)
	// for _, trace := range traces {
	// 	ConstructTrainingData(trace, iteration)
	// }
}
