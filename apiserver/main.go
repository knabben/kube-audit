package main

import (
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"context"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

)

type AuditService interface {
	Save(string) (string, error)
}

type auditService struct{}

func (auditService) SaveAudit(s saveRequest) (string, error) {
	return "", nil
}

type saveRequest struct {

}
type saveResponse struct {

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
	req := saveRequest{}
	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}