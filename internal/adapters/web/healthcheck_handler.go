package web

import "net/http"

type healthcheckResponse struct {
	Status string `json:"status"`
}

func HealthcheckHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		writeOK(w, &healthcheckResponse{"UP"})
	}
}
