package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "././identifier.sqlite"

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

func AddWallet(userUUID string, balance float64, db *sql.DB) {
	_, err := db.Exec("insert into wallets (id, balance) values ($1, $2)",
		userUUID, balance)
	if err != nil {
		panic(err)
	}

}
