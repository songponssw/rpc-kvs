#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v .
RUN go build -o /go/bin/app -v cmd/tracing-server/main.go

#final stage
FROM alpine:latest
ARG config="storage_config.json"
ENV CONFIG=${config}
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/src/app/config/${CONFIG} /config/tracing_server_config.json
RUN ["chmod", "+x", "/app"]
ENTRYPOINT /app
LABEL Name=storage Version=0.0.1
EXPOSE 6666
