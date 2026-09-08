package db

import (
	"context"
	"testing"

	"simple_bank/util"

	"github.com/stretchr/testify/assert"
)

func TestTransferTxDeadlock(t *testing.T) {
	ctx := context.Background()

	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	amount := util.RandomInt(1, 20)

	n := 10

	errs := make(chan error)

	for i := 0; i < n; i++ {
		var fromAccountID, toAccountID int64
		if i%2 == 1 {
			fromAccountID, toAccountID = toAccount.ID, fromAccount.ID
		} else {
			fromAccountID, toAccountID = fromAccount.ID, toAccount.ID
		}
		go func() {
			args := TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			}
			_, err := store.TransferTx(ctx, args)
			errs <- err
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)
	}
	// check accounts' balance after transfer
	fromAccountUpdated, err := store.GetAccount(ctx, fromAccount.ID)
	assert.NoError(t, err)
	toAccountUpdated, err := store.GetAccount(ctx, toAccount.ID)
	assert.NoError(t, err)

	assert.Equal(t, fromAccount.Balance, fromAccountUpdated.Balance)
	assert.Equal(t, toAccount.Balance, toAccountUpdated.Balance)
}
