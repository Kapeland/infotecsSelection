package httpServer

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"infotecsSelection/internal/app/db"
	tp "infotecsSelection/internal/types"
	myUUID "infotecsSelection/internal/uuid"
	"log"
	"net/http"
)

const (
	headerKey  = "Content-Type"
	headerVal  = "application/json; charset=utf-8"
	wltApiPath = "/api/v1/wallet"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Method:", r.Method, "Path:", r.URL.Path)
		f(w, r)
	}
}

type server struct {
	router *mux.Router
	db     db.DB
}

func newServer(db db.DB) *server {
	s := &server{
		router: mux.NewRouter(),
		db:     db,
	}
	s.configureRouter()
	return s
}
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
func (s *server) configureRouter() {
	wltS := s.router.PathPrefix(wltApiPath).Subrouter()
	wltS.HandleFunc("", s.createWalletHandler()).Methods("POST")
	wltS.HandleFunc("/{walletId}", s.WalletInfoHandler()).Methods("GET")
	wltS.HandleFunc("/{walletId}/history", s.WalletHistoryHandler()).Methods("GET")
	wltS.HandleFunc("/{walletId}/send", s.WalletSendHandler()).Methods("POST")
}

func (s *server) createWalletHandler() http.HandlerFunc {
	return logging(func(w http.ResponseWriter, r *http.Request) {
		genUUID := myUUID.CreateUUID()
		var initBalance float64 = 100.0
		err := s.db.Wallet().AddWallet(genUUID, initBalance)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		tmpWlt := tp.Wallet{genUUID, initBalance}

		jsonData, err := json.Marshal(tmpWlt)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		w.Header().Set(headerKey, headerVal)
		_, err = w.Write(jsonData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

	})
}

func (s *server) WalletInfoHandler() http.HandlerFunc {
	return logging(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		wltID := vars["walletId"]
		if err := myUUID.CheckUUID(wltID); err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}
		// Сейчас можно попытаться получить указанный кошелёк
		tmpWlt, err := s.db.Wallet().FindWallet(wltID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}
		jsonData, err := json.Marshal(tmpWlt)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		w.Header().Set(headerKey, headerVal)
		_, err = w.Write(jsonData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	})
}

func (s *server) WalletHistoryHandler() http.HandlerFunc {
	return logging(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reqWltID := vars["walletId"]

		if err := myUUID.CheckUUID(reqWltID); err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}
		// Сейчас можно попытаться получить запрашиваемый кошелёк
		_, err := s.db.Wallet().FindWallet(reqWltID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}

		// Теперь пытаемся получить историю операций

		historyOfOp, err := s.db.Wallet().FindInAndOutOp(reqWltID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}

		jsonData, err := json.Marshal(historyOfOp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		w.Header().Set(headerKey, headerVal)
		_, err = w.Write(jsonData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	})
}

func (s *server) WalletSendHandler() http.HandlerFunc {
	return logging(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		outgoingWltID := vars["walletId"]

		if err := myUUID.CheckUUID(outgoingWltID); err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}
		// Сейчас можно попытаться получить исходящий кошелёк
		outgoingWlt, err := s.db.Wallet().FindWallet(outgoingWltID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}
		incomingWlt := tp.WltForSend{}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("There is no body in request")
		}

		err = json.NewDecoder(r.Body).Decode(&incomingWlt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}

		if err := myUUID.CheckUUID(incomingWlt.To); err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}

		incomingWltFromDB, err := s.db.Wallet().FindWallet(incomingWlt.To)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Not found wallet which recievs money.", err)
		}

		//Сейчас попытаемся совершить операцию

		if incomingWlt.Amount < 0.0 {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Negative amount.")
		}

		//Тут возможны проблемы из-за потери точности
		if outgoingWlt.Balance < incomingWlt.Amount {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Not enough money in wallet.")
		}

		incomingWltFromDB.Balance += incomingWlt.Amount

		if err := s.db.Wallet().UpdateWallet(incomingWltFromDB.Id, incomingWltFromDB.Balance); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		outgoingWlt.Balance -= incomingWlt.Amount

		if err := s.db.Wallet().UpdateWallet(outgoingWlt.Id, outgoingWlt.Balance); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		if err := s.db.Wallet().FillOperationLog(outgoingWlt.Id, incomingWltFromDB.Id, incomingWlt.Amount); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	})
}
