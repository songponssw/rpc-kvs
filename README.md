# RPC-KVS Application

Key-Value Store Application Implement as Microservices application with golang

## Detail of Application

// TODO: Insert Application Architecture Image here

Application including by 3 components

1. **client** - _invoke function (gRPC) to frontend_
   1.1. **kvslib** - intercepter between client and frontend
2. **frontend** - _endpoint accept request (gRPC) from client and create another request (RPC) to storage_
3. **storage** - _handle request (RPC) from frontend_

Configuration read from JSON files in `config/*.json`

- _client_config.json_ - config for client
  - ClientID: define client identity
  - FrontEndAddr: define url to Frontend service
- _frontend_config_ - config for frontend
  - ClientAPIListenAddr: define listening port for client invoking
  - StorageAPIListenAddr: define url connect with storage service
- _storage_config_ - config for storage
  - StorageID: define storage identity
  - ~~StorageAdd: no use~~ TODO: remove
  - FrontEndAddr: define port listening to frontend requests
  - DiskPath: define path to store the data file

Folders Description
**cmd** - command
**config** - contain application configuration files
**k8s** - contain kubernetes version configuration files
**kn** - contain knative version configuration files
**monolith** - contain `Dockerfile` of monolith version application
**proto** - all proto files
**pb** - Generate from proto files

## Protobuf files

gRPC required files generated from `.proto` files with specific syntax. See more https://grpc.io/docs/languages/go/basics/

Example by using `protoc` from root directory of the project

```
protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:.
```

or, just use `Makefile`

```
# generate all go files
make gen
# clean up
make clean-gen
```

## Containerization to Docker images

_Required: Docker_
Write `Dockerfile` for each image. They are all golang code so `Dockerfile` is look alike

\***Caution**: Now, still do not know how to assign config inside microVM to config frontend more dynamically. For deploy in vHive cluster, service url, at least of frontend, should be config for vhive cluster here...
Example: storage url: `storage-service.default.svc.cluster.local`

To make it more flexible, I think vHive current version do not provide the convenient way to do this. I think using environment variable is the most efficient

```
# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v .

# change from "frontend" to the other for different image
RUN go build -o /go/bin/app -v cmd/frontend/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy executable file to root
COPY --from=builder /go/bin/app /app

# Copy configuration files to root
COPY --from=builder /go/src/app/config /config

# Set Entrypoint of the image
RUN ["chmod", "+x", "/app"]
ENTRYPOINT /app
LABEL Name=frontend Version=0.0.1

# Expose port that app running on
EXPOSE 50051
```

Build all application images

```
# Build Frontend Image
docker build -t <username>/frontend:<tag> -f cmd/frontend/Dockerfile .
# Build Storage Image
docker build -t <username>/storage:<tag> -f cmd/storage/Dockerfile .
# Build Client Image
docker build -t <username>/client:<tag> -f cmd/client/Dockerfile .
```

To use these images in vHive, one way is to push these images public repository and pull to vHive

To push the images to docker repository

```
# -a: push all tag
docker push <username>/frontend -a
docker push <username>/storage -a
docker push <username>/client -a
```

Now image is ready to use in vHive cluster

## Installation on vHive

_Required_

- vHive cluster running

To install vhive, read `SETUP_VHIVE.md`

To deploy image function in vHive cluster, able to use this template file

frontend.yaml

```
apiVersion: serving.knative.dev/v1
kind: Service
  metadata:
    namespace: default
spec:
  template:
    spec:
      containers:
        - image: crccheck/hello-world:latest # Stub image
          ports:
            - name: h2c # For GRPC support
              containerPort: 50051
          env:
            - name: GUEST_PORT # Port on which the firecracker-containerd container is accepting requests
              value: "50051"
            - name: GUEST_IMAGE # Container image to use for firecracker-containerd container
              value: "docker.io/<username>/frontend:latest"
```

**containers[\*].ports[\*].name**
Define what protocal using for communicate that can be "http1" for HTTP1 base communication or "h2c" for HTTP2 base

**containers[\*].env["name"=="GUEST_IMAGE"]**
Define what function image

**containers[\*].ports[\*].containerPort** - _queue-proxy_ will try to dial via this port of function microVM
**containers[\*].env["name"=="GUEST_PORT"]** - expose port of function microVM
These option **must** be the same value.

**!Caution** Do not sure is it hard code for port 50051. See in `vhive/cri/create_container.go`

Then apply the config files

```
kn service apply -f frontend.yaml
```

### To deploy application to vHive cluster use this command

```
# Create Persistant Volume to store data file
kubectl apply -f kn/grpc/psv.yaml

# Create Storage pods and service
kubectl apply -f kn/grpc/storage.yaml

# Deploy Frontend services as vHive function
kn service apply -f kn/grpc/frontend.yaml
```

### invoke with service

Suppose being a client from outside the cluster invoke the frontend services, use container running on host-network should be make sense

Example command to invoke with the frontend service
First config `config/client_config.json` key `FrontEndAddr` to url of the service

> Note: url MUST NOT contain "http://" just the domain-name only

```
# example
{
	"FrontEndAddr": "frontend.default.172.16.1.240.sslip.io"
}
```

Create container to invoke with the service

```
sudo ctr run -rm --mount type=bind,src=$PWD/config,dst=/config,options=rbind:ro --net-host docker.io/<username>/client test-client
```

# Appendix

To use function that uses HTTP1 Base e.g. RESTful API, RPC etc. vHive also support
First, config name of containerPort to "http1"

```
apiVersion: serving.knative.dev/v1
kind: Service
  metadata:
    namespace: default
spec:
  template:
    spec:
      containers:
        - image: crccheck/hello-world:latest # Stub image
          ports:
            - name: http1 # use HTTP1
              containerPort: 50051
          env:
            - name: GUEST_PORT # Port on which the firecracker-containerd container is accepting requests
              value: "50051"
            - name: GUEST_IMAGE # Container image to use for firecracker-containerd container
              value: "docker.io/<username>/frontend:latest"
```

1. **RPC** - Go lang only

How to invoke

```
curl -X CONNECT --url <url>/_goRPC_ \
	-d '{"method":"your rpc method","params":["args"],"id":<some number to represent your request>}'
```

2. **gRPC** - create RESETful gateway and convert request to RPC
   https://grpc-ecosystem.github.io/grpc-gateway/

> Written with [StackEdit](https://stackedit.io/).
