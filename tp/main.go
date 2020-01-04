package main

import (
	processor2 "github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/knabben/kube-audit/tp/handler"
	"os"
	"syscall"
)

var endpoint string

func init() {
	endpoint = os.Getenv("VALIDATOR_URL")
	if endpoint == "" {
		endpoint = "tcp://localhost:4004"
	}
}

func main() {
	handler := &handler.AuditHandler{}
	processor := processor2.NewTransactionProcessor(endpoint)
	processor.AddHandler(handler)
	processor.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)
	processor.Start()
}