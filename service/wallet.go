package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

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
