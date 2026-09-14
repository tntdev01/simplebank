package api

import (
	"database/sql"
	"fmt"

	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// createAccountRequest represents the expected payload for creating a new account.
type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// createTransfer handles the HTTP request for creating a new transfer.
func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, result)
}

// validAccount checks if the account with the given ID exists and has the specified currency.
func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, errorResponse(err))
			return false
		}
		ctx.JSON(500, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account currency mismatch: %s vs %s", account.Currency, currency)
		ctx.JSON(400, errorResponse(err))
		return false
	}

	return true
}
