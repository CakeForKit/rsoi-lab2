FROM golang:1.27 AS build
WORKDIR /app
COPY go.work go.work.sum ./
COPY lb-core ./lb-core
COPY lb-core/go.sum ./lb-core/go.sum
COPY ms-bonus ./ms-bonus
WORKDIR /app/ms-bonus
RUN GOWORK=off go build -o /bonus ./cmd/application

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /bonus /app/bonus
COPY ms-bonus/config ./config
EXPOSE 8050
ENTRYPOINT ["/app/bonus"]
