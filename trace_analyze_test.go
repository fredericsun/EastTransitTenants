package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

type TestCase struct {
	traces         map[string][]MySpan
	expectedResult []string
}

func TestCleanUpTraces(t *testing.T) {
	file, _ := ioutil.ReadFile("littletrace.json")
	// file, _ := ioutil.ReadFile("alltraces.json")
	// fmt.Println(string(file))
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	traces := make(map[string][]MySpan)
	err := json.Unmarshal(file, &traces)
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	// fmt.Println(traces)

	// result, err := json.Marshal(traces)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(result))

	criticalServices := CleanupTraces(traces)
	fmt.Println(criticalServices)

	// for _, service := range testCase.expectedResult {
	// 	expected[service] = true
	// }
	// if !reflect.DeepEqual(actual, expected) {
	// 	t.Errorf("Expected critical to be %v, got %v\n", expected, actual)
	// }
}
