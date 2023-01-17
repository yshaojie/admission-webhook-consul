package memory_validating_webhook

import (
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme)
)

type WebhookServer struct {
	server *http.Server
}

func (webhookServer *WebhookServer) name(response http.ResponseWriter, request *http.Request) {
	var body []byte
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		body = data
	}

	ar := v1.AdmissionReview{}
	var admissionResponse *v1.AdmissionResponse
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		glog.Errorf("Can't decode body: %v", err)
		admissionResponse = &v1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		fmt.Println(request.URL.Path)
		if request.URL.Path == "/mutate" {
			admissionResponse = webhookServer.mutate(&ar)
		} else if request.URL.Path == "/validate" {
			admissionResponse = webhookServer.validate(&ar)
		}
	}
}

func (webhookServer WebhookServer) mutate(admissionReview *v1.AdmissionReview) *v1.AdmissionResponse {
	return nil
}

func (webhookServer *WebhookServer) validate(admissionReview *v1.AdmissionReview) *v1.AdmissionResponse {
	return nil
}
