package db

import (
	"context"
	"testing"

	"simple_bank/util"

	"github.com/stretchr/testify/assert"
)

func createRandomTransfer(t *testing.T) Transfer {
	ctx := context.Background()

	arg := CreateTransferParams{
		FromAccountID: createRandomAccount(t).ID,
		ToAccountID:   createRandomAccount(t).ID,
		Amount:        util.RandomBalance(),
	}

	transfer, err := testQueries.CreateTransfer(ctx, arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, transfer)

	assert.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	assert.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	assert.Equal(t, transfer.Amount, arg.Amount)

	assert.NotZero(t, transfer.ID)
	assert.NotZero(t, transfer.CreatedAt)

	return transfer
}

func createRandomTransfers(t *testing.T, n int) []Transfer {
	transfers := make([]Transfer, n)
	for i := range transfers {
		transfers[i] = createRandomTransfer(t)
	}
	return transfers
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	ctx := context.Background()

	transfer := createRandomTransfer(t)

	gotTransfer, err := testQueries.GetTransfer(ctx, transfer.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, gotTransfer)

	assert.Equal(t, gotTransfer.ID, transfer.ID)
	assert.Equal(t, gotTransfer.FromAccountID, transfer.FromAccountID)
	assert.Equal(t, gotTransfer.ToAccountID, transfer.ToAccountID)
	assert.Equal(t, gotTransfer.Amount, transfer.Amount)

	assert.WithinDuration(t, gotTransfer.CreatedAt, transfer.CreatedAt, 0)
}

func TestListTransfers(t *testing.T) {
	ctx := context.Background()

	createRandomTransfers(t, 5)

	args := ListTransfersParams{
		Limit:  3,
		Offset: 1,
	}
	transfers, err := testQueries.ListTransfers(ctx, args)

	assert.NoError(t, err)
	assert.NotEmpty(t, transfers)
	assert.Len(t, transfers, 3)
	for _, transfer := range transfers {
		assert.NotEmpty(t, transfer)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)
	}
}
