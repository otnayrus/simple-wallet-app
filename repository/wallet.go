package repository

import (
	"database/sql"

	types "github.com/otnayrus/simple-wallet-app/types/wallet"
)

type walletRepositiory struct {
	db *sql.DB
}

const (
	createWalletQuery = `
		INSERT INTO wallets (id, owned_by, token, status, updated_at, balance)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
)

func NewWalletRepositiory(db *sql.DB) types.WalletRepository {
	return &walletRepositiory{
		db: db,
	}
}

func (wr *walletRepositiory) Create(req types.Wallet) error {
	_, err := wr.db.Exec(
		createWalletQuery,
		req.ID,
		req.OwnedBy,
		req.Token,
		req.Status,
		req.UpdatedAt,
		req.Balance,
	)
	return err
}
