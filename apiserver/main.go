package main

import (
	"log"
	"context"
	"net/http"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type AuditService interface {
	Save(string) (string, error)
}

type auditService struct{}

func (auditService) SaveAudit(saveRequest saveRequest) (int, error) {
	signer := GenerateSigner()
	var auditEvent AuditEvent = saveRequest.Items[0]

	payload, err := json.Marshal(auditEvent)
	if err != nil {
		return 0, err
	}

	addresses := []string{getAddress(auditEvent.RequestURI)}
	header, err := CreateTransactionHeader(HexDigest(payload), addresses, signer)
	if err != nil {
		return 0, err
	}

	transactionBatch, err := CreateTransactionBatch(header, payload, signer)
	if err != nil {
		return 0, err
	}

	response, err := requestServer(transactionBatch)
	if err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}

// Response
type saveResponse struct {
	Status int `json:"status"`
}

// Request
type User struct {
	Username string
}

type AuditEvent struct {
	RequestURI         string
	Verb               string
	Code               int32
	User               User
	Resource           string
	Namespace          string

	AdmissionWebhookMutationAnnotations map[string]string
	AdmissionWebhookPatchAnnotations    map[string]string
}

type saveRequest struct {
	Items []AuditEvent
}

func makeSaveEndpoint(svc auditService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(saveRequest)
		status, _ := svc.SaveAudit(req)
		return saveResponse{Status: status}, nil
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

func decodeSaveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request = saveRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
