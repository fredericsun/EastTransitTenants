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
		err := json.Unmarshal(config_json, &config)
		if err != nil {
			panic(fmt.Errorf("Error in parsing json %v", err))
		}
		JaegerIP = config.Jaeger_ip
		generateTrainingData(config.Request, config.Bearer, config.Workload, config.Target_serv)
		train()
	}
	if task_type == "profile" {
		var config ProfileConfig
		config_json, _ := ioutil.ReadFile("profile_config.json")
		json.Unmarshal([]byte(config_json), &config)
		JaegerIP = config.Jaeger_ip
		writePath(config.Request.Url, config.Request.Body, config.Request.Bearer, config.Workload, jaeger_sleep, config.Target_serv)
		bottlenecks := predict()

		requestData := RequestData{
			url:    config.Request.Url,
			body:   []byte(config.Request.Body),
			bearer: config.Request.Bearer,
		}

		for name, isBottleneck := range bottlenecks {
			if isBottleneck {
				boundary := lookForBoundary(requestData, name, config.Precision, config.Workload)
				if boundary == config.Precision/2 {
					fmt.Println("Cannot find the boundary of the current service")
					continue
				}
				// end-to-end baseline
				virtualSpeedUp("everything", boundary)
				baseline := SendRequest(requestData.url, requestData.body, config.Workload, requestData.bearer)
				removeVirtualSpeedUp(nil)

				virtualSpeedUp(name, boundary)
				speedupResult := SendRequest(requestData.url, requestData.body, config.Workload, requestData.bearer)
				removeVirtualSpeedUp(nil)

				fmt.Printf("baseline %dms, simulation %dms, improved %dms\n", baseline, speedupResult, (baseline - speedupResult))
			}
		}
	}
}
