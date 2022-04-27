FROM golang:1.17-alpine AS build-env

ENV GO111MODULE=on

RUN apk add --no-cache git

WORKDIR  /workspace/engine

COPY ./gateway/ ./gateway
COPY ./manager/ ./manager
COPY ./proxy/ ./proxy
COPY ./engine.go ./engine.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod tidy

RUN go build -o ./application .

FROM alpine:3.11

WORKDIR /app

COPY --from=build-env /workspace/engine/application /application

EXPOSE 8080
EXPOSE 8081
EXPOSE 8443

ENV SECRET="123456"
ENV NATS_URL=localhost
ENV NATS_USERNAME=""
ENV NATS_PASSWORD=""
ENV DB_HOST=""
ENV DB_USERNAME=""
ENV DB_PASSWORD=""
ENV DB_LOCATION=""

CMD /application --nats_url ${NATS_URL} --nats_username ${NATS_USERNAME} --nats_password ${NATS_PASSWORD} --db_host ${DB_HOST} --db_username ${DB_USERNAME} --db_password ${DB_PASSWORD} --db_location ${DB_LOCATION} --secret ${SECRET}