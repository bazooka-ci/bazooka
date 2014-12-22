FROM golang:1.4

ADD . /go/src/github.com/haklop/bazooka/parserlang/python
RUN go get -d -v ./...
RUN go install -v github.com/haklop/bazooka/parserlang/python

CMD ["/go/bin/python"]
