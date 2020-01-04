FROM golang:buster

WORKDIR /app
RUN apt-get update
RUN apt-get install -y python3 python3-grpcio python3-grpc-tools libssl-dev net-tools

RUN go get github.com/golang/mock/mockgen
RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/hyperledger/sawtooth-sdk-go

RUN cd /go/src/github.com/hyperledger/sawtooth-sdk-go && go generate

COPY apiserver/utils.go utils.go
COPY apiserver/request.go request.go
COPY apiserver/main.go main.go

RUN go get github.com/btcsuite/btcd/btcec
RUN go get github.com/go-kit/kit/endpoint
RUN go get github.com/go-logfmt/logfmt

RUN go build main.go request.go utils.go

CMD ["./main"]
