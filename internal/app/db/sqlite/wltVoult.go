package sqlite

import (
	"database/sql"
	"fmt"
	tp "infotecsSelection/internal/types"
	"log"
	"time"
)

type WltVoult struct {
	db *DB
}

func (wltDb *WltVoult) AddWallet(walletUUID string, balance float64) error {
	_, err := wltDb.db.db.Exec("insert into wallets (id, balance) values ($1, $2)",
		walletUUID, balance)
	if err != nil {
		return fmt.Errorf("can't add wallet to DB: %w", err)
	}
	return nil
}

func (wltDb *WltVoult) FindWallet(walletUUID string) (tp.Wallet, error) {
	row := wltDb.db.db.QueryRow("select * from wallets where id = $1", walletUUID)
	wlt := tp.Wallet{}
	err := row.Scan(&wlt.Id, &wlt.Balance)
	if err == sql.ErrNoRows {
		return tp.Wallet{}, fmt.Errorf("no wallet with UUID=%s in DB: %w", walletUUID, err)
	}
	if err != nil {
		return tp.Wallet{}, fmt.Errorf("error during searching wallet with UUID=%s in DB: %w", walletUUID, err)
	}
	return wlt, nil
}

func (wltDb *WltVoult) UpdateWallet(walletUUID string, balance float64) error {
	_, err := wltDb.db.db.Exec("update wallets set balance = $1 where id = $2", balance, walletUUID)
	if err != nil {
		return fmt.Errorf("can't update wallet balance in DB: %w", err)
	}
	return nil
}

func (wltDb *WltVoult) FillOperationLog(fromUUID, toUUID string, amount float64) error {
	_, err := wltDb.db.db.Exec("insert into op_log (fromID, toID, amount, time) values ($1, $2, $3, $4)",
		fromUUID, toUUID, amount, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("can't add operation info to DB: %w", err)
	}
	return nil
}

func (wltDb *WltVoult) FindInAndOutOp(walletUUID string) ([]tp.Operation, error) {
	rows, err := wltDb.db.db.Query("select * from op_log where op_log.fromID = $1 OR op_log.toID = $1", walletUUID)
	if err == sql.ErrNoRows {
		return []tp.Operation{}, fmt.Errorf("no operations with wallet with UUID=%s in DB: %w", walletUUID, err)
	}
	if err != nil {
		return []tp.Operation{}, fmt.Errorf("error during searching operations with UUID=%s, %w", walletUUID, err)
	}
	defer rows.Close()
	operations := []tp.Operation{}
	tmpStrTime := ""
	for rows.Next() {
		op := tp.Operation{}
		err := rows.Scan(&op.From, &op.To, &op.Amount, &tmpStrTime)
		if err != nil {
			log.Printf("error during finding operations. UUID=%s, %w", walletUUID, err)
			continue
		}
		op.Time, _ = time.Parse(time.RFC3339, tmpStrTime)
		operations = append(operations, op)
	}
	return operations, nil
}
