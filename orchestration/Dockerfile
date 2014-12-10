FROM golang:1.3.3

ADD . /go/src/github.com/haklop/bazooka/orchestration
RUN go get -d -v ./...
RUN go install -v github.com/haklop/bazooka/orchestration

CMD ["/go/bin/orchestration"]
