package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	tp "infotecsSelection/internal/types"
	"log"
	"time"
)

func LaunchDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to the DB: %w", err)
	}
	return db, nil
}

func CloseDB(db *sql.DB) {
	db.Close()
}

func InitDB(db *sql.DB) error {
	qWallets := `CREATE TABLE IF NOT EXISTS wallets (
    id TEXT PRIMARY KEY ,
    balance FLOAT
)`
	_, err := db.Exec(qWallets)
	if err != nil {
		return fmt.Errorf("can't create table wallets: %w", err)
	}
	qLog := `CREATE TABLE IF NOT EXISTS op_log (
    fromID TEXT ,
    toID TEXT ,
    amount FLOAT ,
    time TEXT,
    FOREIGN KEY(fromID, toID) REFERENCES wallets(id, id)
)`
	_, err = db.Exec(qLog)
	if err != nil {
		return fmt.Errorf("can't create table op_log: %w", err)
	}
	qPragma := `PRAGMA foreign_keys = ON`
	_, err = db.Exec(qPragma)
	if err != nil {
		return fmt.Errorf("can't use PRAGMA foreign_keys = ON: %w", err)
	}
	return nil
}

func AddWallet(walletUUID string, balance float64, db *sql.DB) error {
	_, err := db.Exec("insert into wallets (id, balance) values ($1, $2)",
		walletUUID, balance)
	if err != nil {
		return fmt.Errorf("can't add wallet to DB: %w", err)
	}
	return nil
}

func FindWallet(walletUUID string, db *sql.DB) (tp.Wallet, error) {
	row := db.QueryRow("select * from wallets where id = $1", walletUUID)
	wlt := tp.Wallet{}
	err := row.Scan(&wlt.Id, &wlt.Balance)
	if err == sql.ErrNoRows {
		return tp.Wallet{}, nil
	}
	if err != nil {
		return tp.Wallet{}, fmt.Errorf("error during searching wallet with UUID=%s in DB: %w", walletUUID, err)
	}
	return wlt, nil
}

func UpdateWalletDB(walletUUID string, balance float64, db *sql.DB) error {
	_, err := db.Exec("update wallets set balance = $1 where id = $2", balance, walletUUID)
	if err != nil {
		return fmt.Errorf("can't update wallet balance in DB: %w", err)
	}
	return nil
}

func FillOperationLog(fromUUID, toUUID string, amount float64, db *sql.DB) error {
	_, err := db.Exec("insert into op_log (fromID, toID, amount, time) values ($1, $2, $3, $4)",
		fromUUID, toUUID, amount, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("can't add operation info to DB: %w", err)
	}
	return nil
}

func FindInAndOutOp(UUID string, db *sql.DB) ([]tp.Operation, error) {
	rows, err := db.Query("select * from op_log where op_log.fromID = $1 OR op_log.toID = $1", UUID)
	if err == sql.ErrNoRows {
		return []tp.Operation{}, nil
	}
	if err != nil {
		return []tp.Operation{}, fmt.Errorf("error during searching operations with UUID=%s, %w", UUID, err)
	}
	defer rows.Close()
	operations := []tp.Operation{}
	tmpStrTime := ""
	for rows.Next() {
		op := tp.Operation{}
		err := rows.Scan(&op.From, &op.To, &op.Amount, &tmpStrTime)
		if err != nil {
			log.Printf("error during finding operations. UUID=%s, %w", UUID, err)
			continue
		}
		op.Time, _ = time.Parse(time.RFC3339, tmpStrTime)
		operations = append(operations, op)
	}
	return operations, nil
}
