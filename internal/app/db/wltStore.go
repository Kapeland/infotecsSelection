package db

import (
	tp "infotecsSelection/internal/types"
)

type WalletStore interface {
	AddWallet(walletUUID string, balance float64) error
	FindWallet(walletUUID string) (tp.Wallet, error)
	UpdateWallet(walletUUID string, balance float64) error
	FillOperationLog(fromUUID, toUUID string, amount float64) error
	FindInAndOutOp(UUID string) ([]tp.Operation, error)
}
