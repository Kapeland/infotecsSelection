package wallet

import (
	myDB "infotecsSelection/internal/db/sqlite"
	tp "infotecsSelection/internal/types"
)

const initBalance float64 = 100.0

func CreateWallet(walletID string) tp.Wallet {
	db := myDB.LaunchDB()
	defer myDB.CloseDB(db)
	myDB.AddWallet(walletID, initBalance, db)
	return tp.Wallet{walletID, initBalance}
}

// If error returns empty wallet.
// If Error then should be ErrNoRows
func CheckWallet(walletID string) (tp.Wallet, error) {
	db := myDB.LaunchDB()
	defer myDB.CloseDB(db)
	balance, err := myDB.FindWallet(walletID, db)
	if err != nil {
		return tp.Wallet{}, err
	}
	return tp.Wallet{walletID, balance}, nil
}

func UpdateWallet(wlt tp.Wallet) error {
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

func GetInAndOutOp(UUID string) ([]tp.Operation, error) {
	db := myDB.LaunchDB()
	defer myDB.CloseDB(db)

	return myDB.FindInAndOutOp(UUID, db)

}
