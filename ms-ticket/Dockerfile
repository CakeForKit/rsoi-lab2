FROM golang:1.27 AS build
WORKDIR /app
COPY go.work go.work.sum ./
COPY lb-core ./lb-core
COPY lb-core/go.sum ./lb-core/go.sum
COPY ms-ticket ./ms-ticket
WORKDIR /app/ms-ticket
RUN GOWORK=off go build -o /ticket ./cmd/application

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /ticket /app/ticket
COPY ms-ticket/config ./config
EXPOSE 8070
ENTRYPOINT ["/app/ticket"]
