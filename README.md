# Message API

The message api is an API service used behind the Micro API gateway. 

## Dependence on Service
- [message-srv](https://github.com/microhq/message-srv)

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Download and start the service

	```shell
	go get github.com/micro/message-api
	message-api
	```

	OR as a docker container

	```shell
	docker run microhq/message-api --registry_address=YOUR_REGISTRY_ADDRESS
	```
