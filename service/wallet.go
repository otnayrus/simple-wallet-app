package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	types "github.com/otnayrus/simple-wallet-app/types/wallet"
)

type walletService struct {
	walletRepo types.WalletRepository
}

func NewWalletService(wr types.WalletRepository) types.WalletService {
	return &walletService{
		walletRepo: wr,
	}
}

func (ws *walletService) Initialize(req types.InitializeRequest) (types.InitializeResponse, error) {
	newToken, err := makeToken()
	if err != nil {
		return types.InitializeResponse{}, err
	}

	err = ws.walletRepo.Create(types.Wallet{
		ID:      uuid.NewString(),
		OwnedBy: req.CustomerID,
		Token:   newToken,
		Status:  int(types.StatusNewlyCreated),
		Balance: 0,
	})
	if err != nil {
		return types.InitializeResponse{}, err
	}

	return types.InitializeResponse{
		Token: newToken,
	}, nil
}

func (ws *walletService) Enable(req types.EnableRequest) (types.EnableResponse, error) {
	wallet, err := ws.walletRepo.Enable(req.Token)
	if err != nil {
		log.Println("walletService.Enable.Enable", err)
		return types.EnableResponse{}, err
	}

	return types.EnableResponse{
		ID:        wallet.ID,
		OwnedBy:   wallet.OwnedBy,
		Status:    wallet.GetStatusString(),
		EnabledAt: wallet.UpdatedAt.Time,
		Balance:   wallet.Balance,
	}, nil
}

func (ws *walletService) ViewBalance(req types.ViewBalanceRequest) (types.ViewBalanceResponse, error) {
	wallet, err := ws.walletRepo.GetByToken(req.Token)
	if err != nil {
		log.Println("walletService.ViewBalance", err)
		return types.ViewBalanceResponse{}, err
	}

	if wallet.Status != int(types.StatusActive) {
		return types.ViewBalanceResponse{}, errors.New("wallet is disabled")
	}

	return types.ViewBalanceResponse{
		ID:        wallet.ID,
		OwnedBy:   wallet.OwnedBy,
		Status:    wallet.GetStatusString(),
		EnabledAt: wallet.UpdatedAt.Time,
		Balance:   wallet.Balance,
	}, nil
}

func (ws *walletService) Disable(req types.DisableRequest) (types.DisableResponse, error) {
	wallet, err := ws.walletRepo.Disable(req.Token)
	if err != nil {
		log.Println("walletService.Disable.Disable", err)
		return types.DisableResponse{}, err
	}

	return types.DisableResponse{
		ID:         wallet.ID,
		OwnedBy:    wallet.OwnedBy,
		Status:     wallet.GetStatusString(),
		DisabledAt: wallet.UpdatedAt.Time,
		Balance:    wallet.Balance,
	}, nil
}

func (ws *walletService) Deposit(req types.DepositRequest) (types.DepositResponse, error) {
	wallet, err := ws.walletRepo.GetByToken(req.Token)
	if err != nil {
		log.Println("walletService.Deposit", err)
		return types.DepositResponse{}, err
	}

	if wallet.Status != int(types.StatusActive) {
		return types.DepositResponse{}, errors.New("wallet is inactive")
	}

	mutation := types.Mutation{
		ID:          uuid.NewString(),
		ReferenceID: req.ReferenceID,
		CreatedAt:   time.Now(),
		CreatedBy:   wallet.OwnedBy,
		Action:      int(types.MutationActionDeposit),
		Status:      int(types.MutationStatusSuccess),
		Amount:      req.Amount,
	}
	err = ws.walletRepo.Mutate(mutation, wallet.Balance+req.Amount, req.Token)
	if err != nil {
		log.Println("walletService.Deposit", err)
		return types.DepositResponse{}, err
	}

	return types.DepositResponse{
		ID:          mutation.ID,
		DepositedBy: mutation.CreatedBy,
		Status:      mutation.GetStatusString(),
		DepositedAt: mutation.CreatedAt,
		Amount:      mutation.Amount,
		ReferenceID: mutation.ReferenceID,
	}, nil
}

func (ws *walletService) Withdraw(req types.WithdrawRequest) (types.WithdrawResponse, error) {
	wallet, err := ws.walletRepo.GetByToken(req.Token)
	if err != nil {
		log.Println("walletService.Withdraw", err)
		return types.WithdrawResponse{}, err
	}

	if wallet.Status != int(types.StatusActive) {
		return types.WithdrawResponse{}, errors.New("wallet is inactive")
	}

	if wallet.Balance < req.Amount {
		return types.WithdrawResponse{}, errors.New("insufficient funds")
	}

	mutation := types.Mutation{
		ID:          uuid.NewString(),
		ReferenceID: req.ReferenceID,
		CreatedAt:   time.Now(),
		CreatedBy:   wallet.OwnedBy,
		Action:      int(types.MutationActionWithdraw),
		Status:      int(types.MutationStatusSuccess),
		Amount:      req.Amount,
	}
	err = ws.walletRepo.Mutate(mutation, wallet.Balance-req.Amount, req.Token)
	if err != nil {
		log.Println("walletService.Withdraw", err)
		return types.WithdrawResponse{}, err
	}

	return types.WithdrawResponse{
		ID:          mutation.ID,
		WithdrawnBy: mutation.CreatedBy,
		Status:      mutation.GetStatusString(),
		WithdrawnAt: mutation.CreatedAt,
		Amount:      mutation.Amount,
		ReferenceID: mutation.ReferenceID,
	}, nil
}

// helpers

func makeToken() (string, error) {
	length := 20

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println("Error generating random string:", err)
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}
