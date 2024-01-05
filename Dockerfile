FROM golang:1.21.5-bullseye AS build

RUN apt-get update && apt-get install -y git

WORKDIR /app

RUN echo orders-service

RUN git clone https://github.com/akshay0074700747/orders-service-grpc.git .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o bin/order-service

COPY /cmd/.env /app/cmd/bin/

FROM busybox:latest

WORKDIR /order-service

COPY --from=build /app/cmd/bin/order-service .

COPY --from=build /app/cmd/bin/.env .

EXPOSE 50003

CMD ["./order-service"]