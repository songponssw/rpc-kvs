.PHONY: storage frontend client config-gen clean

all: storage frontend client frontend-interface client-interface

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
	rm storage frontend client frontend-interface client-interface mem || true

gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:.

clean-gen:
	rm -rf pb/*.go
