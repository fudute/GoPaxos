# FROM golang:alpine AS builder

# RUN apk update && apk add --no-cache git 

# RUN mkdir /paxos

# WORKDIR /paxos

# COPY . .

# RUN go get -d -v

# # RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/paxos
# RUN go build -v -a -installsuffix cgo -o /go/bin/paxos

FROM ubuntu:18.04

COPY GoPaxos /go/bin/GoPaxos
VOLUME [ "/home/log" ]
ENTRYPOINT ["/go/bin/GoPaxos"]

EXPOSE 8080