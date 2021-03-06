package web

import (
	"errors"
	"fmt"
	"net/http"

	"chainlink/core/services"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/store/presenters"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// TransfersController can send LINK tokens to another address
type TransfersController struct {
	App services.Application
}

// Create sends ETH from the Chainlink's account to a specified address.
//
// Example: "<application>/withdrawals"
func (tc *TransfersController) Create(c *gin.Context) {
	var tr models.SendEtherRequest

	store := tc.App.GetStore()

	if err := c.ShouldBindJSON(&tr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
	} else if from, err := retrieveFromAddress(tr.FromAddress, store); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
	} else if tx, err := store.TxManager.CreateTxWithEth(from, tr.DestinationAddress, tr.Amount); err != nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("Transaction failed: %v", err))
	} else {
		jsonAPIResponse(c, presenters.NewTx(tx), "transaction")
	}
}

func retrieveFromAddress(from common.Address, store *store.Store) (common.Address, error) {
	if from != utils.ZeroAddress {
		return from, nil
	}
	ma := store.TxManager.NextActiveAccount()
	if ma == nil {
		return common.Address{}, errors.New("Must activate an account before creating a transaction")
	}

	return ma.Address, nil
}
