package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"infotecsSelection/internal/app/db"
)

type DB struct {
	db       *sql.DB
	wltVault *WltVoult
}

func (db *DB) Wallet() db.WalletStore {
	if db.wltVault != nil {
		return db.wltVault
	}

	db.wltVault = &WltVoult{
		db: db,
	}

	return db.wltVault
}

func New(db *sql.DB) (*DB, error) {
	if err := InitDB(db); err != nil {
		fmt.Errorf("can't init SQLite DB: %w", err)
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func InitDB(db *sql.DB) error {
	/*qPragma := `PRAGMA foreign_keys = ON;`
	_, err := db.Exec(qPragma)
	if err != nil {
		return fmt.Errorf("can't use PRAGMA foreign_keys = ON: %w", err)
	}
	*/
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
    FOREIGN KEY(fromID) REFERENCES wallets(id),
    FOREIGN KEY(toID) REFERENCES wallets(id)          
)`
	_, err = db.Exec(qLog)
	if err != nil {
		return fmt.Errorf("can't create table op_log: %w", err)
	}

	return nil
}
