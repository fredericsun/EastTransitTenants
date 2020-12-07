# EastTransitTenants

## Introduction
This project is a performance profiling and bottleneck analysis framework built on top of an open-source railway ticketingmicroservices system. From a big picture,our framework consists of five components:<br />
- 1)orchestrator:coordinate the interaction of other functional components;<br /><br />
- 2)load generator: generate workload for training and profiling purpose;<br /><br />
- 3)benchmark application:  an opensource railway ticketing microservices system which contains 41 mi-croservices and whose request can involve up to 200 spans (You can get more details at https://github.com/FudanSELab/train-ticket/wiki);<br /><br />
- 4)bottleneck discover: find all potential bottleneck services ina time-efficient way;<br /><br />
- 5)performance profiler: instrumentallyprofile how a service influencing the end-to-end performanceof the requests it belongs to.

## Architecture Graph
![architecture](./architecture.png)

## Usage
First, build the framework with:
```go build```
<br /><br />
After built-up, you can choose either to train the bottleneck-discovery model or start the performance profiling and bottleneck analysis by run the project with different flags.
### Train Bottleneck-discovery Model
#### 1. Configuration for training
The input configurations are defined in ```train_config.json ```. We have provided a template configuration in ths file. You can modify it based on your application and request types.
```
{
    "request": [{
        "name": "search_tickets",
        "url": "http://35.231.88.215:32677/api/v1/travel2service/trips/left",
        "body": "{\"startingPlace\": \"Shang Hai\",\"endPlace\": \"Tai Yuan\",\"departureTime\": \"2020-12-21\"}"
    }],
    "bearer": "",
    "jaeger_ip": "35.231.88.215:32688",
    "workload": [1, 10, 25, 50, 100],
    "target_serv" : "ts-ui-dashboard.default"
}
```
- request: list of requests to generate training data. Each of them should contain the name of the request, the url of its API call, and the request body in json format
- bearer: authorization bearer token. Only provided if needed
- jaeger_ip: the IP address of Jaeger
- target_serv: the entrance service for all requests to your application, which is normally the front-end service.

#### 2. Train the model
```./easttransittenants -type=train```

#### Presequisite
* Docker
* Docker Compose

#### 1. Clone the Repository
```bash
git clone --depth=1 https://github.com/FudanSELab/train-ticket.git
cd train-ticket/
```

#### 2. Start the Application
```bash
docker-compose -f deployment/docker-compose-manifests/quickstart-docker-compose.yml up
```
Once the application starts, you can visit the Train Ticket web page at [http://localhost:8080](http://localhost:8080).

### Using Kubernetes
Here is the steps to deploy the Train Ticket onto any existing Kubernetes cluster.

#### Presequisite
* An existing Kubernetes cluster

```



