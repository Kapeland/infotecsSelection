package wallet

import myDB "infotecsSelection/internal/db/sqlite"

const initBalance float64 = 100.0

type Wallet struct {
	Id      string  `json:"id"`
	Balance float64 `json:"balance"`
}

func CreateWallet(userUUID string) Wallet {
	db := myDB.LaunchDB()
	defer myDB.CloseDB(db)
	myDB.AddWallet(userUUID, initBalance, db)
	return Wallet{userUUID, initBalance}
}
