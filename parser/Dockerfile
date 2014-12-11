FROM golang:1.4

ADD . /go/src/github.com/haklop/bazooka/parser
RUN go get -d -v ./...
RUN go install -v github.com/haklop/bazooka/parser

CMD ["/go/bin/parser"]
