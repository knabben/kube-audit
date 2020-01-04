package main

import (
	processor2 "github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/knabben/kube-audit/tp/handler"
	"syscall"
)

func main() {
	endpoint := "tcp://localhost:4004"

	handler := &handler.AuditHandler{}
	processor := processor2.NewTransactionProcessor(endpoint)
	processor.AddHandler(handler)
	processor.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)

	processor.Start()
}