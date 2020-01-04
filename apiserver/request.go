package main

import (
	"bytes"
	"fmt"

	"encoding/hex"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
)

// GenerateSigner - generates a new signer with a privatekey
// this must be modified to receive a private key by user in the params.
func GenerateSigner() *signing.Signer {
	ctx = signing.NewSecp256k1Context()
	privateKey := ctx.NewRandomPrivateKey()
	return signing.NewCryptoFactory(ctx).NewSigner(privateKey)
}

// CreateTransactionHandler - returns the
func CreateTransactionHeader(payloadHash string, addresses []string, signer *signing.Signer) ([]byte, error) {
	publicKey := signer.GetPublicKey().AsHex()
	rawTransactionHeader := transaction_pb2.TransactionHeader{
		SignerPublicKey:  publicKey,
		FamilyName:       "audit",
		FamilyVersion:    "1.0",
		Dependencies:     []string{},
		BatcherPublicKey: publicKey,
		Inputs:           addresses,
		Outputs:          addresses,
		PayloadSha512:    payloadHash,
	}
	fmt.Println(rawTransactionHeader)
	return proto.Marshal(&rawTransactionHeader)
}

// CreateTransactionBatch - Create the transaction batch
func CreateTransactionBatch(transactionHeader, payload []byte, signer *signing.Signer) ([]byte, error) {

	// Start transaction serialization
	transaction := transaction_pb2.Transaction{
		Header:          transactionHeader,
		HeaderSignature: hex.EncodeToString(signer.Sign(transactionHeader)),
		Payload:         payload,
	}

	// Batch header
	batchHeader := batch_pb2.BatchHeader{
		SignerPublicKey: signer.GetPublicKey().AsHex(),
		TransactionIds:  []string{transaction.HeaderSignature},
	}
	batchHeaderBytes, _ := proto.Marshal(&batchHeader)

	// Create the batch of transactions
	batch := batch_pb2.Batch{
		Header:          batchHeaderBytes,
		Transactions:    []*transaction_pb2.Transaction{&transaction},
		HeaderSignature: hex.EncodeToString(signer.Sign(batchHeaderBytes)),
	}

	// Batch list
	return proto.Marshal(
		&batch_pb2.BatchList{Batches: []*batch_pb2.Batch{&batch}})
}

// requestServer
func requestServer(batchList []byte) (*http.Response, error) {
	return http.Post(
		"http://localhost:8008/batches",
		"application/octet-stream",
		bytes.NewBuffer(batchList),
	)
}
