package handlers

import (
	"encoding/json"
	wlt "infotecsSelection/cmd/httpServer/internal/wallet"
	myURL "infotecsSelection/internal/url"
	myUUID "infotecsSelection/internal/uuid"
	"log"
	"net/http"
)

func CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		tmpWlt := wlt.CreateWallet(myUUID.CreateUUID())
		jsonData, err := json.Marshal(tmpWlt)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(jsonData)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest) //TODO maybe StatusInternalError?
			log.Println(err)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
func WalletInfoHandler(w http.ResponseWriter, r *http.Request) {
	recvURL := myURL.ParseURL(r.URL)

	if recvURL[0] != "" { // Значит что-то есть после домена и перед API, что не правильно
		w.WriteHeader(http.StatusBadRequest)
	}
	recvURL = recvURL[1:]
	if len(recvURL) != 4 && len(recvURL) != 5 { // Тот случай, когда запрос очень длинный или короткий
		w.WriteHeader(http.StatusBadRequest)
	}

	if len(recvURL) == 4 { // Только информация о кошельке
		switch r.Method {
		case "GET":
			wltID := recvURL[len(recvURL)-1]
			if err := myUUID.CheckUUID(wltID); err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				break
			}
			// Сейчас можно попытаться получить указанный кошелёк
			tmpWlt, err := wlt.CheckWallet(wltID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				break
			}
			jsonData, err := json.Marshal(tmpWlt)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(jsonData)

			if err != nil {
				w.WriteHeader(http.StatusNotFound) //TODO maybe StatusInternalError?
				log.Fatal(err)

			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

}
