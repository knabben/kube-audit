FROM golang:buster
WORKDIR /app

ENV GO111MODULE on
COPY main.go main.go
RUN go build main.go

CMD ["./main"]
