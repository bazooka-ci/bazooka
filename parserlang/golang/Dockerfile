FROM golang:1.4

ADD . /go/src/github.com/haklop/bazooka/parserlang/golang
RUN go get -d -v ./...
RUN go install -v github.com/haklop/bazooka/parserlang/golang

CMD ["/go/bin/golang"]
