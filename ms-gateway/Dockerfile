FROM golang:1.27 AS build
WORKDIR /app
COPY go.work go.work.sum ./
COPY lb-core ./lb-core
COPY lb-core/go.sum ./lb-core/go.sum
COPY ms-gateway ./ms-gateway
WORKDIR /app/ms-gateway
RUN GOWORK=off go build -o /gateway ./cmd/application

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /gateway /app/gateway
COPY ms-gateway/config ./config
EXPOSE 8080
ENTRYPOINT ["/app/gateway"]
