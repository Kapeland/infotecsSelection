package handlers

import (
	"encoding/json"
	"fmt"
	wlt "infotecsSelection/cmd/httpServer/internal/wallet"
	intl "infotecsSelection/internal"
	"log"
	"net/http"
)

func CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		tmpWlt := wlt.CreateWallet(intl.GetUUID())
		jsonData, err := json.Marshal(tmpWlt)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(jsonData)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatal(err)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusBadRequest)
	}
}
