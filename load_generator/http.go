package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func sendRequest(url string, body []byte, iteration int) {
	finished := make(chan bool)

	request := func() {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("X-Custom-Header", "loader")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		if resp.Status == "200 OK" {
			finished <- true
		} else {
			finished <- false
		}
	}
	for i := 0; i < iteration; i++ {
		go request()
	}
	count := 0
	for i := 0; i < iteration; i++ {
		if <-finished {
			count += 1
		}
	}
	fmt.Println("The number of success request is:", count)
}

func main() {
	// searcg tickets advance
	// url := "http://34.123.147.40:32677/api/v1/travelplanservice/travelPlan/cheapest"
	// search tickets
	url := "http://34.123.147.40:32677/api/v1/travel2service/trips/left"
	body := []byte(`{"startingPlace": "Shang Hai", "endPlace": "Tai Yuan", "departureTime": "2020-11-25"}`)
	iteration := 3
	sendRequest(url, body, iteration)
}
