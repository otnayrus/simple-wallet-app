package repository

import (
	"database/sql"
	"time"

	types "github.com/otnayrus/simple-wallet-app/types/wallet"
)

type walletRepository struct {
	db *sql.DB
}

const (
	createWalletQuery = `
		INSERT INTO wallets (id, owned_by, token, status, updated_at, balance)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	enableWalletQuery = `
		UPDATE wallets
		SET
			status = $1,
			updated_at = $2
		WHERE
			token = $3
		RETURNING id, owned_by, token, status, updated_at, balance;
	`

	getWalletByTokenQuery = `
		SELECT id, owned_by, token, status, updated_at, balance
		FROM wallets
		WHERE token = $1;
	`
)

func NewWalletRepositiory(db *sql.DB) types.WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (wr *walletRepository) Create(req types.Wallet) error {
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

func (wr *walletRepository) Enable(token string) (types.Wallet, error) {
	var data types.Wallet
	err := wr.db.QueryRow(
		enableWalletQuery,
		types.StatusActive,
		time.Now(),
		token,
	).Scan(
		&data.ID,
		&data.OwnedBy,
		&data.Token,
		&data.Status,
		&data.UpdatedAt,
		&data.Balance,
	)

	return data, err
}

func (wr *walletRepository) GetByToken(token string) (types.Wallet, error) {
	var data types.Wallet
	err := wr.db.QueryRow(getWalletByTokenQuery, token).Scan(
		&data.ID,
		&data.OwnedBy,
		&data.Token,
		&data.Status,
		&data.UpdatedAt,
		&data.Balance,
	)

	return data, err
}
