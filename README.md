# EastTransitTenants

This project is a performance profiling and bottleneck analysis framework built on top of an open-source railway ticketingmicroservices system. From a big picture,our framework consists of five components:
1)orchestrator:coordinate the interaction of other functional components;
2)load generator: generate workload for training and profiling purpose;  
3)benchmark application:  an opensource railway ticketing microservices system which contains 41 mi-croservices and whose request can involve up to 200 spans (You can get more details at [Wiki Pages](https://github.com/FudanSELab/train-ticket/wiki).); 
4)bottleneck discover: find all potential bottleneck services ina time-efficient way; 
5)performance profiler: instrumentallyprofile how a service influencing the end-to-end performanceof the requests it belongs to.

## Service Architecture Graph
![architecture](./image/architecture.png)

## Quick Start
We provide two options to quickly deploy our application: [Using Docker Compose](#Using-Docker-Compose) and [Using Kubernetes](#Using-Kubernetes).

### Using Docker Compose
The easiest way to get start with the Train Ticket application is by using [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).

> If you don't have Docker and Docker Compose installed, you can refer to [the Docker website](https://www.docker.com/) to install them.

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



