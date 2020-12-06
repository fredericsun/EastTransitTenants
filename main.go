package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

var JaegerIP = "35.231.88.215:32688"
var jaeger_sleep = 5

// Some example requests
var Example_ToRun = map[string]RequestData{
	"search_tickets": RequestData{
		url:  "http://35.231.88.215:32677/api/v1/travel2service/trips/left",
		body: []byte(`{"startingPlace":"Shang Hai","endPlace":"Tai Yuan","departureTime":"2020-12-21"}`),
	},
	"get_orders": RequestData{
		url:  "http://35.225.46.132:32677/api/v1/orderservice/order/refresh",
		body: []byte(`{"loginId":"4d2a46c7-71cb-4cf1-b5bb-b68406d9da6f","enableStateQuery":false,"enableTravelDateQuery":false,"enableBoughtDateQuery":false,"travelDateStart":null,"travelDateEnd":null,"boughtDateStart":null,"boughtDateEnd":null}`),
	},
	"preserve_tickets": RequestData{
		url:  "http://35.231.88.215:32677/api/v1/preserveotherservice/preserveOther",
		body: []byte(`{"accountId":"4d2a46c7-71cb-4cf1-b5bb-b68406d9da6f","contactsId":"dfbe3b5f-ca3a-4d1d-90dd-d6ba5fb15ab2","tripId":"T1235","seatType":"2","date":"2020-12-21","from":"Shang Hai","to":"Tai Yuan","assurance":"0","foodType":1,"foodName":"Bone Soup","foodPrice":2.5,"stationName":"","storeName":""}`),
	},
}

func main() {
	t := flag.String("type", "", "task type")
	flag.Parse()
	task_type := *t
	if task_type == "train" {
		var config TrainConfig
		config_json, _ := ioutil.ReadFile("train_config.json")
		json.Unmarshal([]byte(config_json), &config)
		JaegerIP = config.jaeger_ip
		generateTrainingData(config.request, config.bearer, config.workload, config.target_serv)
		train()
	}
	if task_type == "profile" {
		var config ProfileConfig
		config_json, _ := ioutil.ReadFile("profile_config.json")
		json.Unmarshal([]byte(config_json), &config)
		JaegerIP = config.jaeger_ip
		writePath(config.request.url, config.request.body, config.request.bearer, config.workload, jaeger_sleep, config.target_serv)
		bottlenecks := predict()

		requestData := RequestData{
			url:    config.request.url,
			body:   []byte(config.request.body),
			bearer: config.request.bearer,
		}

		for name, isBottleneck := range bottlenecks {
			if isBottleneck {
				boundary := lookForBoundary(requestData, name, config.precision, config.workload)
				if boundary == config.precision/2 {
					fmt.Println("Cannot find the boundary of the current service")
					continue
				}
				// end-to-end baseline
				virtualSpeedUp("everything", boundary)
				baseline := SendRequest(requestData.url, requestData.body, config.workload, requestData.bearer)
				removeVirtualSpeedUp(nil)

				virtualSpeedUp(name, boundary)
				speedupResult := SendRequest(requestData.url, requestData.body, config.workload, requestData.bearer)
				removeVirtualSpeedUp(nil)

				fmt.Printf("baseline %dms, simulation %dms, improved %dms\n", baseline, speedupResult, (baseline - speedupResult))
			}
		}
	}
}
