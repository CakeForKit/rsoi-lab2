FROM golang:1.27 AS build
WORKDIR /app
COPY go.work go.work.sum ./
COPY lb-core ./lb-core
COPY lb-core/go.sum ./lb-core/go.sum
COPY ms-flight ./ms-flight
WORKDIR /app/ms-flight
RUN GOWORK=off go build -o /flight ./cmd/application

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /flight /app/flight
COPY ms-flight/config ./config
EXPOSE 8060
ENTRYPOINT ["/app/flight"]
