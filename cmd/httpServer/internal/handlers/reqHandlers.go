package handlers

import (
	"encoding/json"
	wlt "infotecsSelection/cmd/httpServer/internal/wallet"
	tp "infotecsSelection/internal/types"
	myURL "infotecsSelection/internal/url"
	myUUID "infotecsSelection/internal/uuid"
	"log"
	"net/http"
)

const (
	headerKey = "Content-Type"
	headerVal = "application/json; charset=utf-8"
)

func CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		tmpWlt := wlt.CreateWallet(myUUID.CreateUUID())
		jsonData, err := json.Marshal(tmpWlt)
		w.Header().Set(headerKey, headerVal)
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

	// Только информация о кошельке
	if len(recvURL) == 4 {
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
			w.Header().Set(headerKey, headerVal)
			w.Write(jsonData)

			if err != nil {
				w.WriteHeader(http.StatusNotFound) //TODO maybe StatusInternalError?
				log.Fatal(err)

			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	//history and send
	if len(recvURL) == 5 {
		switch r.Method {
		case "GET":
			if recvURL[len(recvURL)-1] != "history" {
				w.WriteHeader(http.StatusNotFound)
				log.Println("GET request but not /history endpoint.")
				break
			}
			reqWltID := recvURL[len(recvURL)-2]
			if err := myUUID.CheckUUID(reqWltID); err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				break
			}
			// Сейчас можно попытаться получить запрашиваемый кошелёк
			_, err := wlt.CheckWallet(reqWltID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				break
			}

			// Теперь пытаемся получить историю операций

			outgoinOp, err := wlt.GetOutgoingOp(reqWltID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				break
			}
			incomingOp, err := wlt.GetIncomingOp(reqWltID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				break
			}

			jsonData, err := json.Marshal(append(outgoinOp, incomingOp...))
			w.Header().Set(headerKey, headerVal)
			w.Write(jsonData)

			if err != nil {
				w.WriteHeader(http.StatusNotFound) //TODO maybe StatusInternalError?
				log.Println(err)
				break
			}
		case "POST":
			if recvURL[len(recvURL)-1] != "send" {
				w.WriteHeader(http.StatusNotFound)
				log.Println("POST request but not /send endpoint.")
				break
			}
			outgoingWltID := recvURL[len(recvURL)-2]
			if err := myUUID.CheckUUID(outgoingWltID); err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				break
			}
			// Сейчас можно попытаться получить исходящий кошелёк
			outgoingWlt, err := wlt.CheckWallet(outgoingWltID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println(err)
				break
			}
			incomingWlt := tp.WltForSend{}

			if r.Body == nil {
				w.WriteHeader(http.StatusBadRequest) //TODO maybe StatusInternalError?
				log.Println("There is no body")
				break
			}

			err = json.NewDecoder(r.Body).Decode(&incomingWlt)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest) //TODO maybe StatusInternalError?
				log.Println(err)
				break
			}

			if err := myUUID.CheckUUID(incomingWlt.To); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				break
			}

			incomingWltFromDB, err := wlt.CheckWallet(incomingWlt.To)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println("Not found wallet which recievs money.", err)
				break
			}

			//Сейчас попытаемся совершить операцию

			//Тут возможны проблемы из-за потери точности
			if outgoingWlt.Balance < incomingWlt.Amount {
				w.WriteHeader(http.StatusBadRequest)
				log.Println("Not enough money.")
				break
			}

			incomingWltFromDB.Balance += incomingWlt.Amount

			if err := wlt.UpdateWallet(incomingWltFromDB); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				break
			}

			outgoingWlt.Balance -= incomingWlt.Amount

			if err := wlt.UpdateWallet(outgoingWlt); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				break
			}
			if err := wlt.RegisterOperation(outgoingWlt.Id, incomingWltFromDB.Id, incomingWlt.Amount); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				break
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

}
