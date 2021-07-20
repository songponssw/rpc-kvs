.PHONY: storage frontend client config-gen clean

all: storage frontend client

storage:
	go build -o storage cmd/storage/main.go

frontend:
	go build -o frontend cmd/frontend/main.go

client:
	go build -o client cmd/client/main.go

clean:
	rm storage frontend client mem || true
