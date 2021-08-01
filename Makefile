.PHONY: storage frontend client config-gen clean

all: storage frontend client 

storage:
	go build -o storage cmd/storage/main.go

frontend-interface:
	go build -o frontend-interface cmd/frontend-interface/main.go

frontend:
	go build -o frontend cmd/frontend/main.go

client-interface:
	go build -o client-interface cmd/client-interface/main.go

client:
	go build -o client cmd/client/main.go

clean:
	rm storage frontend client mem || true

images: 
	docker build -t kofeebrian/grpc-kvs-storage:${tag} -f cmd/storage/Dockerfile . & 
	docker build -t kofeebrian/grpc-kvs-frontend:${tag} -f cmd/frontend/Dockerfile . &
	docker build -t kofeebrian/grpc-kvs-client:${tag} -f cmd/client/Dockerfile .

push: 
	docker push kofeebrian/grpc-kvs-storage:${tag} &
	docker push kofeebrian/grpc-kvs-frontend:${tag} &
	docker push kofeebrian/grpc-kvs-client:${tag}

gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:.

clean-gen:
	rm -rf pb/*.go
