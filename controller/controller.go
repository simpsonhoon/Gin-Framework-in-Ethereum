package controller

import (
	conf "lecture/go-contracts/conf"
	"lecture/go-contracts/model"
	service "lecture/go-contracts/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
}

func NewCTL(rep *model.Model) (*Controller, error) {
	r := &Controller{md: rep}
	return r, nil
}

func (p *Controller) GetTokenSymbol(c *gin.Context) {
	symbol, err := service.GetSymbol(conf.Config.Contract.contractaddress, c.Param("tokenName"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"symbol": symbol})
}

func (p *Controller) GetTokenBalanceByAddress(c *gin.Context) {
	balance, err := service.GetTokenBalaceByAddress(conf.Config.Contract.contractaddress, c.Param("address"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"balance": balance})
}

func (p *Controller) CoinTransfer(c *gin.Context) {
	var account model.RequestTransferFrom
	receipt, err := service.TransferWemix(conf.Config.Address.privateKey, c.Param("address"), account.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"Transaction": receipt})
}

func (p *Controller) CoinTransferFrom(c *gin.Context) {
	var account model.RequestTransferFrom
	receipt, err := service.TransferWemix(account.PrivateKey, c.Param("address"), account.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"Transaction": receipt})
}

func (p *Controller) TokenTransfer(c *gin.Context) {
	var account model.RequestTransferFrom
	receipt, err := service.TransferCtxSSH(conf.Config.Address.privateKey, c.Param("address"), account.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"Transaction": receipt})
}

func (p *Controller) TokenTransferFrom(c *gin.Context) {
	var account model.RequestTransferFrom
	receipt, err := service.TransferCtxSSH(account.PrivateKey, c.Param("address"), account.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"Transaction": receipt})
}
