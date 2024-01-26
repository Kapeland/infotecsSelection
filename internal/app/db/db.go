package db

type DB interface {
	Wallet() WalletStore
}
