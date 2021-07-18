.PHONY: storage frontend client tracing-server config-gen clean docker-storage docker-frontend docker-tracing-server

all: storage frontend client tracing-server

storage:
	go build -o storage cmd/storage/main.go

frontend:
	go build -o frontend cmd/frontend/main.go

client:
	go build -o client cmd/client/main.go

tracing-server:
	go build -o tracing-server cmd/tracing-server/main.go

clean:
	rm storage frontend client tracing-server *".log" *"-Log.txt" 2> /dev/null || true

docker-storage:
	docker build --progress plain --build-arg config="storage_config.k8s.json" -t kofeebrian/kvs-storage:k8s -f cmd/storage/Dockerfile .

docker-frontend:
	docker build --progress plain --build-arg config="frontend_config.k8s.json" -t kofeebrian/kvs-frontend:k8s -f cmd/frontend/Dockerfile .

docker-tracing-server:
	docker build --progress plain --build-arg config="tracing_server_config.k8s.json" -t kofeebrian/kvs-tracing-server:k8s -f cmd/tracing-server/Dockerfile .