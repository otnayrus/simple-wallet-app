package rest

import (
	"errors"
	"log"
	"net/http"
	"strings"

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

func (wh *walletHandler) Enable(c *gin.Context) {
	var header types.WalletRequestHeader

	if err := c.BindHeader(&header); err != nil {
		log.Println("fail to decode request header")
		return
	}

	split := strings.Split(header.Authorization, " ")
	if len(split) != 2 && split[0] != "Token" {
		utils.MakeRestResponse(c.Writer, nil, http.StatusBadRequest, errors.New("invalid request header"))
		return
	}
	token := split[1]

	res, err := wh.walletService.Enable(types.EnableRequest{Token: token})
	if err != nil {
		utils.MakeRestResponse(c.Writer, nil, http.StatusInternalServerError, err)
		return
	}

	utils.MakeRestResponse(c.Writer, utils.AddWalletWrapper(res), http.StatusCreated, nil)
}

func (wh *walletHandler) ViewBalance(c *gin.Context) {
	var header types.WalletRequestHeader

	if err := c.BindHeader(&header); err != nil {
		log.Println("fail to decode request header")
		return
	}

	split := strings.Split(header.Authorization, " ")
	if len(split) != 2 && split[0] != "Token" {
		utils.MakeRestResponse(c.Writer, nil, http.StatusBadRequest, errors.New("invalid request header"))
		return
	}
	token := split[1]

	res, err := wh.walletService.ViewBalance(types.ViewBalanceRequest{Token: token})
	if err != nil {
		utils.MakeRestResponse(c.Writer, nil, http.StatusInternalServerError, err)
		return
	}

	utils.MakeRestResponse(c.Writer, utils.AddWalletWrapper(res), http.StatusCreated, nil)
}

func (wh *walletHandler) Disable(c *gin.Context) {
	var header types.WalletRequestHeader

	if err := c.BindHeader(&header); err != nil {
		log.Println("fail to decode request header")
		return
	}

	split := strings.Split(header.Authorization, " ")
	if len(split) != 2 && split[0] != "Token" {
		utils.MakeRestResponse(c.Writer, nil, http.StatusBadRequest, errors.New("invalid request header"))
		return
	}
	token := split[1]

	var req types.DisableRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println("fail to decode request body")
		return
	}
	if !req.IsDisabled {
		utils.MakeRestResponse(c.Writer, nil, http.StatusBadRequest, errors.New("is disabled flag is false"))
		return
	}

	res, err := wh.walletService.Disable(types.DisableRequest{Token: token})
	if err != nil {
		utils.MakeRestResponse(c.Writer, nil, http.StatusInternalServerError, err)
		return
	}

	utils.MakeRestResponse(c.Writer, utils.AddWalletWrapper(res), http.StatusCreated, nil)
}
