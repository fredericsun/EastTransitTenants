package main

import (
	"fmt"
	"os/exec"
)

var format = `apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: %s-latency
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

func constructYaml(service string, latency_ms int) string {
	content := fmt.Sprintf(format, service, service, latency_ms, service, service)
	return content
}

func injectLatency(service string, latency_ms int) {
	cmdFormat := `cat <<EOF | kubectl apply -f -
  %s
  EOF`
	yamlContent := constructYaml(service, latency_ms)
	command := fmt.Sprintf(cmdFormat, yamlContent)
	cmd := exec.Command(command)
	res, _ := cmd.CombinedOutput()
	fmt.Println(string(res))
}

func deleteLatency() {}
