package handler

import (
	"encoding/json"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

type AuditEvent struct {
	RequestURI         string
	Verb               string
	Code               int32
	Resource           string
	Namespace          string

	AdmissionWebhookMutationAnnotations map[string]string
	AdmissionWebhookPatchAnnotations    map[string]string
}

type AuditHandler struct {}

func (*AuditHandler) FamilyName() string {
	return "audit"
}

func (*AuditHandler) FamilyVersions() []string{
	return []string{"1.0"}
}

func (*AuditHandler) Namespaces() []string {
	return []string{Namespace}
}

func (*AuditHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	var event AuditEvent

	data := request.GetPayload()
	json.Unmarshal(data, &event)
	addresses, err := context.SetState(map[string][]byte{
		getAddress(event.RequestURI): data,
	})
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return &processor.InternalError{Msg: "No addresses in set response"}
	}

	return nil
}