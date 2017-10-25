FROM golang:1.9
ENV GOPATH=/go
COPY . /go/src/github.com/daemonl/go_robot
RUN go build -o /robot github.com/daemonl/go_robot/cmd/robot
ENTRYPOINT ["/robot"]
