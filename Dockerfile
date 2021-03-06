FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . ./
RUN go build -o main cmd/main/main.go

FROM golang:1.18-alpine as builder_auth

WORKDIR /auth

COPY . ./

RUN go build -o auth cmd/auth/authorisathion.go

FROM alpine:latest as exec

WORKDIR /cmd

RUN mkdir -p /opt/pics
RUN mkdir /cmd/configs
VOLUME ["/cmd/configs"]

ARG CONFIG
ENV CONFIG=${CONFIG}

COPY --from=builder /app/main ./
COPY --from=builder /app/${CONFIG} ./configs

ENTRYPOINT ["/cmd/main"]

FROM alpine:latest as exec_auth

WORKDIR /auth_cmd

RUN mkdir /auth_cmd/configs
VOLUME ["/auth_cmd/configs"]

ARG CONFIG
ENV CONFIG=${CONFIG}

COPY --from=builder_auth /auth/auth ./
COPY --from=builder /app/${CONFIG} ./configs

ENTRYPOINT ["/auth_cmd/auth"]
