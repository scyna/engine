# Build stage
FROM golang:1.19.3-alpine3.16 AS build-env
ENV GO111MODULE=on
WORKDIR /workspace/engine
COPY . .
RUN go mod tidy 
RUN go build -o ./application .

# Runtime stage
FROM alpine:3.17
WORKDIR /app
COPY --from=build-env /workspace/engine/application /engine
COPY ./.cert/ ./cert/
ENV MANAGER_PORT=6081 \
    PROXY_PORT=6060 \
    GATEWAY_PORT=6443 \
    NATS_URL=localhost \
    NATS_USERNAME="" \
    NATS_PASSWORD="" \
    DB_HOST="" \
    DB_USERNAME="" \
    DB_PASSWORD="" \
    DB_LOCATION="" \
    CERTIFICATE_FILE="/cert/localhost.crt" \
    CERTIFICATE_KEY="/cert/localhost.key" \
    CERTIFICATE_ENABLE="false"

EXPOSE $MANAGER_PORT $PROXY_PORT $GATEWAY_PORT
CMD ["/engine", \
    "--manager_port=${MANAGER_PORT}", \
    "--proxy_port=${PROXY_PORT}", \
    "--gateway_port=${GATEWAY_PORT}", \
    "--nats_url=${NATS_URL}", \
    "--nats_username=${NATS_USERNAME}", \
    "--nats_password=${NATS_PASSWORD}", \
    "--db_host=${DB_HOST}", \
    "--db_username=${DB_USERNAME}", \
    "--db_password=${DB_PASSWORD}", \
    "--db_location=${DB_LOCATION}", \
    "--certificateFile=${CERTIFICATE_FILE}", \
    "--certificateKey=${CERTIFICATE_KEY}", \
    "--certificateEnable=${CERTIFICATE_ENABLE}"
    ]
