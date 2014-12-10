FROM golang:1.3.3

ADD . /go/src/github.com/haklop/bazooka/server
RUN go get -d -v ./...
RUN go install -v github.com/haklop/bazooka/server

CMD ["/go/bin/server"]
EXPOSE 3000
