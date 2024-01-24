package wallet

import myDB "infotecsSelection/internal/db/sqlite"

const initBalance float64 = 100.0

type Wallet struct {
	Id      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type WltForSend struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func CreateWallet(walletID string) Wallet {
	db := myDB.LaunchDB()
	defer myDB.CloseDB(db)
	myDB.AddWallet(walletID, initBalance, db)
	return Wallet{walletID, initBalance}
}

// If error returns empty wallet
func CheckWallet(walletID string) (Wallet, error) {
	db := myDB.LaunchDB()
	defer myDB.CloseDB(db)
	balance, err := myDB.FindWallet(walletID, db)
	if err != nil {
		return Wallet{}, err
	}
	return Wallet{walletID, balance}, nil
}

func UpdateWallet(wlt Wallet) error {
	db := myDB.LaunchDB()
	defer myDB.CloseDB(db)
	err := myDB.UpdateWalletDB(wlt.Id, wlt.Balance, db)
	if err != nil {
		return err
	}
	return nil
}

func RegisterOperation(fromUUID, toUUID string, amount float64) error {
	db := myDB.LaunchDB()
	defer myDB.CloseDB(db)
	err := myDB.FillOperationLog(fromUUID, toUUID, amount, db)
	if err != nil {
		return err
	}
	return nil
}
