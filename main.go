package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const JaegerIP = "localhost:16686"
const RequestBearer = "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJmZHNlX21pY3Jvc2VydmljZSIsInJvbGVzIjpbIlJPTEVfVVNFUiJdLCJpZCI6IjRkMmE0NmM3LTcxY2ItNGNmMS1iNWJiLWI2ODQwNmQ5ZGE2ZiIsImlhdCI6MTYwNjY4ODkwNCwiZXhwIjoxNjA2NjkyNTA0fQ.jnTJBJg_PzZlbJw6QTPBsCk9jPpGqkh53JMs94t6Jkk"

// const RequestUrl = "http://35.231.88.215:32677/api/v1/travel2service/trips/left"
// const RequestBody = `{"startingPlace":"Shang Hai","endPlace":"Tai Yuan","departureTime":"2020-12-21"}`
// const filename = "ticket_reserve"

type RequstData struct {
	url    string
	body   []byte
	bearer string
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
		// Write Traces out
		// traceLog, err := json.Marshal(traces)
		// if err != nil {
		// 	fmt.Println("Error in marshalling")
		// }
		// if err := ioutil.WriteFile(filepath.Join("data", "traces"), traceLog, 0644); err != nil {
		// 	fmt.Println(err)
		// }

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

// Search Tickets
// url := "http://35.231.88.215:32677/api/v1/travelplanservice/travelPlan/cheapest"
// url := "http://35.231.88.215:32677/api/v1/travel2service/trips/left"
// body := []byte(`{"startingPlace": "Shang Hai", "endPlace": "Su Zhou", "departureTime": "2020-11-28"}`)

// Get Orders
// url := "http://35.231.88.215:32677/api/v1/orderservice/order/refresh"
// body := []byte(`{"loginId":"4d2a46c7-71cb-4cf1-b5bb-b68406d9da6f","enableStateQuery":false,"enableTravelDateQuery":false,"enableBoughtDateQuery":false,"travelDateStart":null,"travelDateEnd":null,"boughtDateStart":null,"boughtDateEnd":null}`)

// Preserve Tickets
// url := "http://35.231.88.215:32677/api/v1/preserveotherservice/preserveOther"
// body := []byte(`{"accountId":"4d2a46c7-71cb-4cf1-b5bb-b68406d9da6f","contactsId":"8bfcba2c-c777-459b-aeb8-c07a22b489de","tripId":"Z1234","seatType":"2","date":"2020-12-22","from":"Shang Hai","to":"Tai Yuan","assurance":"0","foodType":1,"foodName":"Bone Soup","foodPrice":2.5,"stationName":"","storeName":""}`)

func main() {
	// batch := flag.Int("load", 100, "requests batch size")
	// target := flag.String("target", "", "get traces involving such service")
	sleep := flag.Int("sleep", 5, "Wait until Jaeger is updated")
	flag.Parse()
	// iteration := *batch
	target_service := "ts-ui-dashboard.default"
	jaeger_sleep := *sleep

	toRun := make(map[string]RequstData)
	toRun["search_tickets"] = RequstData{
		url:  "http://35.225.46.132:32677/api/v1/travel2service/trips/left",
		body: []byte(`{"startingPlace":"Shang Hai","endPlace":"Tai Yuan","departureTime":"2020-12-21"}`),
	}
	toRun["get_orders"] = RequstData{
		url:  "http://35.225.46.132:32677/api/v1/orderservice/order/refresh",
		body: []byte(`{"loginId":"4d2a46c7-71cb-4cf1-b5bb-b68406d9da6f","enableStateQuery":false,"enableTravelDateQuery":false,"enableBoughtDateQuery":false,"travelDateStart":null,"travelDateEnd":null,"boughtDateStart":null,"boughtDateEnd":null}`),
	}
	toRun["preserve_tickets"] = RequstData{
		url:    "http://35.225.46.132:32677/api/v1/preserveotherservice/preserveOther",
		body:   []byte(`{"accountId":"4d2a46c7-71cb-4cf1-b5bb-b68406d9da6f","contactsId":"9b15edd9-ead7-477e-836b-fae38223bb40","tripId":"T1235","seatType":"2","date":"2020-12-21","from":"Shang Hai","to":"Tai Yuan","assurance":"0","foodType":1,"foodName":"Bone Soup","foodPrice":2.5,"stationName":"","storeName":""}`),
		bearer: RequestBearer,
	}

	for filename, requestData := range toRun {
		for _, iteration := range []int{1, 10, 100} {
			generate(requestData.url, requestData.body, requestData.bearer, filename, iteration, target_service, jaeger_sleep)
		}
	}
}

// func main() {
// 	inject_latency()
// }
