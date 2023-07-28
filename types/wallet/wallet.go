package types

import (
	"database/sql"
	"time"
)

type WalletService interface {
	Initialize(InitializeRequest) (InitializeResponse, error)
	Enable(EnableRequest) (EnableResponse, error)
	ViewBalance(ViewBalanceRequest) (ViewBalanceResponse, error)
}

type WalletRepository interface {
	Create(Wallet) error
	Enable(string) (Wallet, error)
	GetByToken(string) (Wallet, error)
}

type WalletStatus int

const (
	StatusInactive WalletStatus = iota
	StatusNewlyCreated
	StatusActive

	StatusStringEnabled  string = "enabled"
	StatusStringDisabled string = "disabled"
)

var (
	walletStatusToString = map[WalletStatus]string{
		StatusInactive: StatusStringDisabled,
		StatusActive:   StatusStringEnabled,
	}
)

func (w *Wallet) GetStatusString() string {
	return walletStatusToString[WalletStatus(w.Status)]
}

type (
	Wallet struct {
		ID        string       `db:"id"`
		OwnedBy   string       `db:"owned_by"`
		Token     string       `db:"token"`
		Status    int          `db:"status"`
		UpdatedAt sql.NullTime `db:"updated_at"`
		Balance   float64      `db:"balance"`
	}

	InitializeRequest struct {
		CustomerID string `form:"customer_xid"`
	}

	InitializeResponse struct {
		Token string `json:"token"`
	}

	WalletRequestHeader struct {
		Authorization string `header:"Authorization"`
	}

	EnableRequest struct {
		Token string
	}

	EnableResponse struct {
		ID        string    `json:"id"`
		OwnedBy   string    `json:"owned_by"`
		Status    string    `json:"status"`
		EnabledAt time.Time `json:"enabled_at"`
		Balance   float64   `json:"balance"`
	}

	ViewBalanceRequest struct {
		Token string
	}

	ViewBalanceResponse struct {
		ID        string    `json:"id"`
		OwnedBy   string    `json:"owned_by"`
		Status    string    `json:"status"`
		EnabledAt time.Time `json:"enabled_at"`
		Balance   float64   `json:"balance"`
	}
)
