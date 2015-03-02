FROM golang:1.4

ADD . /go/src/github.com/bazooka-ci/bazooka
RUN go get -d -v ./...
RUN go install -v github.com/bazooka-ci/bazooka/parserlang/java

CMD ["/go/bin/java"]
