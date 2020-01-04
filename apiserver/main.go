package main

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"

	"context"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"k8s.io/api/auditregistration/v1alpha1"

)

type AuditService interface {
	Save(string) (string, error)
}

type auditService struct{}

func (auditService) SaveAudit(s v1alpha1.AuditSinkList) (string, error) {
	return "", nil
}


type saveResponse struct {

}

type AuditEvent struct {
	RequestURI         string
	Verb               string
	Code               int32
	User               string
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
		req := request.(v1alpha1.AuditSinkList)
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
	//
	//b, _:= ioutil.ReadAll(r.Body)
	//fmt.Println(string(b))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	fmt.Println(request.Items)
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}