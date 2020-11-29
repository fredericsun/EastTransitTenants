package main

import (
	"fmt"
	"os/exec"
)

var format = `apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: seat-latency
spec:
  hosts:
  - %s
  http:
  - fault:
      delay:
        percentage:
          value: 100.0
        fixedDelay: %dms
    route:
    - destination:
        host: %s
  - route:
    - destination:
        host: %s`

func contruct_yaml(service string, latency_ms int) string {
	content := fmt.Sprintf(format, service, latency_ms, service, service)
	return content
}

func inject_latency() {
	cmd := exec.Command("sh", "latency.sh")
	res, _ := cmd.Output()
	fmt.Println(string(res))
}
