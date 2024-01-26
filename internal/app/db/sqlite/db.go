package sqlite

import (
	"database/sql"
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

func New(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}
