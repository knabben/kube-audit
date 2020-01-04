package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"context"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

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

func decodeSaveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request = saveRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	for _, i := range request.Items {
		fmt.Println(i.RequestURI, i.Verb, i.User.Username)
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
