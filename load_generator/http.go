package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	finished := make(chan bool)

	fn := func() {
		url := "http://34.123.147.40:32677/api/v1/travelservice/trips/left"
		body := []byte(`{"startingPlace": "Shang Hai", "endPlace": "Su Zhou", "departureTime": "2020-11-21"}`)

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
		finished <- true
	}
	for i := 0; i < 10; i++ {
		go fn()
	}
	for i := 0; i < 10; i++ {
		<-finished
	}
}
