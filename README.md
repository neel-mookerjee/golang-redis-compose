# A Golang API and Redis Composed Together
This sample application uses Docker compose to host a Golang RESTFul service with a backend storage.

## Prerequisite

To use this project:

- Docker to build images
- Golang version 1.14 `brew install go`
* Swag v1.4.1 is needed
- Nice to have [make](https://www.gnu.org/software/make/manual/make.html)

## Usecase

- Use a Docker container to host your service.
- Provide a service endpoint that will generate a unique URI.
- Allow POST operations to that generated endpoint.
- Store the contents (BODY) of the POST to a persistent datastore.
- Include a health check that will monitor your service.
- Accept GET requests to the service endpoint that will return the body of the most recent POST operation.

## Design
1. An Http Endpoint accepts PUT request to generate unique endpoint.
1. The docker-compose network is named `backend`. The internal service API base would be http://app:8080. However ports are mapped to make the URI accessible from the host machine using _localhost_.
1. A Redis instance is used to store the newly generated URIs as KVP.
1. Payloads submitted to the endpoint are stored in a Redis List (add to the head of the linkedlist).
1. Most recently submitted payload is retrieved using list range [0, 0].
1. A health endpoint is also made available for the generated endpoints.
1. A canary is deployed on its own container to check the health of the generated endpoints.
1. Project [api](api) has the RESTFul endpoints and project [canary](canary) has the health checks implemented.
1. Sample basic tests are added to the api project.
1. The canary uses the Redis instance as the repo for the URIs (lol). Ideally every endpoint creation should register the endpoint in a telemetry system.
1. A client container is added to use the endpoint with the network `backend`.
1. Swagger specs are generated and served from a static file.
1. Check the specs for error codes and responses.

## Development

```bash
$ make

  ## GOLANG-REDIS-COMPOSE

  golang-redis-compose =>    make swagger/doc               Generate Swagger specs. Requires swag.
  golang-redis-compose =>    make go/test                   Run the tests
  golang-redis-compose =>    make docker/build              Build the docker images
  golang-redis-compose =>    make docker/deploy             Deploy the containers using docker-compose
  golang-redis-compose =>    make docker/destroy            Stop and remove the containers
  golang-redis-compose =>    make help                      This is default and it helps
```

Refer the commands in the [Makefile](Makefile) to run without Make!

## Endpoints

Locally deployed instances can be accessed as the following:

| Endpoint | Description |
| ----- | ------ |
| [http://localhost:8080/generate](http://localhost:8080/generate) | PUT to Create an Endpoint which is unique and accepts GET/POST |
| Generated URL (e.g. [http://localhost:8080/R6jxp3z1wGXJw](http://localhost:8080/R6jxp3z1wGXJw])) | POST to submit payloads and GET to retrieve the most recent payload |
| [http://localhost:8080/](http://localhost:8080/) | API Spec that can be rendered using Swagger UI |

### Swagger API Spec
The Swagger Spec for the APIs can be found here:<br/>
[http://localhost:8080/swagger/index.html?url=http://localhost:8080](http://localhost:8080/swagger/index.html?url=http://localhost:8080)

## Using the Client Container
From within the client container the endpoints can be invoked using the network `backend`. A sample flow of commands and their outputs are as the following:

```bash
bash$ docker exec -it golangrediscompose_client_1 bash
bash-4.4# curl -X PUT http://app:8080/generate
{"appurl":"http://app:8080/W6oDZwzEKnOJQ","note":"To use from outside the docker environment use: http://localhost:8080/W6oDZwzEKnOJQ"}
bash-4.4# curl -d '{"message":"hello"}' -H 'Content-Type: application/json' http://app:8080/W6oDZwzEKnOJQ
{"status":"submitted"}
bash-4.4# curl -d '{"message":"hello again"}' -H 'Content-Type: application/json' http://app:8080/W6oDZwzEKnOJQ
{"status":"submitted"}
bash-4.4# curl http://app:8080/W6oDZwzEKnOJQ
{"message":"hello again"}
```

## Callouts
1. The Redis instance is not having a persistence storage. Currently the redis storage works as a cache valid for 24 hours! If persistence is enabled, data is stored in the VOLUME /data, which can be used with `-v /docker/host/dir:/data`. The docker command must run with the following params: `redis-server --appendonly yes`
1. A scalable orchestration (docker swarm, mesos, k8) is recommended for serving the newly generated endpoints with good performance. Multiple instances will distribute the load and provide high availability.
1. Can concurrency cause problems?
1. Size limits?
1. Throttling and rate-limiting are necessary. User based access token can be used to, among other things, throttle users based on their allocated quota.
1. The canary uses the Redis instance as the repo for the URIs (lol). Ideally every endpoint creation should register the endpoint in a telemetry system for http health check. 




curl http://app:8080/W6oDZwzEKnOJQ