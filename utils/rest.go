package utils

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

const (
	StatusSuccess = "success"
	StatusError   = "fail"
)

type RestResponse struct {
	Status string
	Data   interface{}
}

func MakeRestResponse(w http.ResponseWriter, data interface{}, status int, err error) {
	var statusStr string = StatusSuccess
	if err != nil {
		statusStr = StatusError
		data = struct {
			Error error
		}{
			Error: err,
		}
	}

	res := RestResponse{
		Status: statusStr,
		Data:   data,
	}
	marshal, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to marshal rest response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(marshal)
}

func InitRestContext() context.Context {
	return context.Background()
}

type WalletWrapper struct {
	Wallet interface{} `json:"wallet"`
}

func AddWalletWrapper(data interface{}) WalletWrapper {
	return WalletWrapper{
		Wallet: data,
	}
}
