FROM golang:1.3.3

ADD . /go/src/github.com/haklop/bazooka/parserlang/java
RUN go get -d -v ./...
RUN go install -v github.com/haklop/bazooka/parserlang/java

CMD ["/go/bin/java"]
