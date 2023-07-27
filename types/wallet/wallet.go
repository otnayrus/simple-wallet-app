package types

import (
	"database/sql"
)

type WalletService interface {
	Initialize(InitializeRequest) (InitializeResponse, error)
}

type WalletRepository interface {
	Create(Wallet) error
}

type WalletStatus int

const (
	StatusInactive WalletStatus = iota
	StatusNewlyCreated
	StatusActive
)

type (
	Wallet struct {
		ID        string       `db:"id"`
		OwnedBy   string       `db:"owned_by"`
		Token     string       `db:"token"`
		Status    int          `db:"status"`
		UpdatedAt sql.NullTime `db:"updated_at"`
		Balance   int          `db:"balance"`
	}

	InitializeRequest struct {
		CustomerID string `form:"customer_xid"`
	}

	InitializeResponse struct {
		Token string `json:"token"`
	}
)
