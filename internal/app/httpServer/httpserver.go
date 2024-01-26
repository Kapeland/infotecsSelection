package httpServer

import (
	"database/sql"
	"fmt"
	"infotecsSelection/internal/app/db/sqlite"
	"net/http"
)

const path = "././identifier.sqlite"

// Start ...
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlite.New(db)
	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, fmt.Errorf("can't open DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to the DB: %w", err)
	}
	return db, nil
}
