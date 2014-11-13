FROM golang:1.3.3

VOLUME ["/bazooka"]

ADD . /go/src/bitbucket.org/bywan/bazooka-api/bazooka-server
RUN go get -d -v ./...
RUN go install -v bitbucket.org/bywan/bazooka-api/bazooka-server/server

CMD ["/go/bin/server"]
EXPOSE 3000
