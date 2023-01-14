package model

import (
	"context"
	"lecture/go-contracts/conf"
)

type RequestTransferFrom struct {
	Value      int64  `json:"value"`
	PrivateKey string `json:"privatekey"`
}

type Model struct {
	ctx context.Context
	cfg *conf.Config
}
