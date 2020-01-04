package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"context"
	"strings"

	"encoding/json"
	"crypto/sha512"
	"encoding/hex"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"

	"github.com/hyperledger/sawtooth-sdk-go/signing"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	_ "github.com/hyperledger/sawtooth-sdk-go/signing"

)

var ctx signing.Context

type AuditService interface {
	Save(string) (string, error)
}

type auditService struct{}

func (auditService) SaveAudit(saveRequest) (string, error) {
	return "", nil
}

type saveResponse struct {
}

type User struct {
	Username string
}

type AuditEvent struct {
	RequestURI         string
	Verb               string
	Code               int32
	User               User
	ImpersonatedUser   string
	ImpersonatedGroups string
	Resource           string
	Namespace          string
	RequestObject      bool
	ResponseObject     bool
	AuthorizeDecision  string

	AdmissionWebhookMutationAnnotations map[string]string
	AdmissionWebhookPatchAnnotations    map[string]string
}

type saveRequest struct {
	Items []AuditEvent
}

func makeSaveEndpoint(svc auditService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(saveRequest)
		svc.SaveAudit(req)
		return saveResponse{}, nil
	}
}

func main() {
	auditSvc := auditService{}

	auditHandler := httptransport.NewServer(
		makeSaveEndpoint(auditSvc),
		decodeSaveRequest,
		encodeResponse,
	)

	http.Handle("/", auditHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateSigner() *signing.Signer {
	ctx = signing.NewSecp256k1Context()
	privateKey := ctx.NewRandomPrivateKey()
	fmt.Println(privateKey)
	return signing.NewCryptoFactory(ctx).NewSigner(privateKey)
}

func decodeSaveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request = saveRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	// We must save the key for each user and reuse this, so it is possible to
	// create permissions, may be etcd(?)
	signer := generateSigner()
	payload := []byte("teste")

	payloadHash := encodePayload(payload)
	header, err := createTransactionHeader(payloadHash, signer)
	if err != nil {
		fmt.Println(err)
	}

	transactionBatch, err := createTransactionBatch(header, payload, signer)
	if err != nil {
		fmt.Println(err)
	}

	requestServer(transactionBatch)

	for _, i := range request.Items {
		fmt.Println(i.RequestURI, i.Verb, i.User.Username)
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodePayload(payload []byte) string {
	hashHandler := sha512.New()
	hashHandler.Write(payload)
	return strings.ToLower(hex.EncodeToString(hashHandler.Sum(nil)))
}

func createTransactionHeader(payloadHash string, signer *signing.Signer) ([]byte, error) {
	publicKey := signer.GetPublicKey().AsHex()
	fmt.Println(publicKey)
	rawTransactionHeader := transaction_pb2.TransactionHeader{
		SignerPublicKey:  publicKey,
		FamilyName:       "audit",
		FamilyVersion:    "1.0",
		Dependencies:     []string{},
		BatcherPublicKey: publicKey ,
		Inputs:           []string{"1cf1266e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7"},  // TODO - fix this with correct namespace
		Outputs:          []string{"1cf1266e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7"},
		PayloadSha512:    payloadHash,
	}
	return proto.Marshal(&rawTransactionHeader)
}

func createTransactionBatch(transactionHeader []byte, payload []byte, signer *signing.Signer) ([]byte, error) {
	signature := hex.EncodeToString(signer.Sign(transactionHeader))

	transaction := transaction_pb2.Transaction{
		Header:          transactionHeader,
		HeaderSignature: signature,
		Payload:         payload,
	}


	batchatchHeader := batch_pb2.BatchHeader{
		SignerPublicKey: signer.GetPublicKey().AsHex(),
		TransactionIds:  []string{transaction.HeaderSignature},
	}

	batchHeaderBytes, _ := proto.Marshal(&batchatchHeader)
	signatureBatch := hex.EncodeToString(signer.Sign(batchHeaderBytes))
	batch := batch_pb2.Batch{
		Header:          batchHeaderBytes,
		Transactions:    []*transaction_pb2.Transaction{&transaction},
		HeaderSignature: signatureBatch,
	}

	rawBatchList := batch_pb2.BatchList{
		Batches: []*batch_pb2.Batch{&batch},
	}

	return proto.Marshal(&rawBatchList)
}

func requestServer(batchList []byte) {
	response, err := http.Post(
		"http://localhost:8008/batches",
		"application/octet-stream",
		bytes.NewBuffer(batchList),
	)
	fmt.Println(response, err)
}