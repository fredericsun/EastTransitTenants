package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func SendRequest(url string, body []byte, iteration int, bearer string) int {
	finished := make(chan int64)

	request := func() {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Close = true
		req.Header.Set("X-Custom-Header", "loader")
		req.Header.Set("Content-Type", "application/json")
		if bearer != "" {
			req.Header.Add("Authorization", bearer)
		}

		client := &http.Client{}
		start := time.Now()
		resp, err := client.Do(req)
		duration := time.Since(start)
		if err != nil {
			finished <- -1
		}
		status := resp.Status
		resp.Body.Close()
		if status == "200 OK" {
			finished <- duration.Milliseconds()
		} else {
			finished <- -1
		}
		fmt.Println("response Status:", status)
	}
	for i := 0; i < iteration; i++ {
		go request()
	}
	total := int64(0)
	count := 0
	for i := 0; i < iteration; i++ {
		duration := <-finished
		if duration > 0 {
			count += 1
			total += duration
		}
	}
	fmt.Println("The number of success request is:", count)
	return int(total) / count
}
