package web

import (
	"encoding/json"
	"errors"
	"github.com/mikhail-karpov/url-shortener/internal/domain"
	"log"
	"net/http"
)

type errEnvelope struct {
	Error string `json:"error"`
}

func readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeOK(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, data)
}

func writeBadRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, &errEnvelope{Error: msg})
}

func writeErr(w http.ResponseWriter, err error) {

	if errors.Is(err, domain.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, &errEnvelope{Error: err.Error()})
	} else {
		log.Printf("unexpected error: %s\n", err)
		writeJSON(w, http.StatusInternalServerError, &errEnvelope{Error: "something went wrong"})
	}
}
