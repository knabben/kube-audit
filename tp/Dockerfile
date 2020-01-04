FROM audit-base:latest

WORKDIR /go/src/github.com/knabben/kube-audit/tp

ENV GO111MODULE on
RUN apt-get install -y libzmq3-dev

COPY tp/ .

RUN go mod tidy ; exit 0
RUN cd /go/pkg/mod/github.com/hyperledger/sawtooth-sdk-go@v0.1.3 && \
chmod +x gen.sh && chmod +x protogen && ./gen.sh

RUN go build -o tp main.go
CMD ["./tp"]
