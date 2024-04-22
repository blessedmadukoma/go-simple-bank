package api

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/blessedmadukoma/go-simple-bank/token"
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

	fromAccount, valid := srv.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse("could not transfer money", err))
		// ctx.JSON(http.StatusUnauthorized, errorResponse("from account does not belong to the authenticated user", nil))
		return
	}

	_, valid = srv.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	senderAccount, _ := srv.store.GetAccountByID(ctx, req.FromAccountID)
	if req.Amount > senderAccount.Balance {
		err := fmt.Errorf("Amount can not be greater than balance")
		ctx.JSON(http.StatusBadRequest, errorResponse("cannot perform transfer", err))
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
func (srv *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := srv.store.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse("No record of this user", err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving user account", err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse("Currency mismatch!", err))
		return account, false
	}

	return account, true
}
