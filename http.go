package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendRequest(url string, body []byte, iteration int, bearer string) {
	finished := make(chan bool)

	request := func() {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("X-Custom-Header", "loader")
		req.Header.Set("Content-Type", "application/json")
		if bearer != "" {
			req.Header.Add("Authorization", bearer)
		}

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
