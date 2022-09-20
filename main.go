package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"k8s.io/api/admission/v1beta1"
)

func main() {
	tlsConfig, err := tlsConfig()
	if err != nil {
		panic(err)
	}
	listener, err := tls.Listen("tcp", ":8443", tlsConfig)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:              listener.Addr().String(),
		Handler:           handler(),
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Minute,
	}

	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func handler() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/validate", validateHandler())
	return router
}

func validateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		admissionRequest, err := decodeAdmissionRequest(req)
		if err != nil {
			http.Error(w, "DECODE_ERROR", http.StatusBadRequest)
			return
		}

		response := v1beta1.AdmissionReview{
			TypeMeta: admissionRequest.TypeMeta,
			Response: &v1beta1.AdmissionResponse{
				UID:     admissionRequest.Request.UID,
				Allowed: true,
			},
		}
		data, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "ENCODE_ERROR", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})
}

func decodeAdmissionRequest(r *http.Request) (v1beta1.AdmissionReview, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return v1beta1.AdmissionReview{}, err
	}
	var review v1beta1.AdmissionReview
	err = json.Unmarshal(body, &review)
	if err != nil {
		return v1beta1.AdmissionReview{}, err
	}
	return review, nil
}

func tlsConfig() (*tls.Config, error) {
	fnGetCertificate, err := tlsGetCertificate()
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		GetCertificate:   fnGetCertificate,
		ClientAuth:       tls.NoClientCert,
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		Renegotiation:          tls.RenegotiateNever,
		SessionTicketsDisabled: true,
	}, nil
}

type tlsGetCertificateFn func(*tls.ClientHelloInfo) (*tls.Certificate, error)

func tlsGetCertificate() (tlsGetCertificateFn, error) {
	const certificatePath = "./dev/certs/server.pem"
	const keyPath = "./dev/certs/server.key"

	cert, err := tls.LoadX509KeyPair(certificatePath, keyPath)
	if err != nil {
		return nil, err
	}
	return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return &cert, nil
	}, nil
}
