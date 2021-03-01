package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	admissionV1 "k8s.io/api/admission/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not read request body: %s", err.Error()), http.StatusBadRequest)
			return
		}

		admissionReviewRequest := admissionV1.AdmissionReview{}
		if err := json.Unmarshal(body, &admissionReviewRequest); err != nil {
			http.Error(w, "Body does not contain a valid AdmissionReview object", http.StatusBadRequest)
			return
		}

		admissionReviewResponse := admissionV1.AdmissionReview{
			TypeMeta: metaV1.TypeMeta{
				APIVersion: "admission.k8s.io/v1",
				Kind:       "AdmissionReview",
			},
			Response: &admissionV1.AdmissionResponse{
				UID:     admissionReviewRequest.Request.UID,
				Allowed: false,
				Result: &metaV1.Status{
					Message: "Pod is not allowed",
				},
			},
		}

		resp, err := json.Marshal(&admissionReviewResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Can't encode response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(resp)
		if err != nil {
			msg := fmt.Sprintf("Can not write response: %s", err.Error())
			log.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

	})

	err := http.ListenAndServeTLS(":443", "/tmp/tls-certs/tls.crt", "/tmp/tls-certs/tls.key", r)
	if err != nil {
		panic(err.Error())
	}
}
