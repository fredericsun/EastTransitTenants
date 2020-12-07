package main

type RequestData struct {
	url    string
	body   []byte
	bearer string
}

type ProfileRequestData struct {
	Url    string `json:"url"`
	Body   string `json:"body"`
	Bearer string `json:"bearer"`
}

type TrainRequestData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Body string `json:"body"`
}

type TrainConfig struct {
	Request     []TrainRequestData `json:"request"`
	Bearer      string             `json:"bearer"`
	Jaeger_ip   string             `json:"jaeger_ip"`
	Workload    []int              `json:"workload"`
	Target_serv string             `json:"target_serv"`
}

type ProfileConfig struct {
	Jaeger_ip   string             `json:"jaeger_ip"`
	Request     ProfileRequestData `json:"request"`
	Workload    int                `json:"workload"`
	Target_serv string             `json:"target_serv"`
	Precision   int                `json:"precision_ms"`
}
