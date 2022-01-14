##
## Build
##

FROM golang:1.16-alpine AS build

RUN apk add --no-cache --update make

WORKDIR /src

ADD go.mod go.sum ./
RUN go mod download

ADD Makefile ./

COPY main.go ./
COPY bot/ ./bot
COPY common/ ./common
COPY logger/ ./logger
COPY shttp/ ./shttp
COPY vkapi/ ./vkapi

RUN make

##
## Deploy
##

FROM alpine:latest

WORKDIR /app

COPY --from=build /src/bin ./

ENTRYPOINT [ "./hwbot" ]
