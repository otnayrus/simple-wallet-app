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

	updateWalletStatusQuery = `
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

	setWalletBalanceByTokenQuery = `
		UPDATE wallets
		SET
			balance = $1,
			updated_at = $2
		WHERE
			token = $3;
	`

	createMutationQuery = `
		INSERT INTO mutations (id, reference_id, created_at, created_by, action, status, amount)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	getMutationListQuery = `
		SELECT id, reference_id, created_at, created_by, action, status, amount
		FROM mutations
		WHERE created_by = $1
		ORDER BY created_at DESC;
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
		updateWalletStatusQuery,
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

func (wr *walletRepository) Disable(token string) (types.Wallet, error) {
	var data types.Wallet
	err := wr.db.QueryRow(
		updateWalletStatusQuery,
		types.StatusInactive,
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

func (wr *walletRepository) SetWalletBalanceByToken(balance float64, token string) error {
	_, err := wr.db.Exec(
		setWalletBalanceByTokenQuery,
		balance,
		time.Now(),
		token,
	)

	return err
}

func (wr *walletRepository) CreateMutation(req types.Mutation) error {
	_, err := wr.db.Exec(
		createMutationQuery,
		req.ID,
		req.ReferenceID,
		req.CreatedAt,
		req.CreatedBy,
		req.Action,
		req.Status,
		req.Amount,
	)

	return err
}

func (wr *walletRepository) Mutate(req types.Mutation, expectedBalance float64, token string) error {
	tx, err := wr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = wr.CreateMutation(req)
	if err != nil {
		return err
	}

	err = wr.SetWalletBalanceByToken(expectedBalance, token)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (wr *walletRepository) ListMutation(ownerID string) ([]types.Mutation, error) {
	rows, err := wr.db.Query(getMutationListQuery, ownerID)
	defer rows.Close()

	var mutations []types.Mutation
	for rows.Next() {
		var mutation types.Mutation
		err := rows.Scan(
			&mutation.ID,
			&mutation.ReferenceID,
			&mutation.CreatedAt,
			&mutation.CreatedBy,
			&mutation.Action,
			&mutation.Status,
			&mutation.Amount,
		)
		if err != nil {
			return nil, err
		}

		mutations = append(mutations, mutation)
	}

	return mutations, err
}
