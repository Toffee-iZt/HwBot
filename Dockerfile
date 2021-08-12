FROM golang:1.16-alpine

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

CMD [ "./bin/hwbot" ]
