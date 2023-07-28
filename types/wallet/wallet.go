package types

import (
	"database/sql"
	"time"
)

type WalletService interface {
	Initialize(InitializeRequest) (InitializeResponse, error)
	Enable(EnableRequest) (EnableResponse, error)
	ViewBalance(ViewBalanceRequest) (ViewBalanceResponse, error)
	Disable(DisableRequest) (DisableResponse, error)
	Deposit(DepositRequest) (DepositResponse, error)
	Withdraw(WithdrawRequest) (WithdrawResponse, error)
}

type WalletRepository interface {
	Create(Wallet) error
	Enable(string) (Wallet, error)
	GetByToken(string) (Wallet, error)
	Disable(token string) (Wallet, error)
	Mutate(Mutation, float64, string) error
	SetWalletBalanceByToken(float64, string) error
	CreateMutation(Mutation) error
}

type WalletStatus int

const (
	StatusInactive WalletStatus = iota
	StatusNewlyCreated
	StatusActive
)

const (
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

	DisableRequest struct {
		Token      string
		IsDisabled bool `form:"is_disabled"`
	}

	DisableResponse struct {
		ID         string    `json:"id"`
		OwnedBy    string    `json:"owned_by"`
		Status     string    `json:"status"`
		DisabledAt time.Time `json:"disabled_at"`
		Balance    float64   `json:"balance"`
	}

	DepositRequest struct {
		Token       string
		ReferenceID string  `form:"reference_id"`
		Amount      float64 `form:"amount"`
	}

	DepositResponse struct {
		ID          string    `json:"id"`
		DepositedBy string    `json:"deposited_by"`
		Status      string    `json:"status"`
		DepositedAt time.Time `json:"deposited_at"`
		Amount      float64   `json:"amount"`
		ReferenceID string    `json:"reference_id"`
	}

	WithdrawRequest struct {
		Token       string
		ReferenceID string  `form:"reference_id"`
		Amount      float64 `form:"amount"`
	}

	WithdrawResponse struct {
		ID          string    `json:"id"`
		WithdrawnBy string    `json:"withdrawn_by"`
		Status      string    `json:"status"`
		WithdrawnAt time.Time `json:"withdrawn_at"`
		Amount      float64   `json:"amount"`
		ReferenceID string    `json:"reference_id"`
	}
)
