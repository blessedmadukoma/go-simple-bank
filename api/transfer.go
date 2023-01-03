package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"` // gt - greater than
	Currency      string `json:"currency" binding:"required,currency"`
}

// createTransfer creates a transfer between two accounts
func (srv *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("User input not valid", err))
		return
	}

	if !srv.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
 
	if !srv.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := srv.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Unable to perform transfer", err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
	return
}

// validAccount checks if an account exists and the currency matches the input currency
func (srv *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := srv.store.GetAccountByID(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse("No record of this user", err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving user account", err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse("Currency mismatch!", err))
		return false
	}

	return true
}
