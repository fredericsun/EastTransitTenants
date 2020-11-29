package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	/*************************
	Generate some requests
	*************************/

	filename := flag.String("name", "trainingData", "train data name")
	batch := flag.Int("iter", 10, "requests batch size")
	flag.Parse()
	iteration := *batch
	train_filename := *filename
	fmt.Printf("Running with iteration: %d\n", iteration)

	bearer := ""
	// Search Tickets
	// url := "http://35.231.88.215:32677/api/v1/travelplanservice/travelPlan/cheapest"
	url := "http://35.231.88.215:32677/api/v1/travel2service/trips/left"
	body := []byte(`{"startingPlace": "Shang Hai", "endPlace": "Su Zhou", "departureTime": "2020-11-28"}`)

	// Get Orders
	// url := "http://35.231.88.215:32677/api/v1/orderservice/order/refresh"
	// body := []byte(`{"loginId":"4d2a46c7-71cb-4cf1-b5bb-b68406d9da6f","enableStateQuery":false,"enableTravelDateQuery":false,"enableBoughtDateQuery":false,"travelDateStart":null,"travelDateEnd":null,"boughtDateStart":null,"boughtDateEnd":null}`)

	// Preserve Tickets
	// url := "http://35.231.88.215:32677/api/v1/preserveotherservice/preserveOther"
	// body := []byte(`{"accountId":"4d2a46c7-71cb-4cf1-b5bb-b68406d9da6f","contactsId":"8bfcba2c-c777-459b-aeb8-c07a22b489de","tripId":"Z1234","seatType":"2","date":"2020-12-22","from":"Shang Hai","to":"Tai Yuan","assurance":"0","foodType":1,"foodName":"Bone Soup","foodPrice":2.5,"stationName":"","storeName":""}`)

	// Pay Ticket
	// url := "http://35.231.88.215:32677/api/v1/inside_pay_service/inside_payment"
	// body := []byte(`{"orderId":"be83b816-87cf-43fc-9988-0b50cd733e82","tripId":"Z1234"}`)
	var trainData []ServiceMetric
	for i := 0; i < 100/iteration; i++ {
		startTime := time.Now()
		SendRequest(url, body, iteration, bearer)
		fmt.Println("Request generation completed")
		time.Sleep(time.Duration(10) * time.Second)

		/*************************
		Get tracing information
		*************************/
		// Create jaeger service client
		jc := NewJaegerClient("35.231.88.215:32688")
		// Query traces
		traces, err := jc.QueryTraces("ts-travel2-service", "", startTime, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		blah, err := json.Marshal(traces)
		if err != nil {
			fmt.Println("Error in marshalling")
		}
		if err := ioutil.WriteFile(filepath.Join("data", "traces"), blah, 0644); err != nil {
			fmt.Println(err)
		}
		critcalPath := CleanupTraces(traces)
		fmt.Println(critcalPath)
		// Construct training data
		for _, trace := range traces {
			serMetric := ConstructTrainingData(trace, iteration)
			trainData = append(trainData, serMetric...)
		}
	}
	result, err := json.Marshal(trainData)
	if err != nil {
		fmt.Println(err)
	}
	// if err := ioutil.WriteFile(filepath.Join("data", "training"), result, 0644); err != nil {
	// 	fmt.Println(err)
	// }
	f, err := os.OpenFile(filepath.Join("data", fmt.Sprintf("%s_%d.json", train_filename, iteration)), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := f.Write(result); err != nil {
		fmt.Println(err)
	}
}

// func main() {
// 	inject_latency()
// }
