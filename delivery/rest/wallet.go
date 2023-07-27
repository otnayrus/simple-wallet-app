package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	types "github.com/otnayrus/simple-wallet-app/types/wallet"
	"github.com/otnayrus/simple-wallet-app/utils"
)

type walletHandler struct {
	walletService types.WalletService
}

func NewWalletHandler(ws types.WalletService) walletHandler {
	return walletHandler{
		walletService: ws,
	}
}

func (wh *walletHandler) Initialize(c *gin.Context) {
	var req types.InitializeRequest

	if err := c.Bind(&req); err != nil {
		log.Println("fail to decode request body")
		return
	}

	res, err := wh.walletService.Initialize(req)
	if err != nil {
		utils.MakeRestResponse(c.Writer, nil, http.StatusInternalServerError, err)
		return
	}

	utils.MakeRestResponse(c.Writer, res, http.StatusCreated, nil)
}
