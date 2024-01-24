package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	tp "infotecsSelection/internal/types"
	"log"
	"time"
)

const dbPath = "././identifier.sqlite?parseTime=true"

type wallet struct {
	id      string
	balance float64
}

func LaunchDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	return db
}

func CloseDB(db *sql.DB) {
	db.Close()
}

func PrintDB(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM wallets")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	wallets := []wallet{}

	for rows.Next() {
		tmpWallet := wallet{}
		err = rows.Scan(&tmpWallet.id, &tmpWallet.balance)
		if err != nil {
			fmt.Println(err)
			continue
		}
		wallets = append(wallets, tmpWallet)
	}
	for _, w := range wallets {
		fmt.Println(w.id, w.balance)
	}
}

func AddWallet(walletUUID string, balance float64, db *sql.DB) {
	_, err := db.Exec("insert into wallets (id, balance) values ($1, $2)",
		walletUUID, balance)
	if err != nil {
		panic(err)
	}

}

func FindWallet(walletUUID string, db *sql.DB) (float64, error) {
	row := db.QueryRow("select * from wallets where id = $1", walletUUID)
	wlt := wallet{}
	if err := row.Scan(&wlt.id, &wlt.balance); err != nil {
		return 0.0, err
	}
	return wlt.balance, nil
}

func UpdateWalletDB(walletUUID string, balance float64, db *sql.DB) error {
	_, err := db.Exec("update wallets set balance = $1 where id = $2", balance, walletUUID)
	if err != nil {
		return err
	}
	return nil
}

func FillOperationLog(fromUUID, toUUID string, amount float64, db *sql.DB) error {
	_, err := db.Exec("insert into op_log (fromID, toID, amount, time) values ($1, $2, $3, $4)",
		fromUUID, toUUID, amount, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func FindOutgoindOp(fromUUID string, db *sql.DB) ([]tp.Operation, error) {
	rows, err := db.Query("select * from op_log where op_log.fromID = $1", fromUUID)
	if err != nil {
		return []tp.Operation{}, err
	}
	defer rows.Close()
	operations := []tp.Operation{}
	tmpStrTime := ""
	for rows.Next() {
		op := tp.Operation{}
		err := rows.Scan(&op.From, &op.To, &op.Amount, &tmpStrTime)
		if err != nil {
			log.Println(err)
			continue
		}
		op.Time, _ = time.Parse(time.RFC3339, tmpStrTime)
		operations = append(operations, op)
	}
	return operations, nil
}
