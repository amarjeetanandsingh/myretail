FROM golang:1.14.6-alpine

#set go path
RUN export GOPATH=/go
RUN export PATH=$GOPATH:$PATH

# install required tools
RUN apk add git bash

# clone code
WORKDIR /go/src
RUN git clone https://github.com/amarjeetanandsingh/myretail.git
WORKDIR myretail

# Download dependencies
RUN go mod download

RUN GOOS=linux go build -o myretail .


# Run
ENTRYPOINT ["./myretail"]
