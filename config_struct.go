package main

type RequestData struct {
	url    string
	body   []byte
	bearer string
}

type ProfileRequestData struct {
	url    string
	body   string
	bearer string
}

type TrainRequestData struct {
	name string
	url  string
	body string
}

type TrainConfig struct {
	request     []TrainRequestData
	bearer      string
	jaeger_ip   string
	workload    []int
	target_serv string
}

type ProfileConfig struct {
	jaeger_ip     string
	request       ProfileRequestData
	workload      int
	entrance_serv string
	target_serv   string
	precision     int
}
