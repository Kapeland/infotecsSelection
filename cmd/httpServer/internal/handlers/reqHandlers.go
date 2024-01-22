package handlers

import (
	"fmt"
	"net/http"
)

func CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		_, err := fmt.Fprint(w, "Create wallet.")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func WalletInfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		_, err := fmt.Fprint(w, "Update wallet.")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
