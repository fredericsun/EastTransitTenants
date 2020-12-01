package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var yamlFormat = `apiVersion: networking.istio.io/v1alpha3
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

var injectCmdFormat = `cat <<EOF | kubectl apply -f -
%s
EOF`

var deleteCmdFormat = `cat <<EOF | kubectl delete -f -
%s
EOF`

func constructYaml(service string, latency_ms int) string {
	content := fmt.Sprintf(yamlFormat, service, service, latency_ms, service, service)
	return content
}

func writeProfile(profilePath string, services []string, latencies []int) {
	if len(services) != len(latencies) {
		panic(fmt.Errorf("expecting same size but got service: %d, latencies: %d", len(services), len(latencies)))
	}

	objects := make([]string, 0)
	for i, service := range services {
		latency := latencies[i]
		virtService := constructYaml(service, latency)
		objects = append(objects, virtService)
	}
	profile := strings.Join(objects, "\n---\n")

	f, err := os.OpenFile(profilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := f.Write([]byte(profile)); err != nil {
		fmt.Println(err)
	}
}

func buildProfile(profileName string, targetService string, latencyAmount int) string {
	profilePath := filepath.Join("tmp", fmt.Sprintf("%s.yaml", profileName))
	if _, err := os.Stat(profilePath); err == nil {
		// path/to/whatever does not exist
		return profilePath
	}

	targetService = strings.Split(targetService, ".")[0]
	jc := NewJaegerClient(JaegerIP)
	allServices, err := jc.QueryServices()
	if err != nil {
		panic(err)
	}
	services := make([]string, 0)
	latencies := make([]int, 0)
	for _, ser := range allServices.Services {
		if !strings.HasPrefix(ser, "ts-") {
			continue
		}
		ser = strings.Split(ser, ".")[0]
		if ser == targetService {
			continue
		}
		services = append(services, ser)
		latencies = append(latencies, latencyAmount)
	}

	writeProfile(profilePath, services, latencies)
	return profilePath
}

// func injectLatency(service string, latency_ms int) {
// 	script := createScript("apply", service, latency_ms)
// 	cmd := exec.Command("sh", script)
// 	res, err := cmd.Output()
// 	if err != nil {
// 		fmt.Printf("error: %s\n", err)
// 	}
// 	fmt.Println(string(res))
// }

// func deleteLatency(service string, latency_ms int) {
// 	script := createScript("delete", service, latency_ms)
// 	cmd := exec.Command("sh", script)
// 	res, err := cmd.Output()
// 	if err != nil {
// 		fmt.Printf("error: %s\n", err)
// 	}
// 	fmt.Println(string(res))
// }

func virtualSpeedUp(service string, amount_ms int) {
	profileName := fmt.Sprintf("target_%s_latency_%dms", service, amount_ms)
	profilePath := buildProfile(profileName, service, amount_ms)

	cmd := exec.Command("kubectl", "apply", "-f", profilePath)
	res, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	fmt.Println(string(res))
}

func removeVirtualSpeedUp(result chan bool) {
	profileName := "resumeProfile"
	profilePath := buildProfile(profileName, "", 0)

	cmd := exec.Command("kubectl", "delete", "-f", profilePath)
	res, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	fmt.Println(string(res))
	if result != nil {
		result <- true
	}
}
