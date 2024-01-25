package wallet

import (
	"database/sql"
	myDB "infotecsSelection/internal/db/sqlite"
	tp "infotecsSelection/internal/types"
	"log"
)

const initBalance float64 = 100.0
const dbPath = "././identifier.sqlite"

var db *sql.DB = nil

func init() {
	var err error
	db, err = myDB.LaunchDB(dbPath)
	if err != nil {
		log.Fatal("Can't launch db.")
	}
	if err := myDB.InitDB(db); err != nil {
		log.Fatal("Can't init db.")
	}
}

func CreateWallet(walletID string) tp.Wallet {
	myDB.AddWallet(walletID, initBalance, db)
	return tp.Wallet{walletID, initBalance}
}

// If error returns empty wallet.
// If Error then should be ErrNoRows
func CheckWallet(walletID string) (tp.Wallet, error) {
	wlt, err := myDB.FindWallet(walletID, db)
	if err != nil {
		return tp.Wallet{}, err
	}
	return wlt, nil
}

func UpdateWallet(wlt tp.Wallet) error {
	err := myDB.UpdateWalletDB(wlt.Id, wlt.Balance, db)
	if err != nil {
		return err
	}
	return nil
}

func RegisterOperation(fromUUID, toUUID string, amount float64) error {
	err := myDB.FillOperationLog(fromUUID, toUUID, amount, db)
	if err != nil {
		return err
	}
	return nil
}

func GetInAndOutOp(UUID string) ([]tp.Operation, error) {
	return myDB.FindInAndOutOp(UUID, db)
}
