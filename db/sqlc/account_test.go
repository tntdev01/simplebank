package db

import (
	"context"
	"database/sql"
	"testing"

	"simple_bank/util"

	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T) Account {
	ctx := context.Background()
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(ctx, arg)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, account.Owner, arg.Owner)
	assert.Equal(t, account.Balance, arg.Balance)
	assert.Equal(t, account.Currency, arg.Currency)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func createRandomAccounts(t *testing.T, n int) []Account {
	accounts := make([]Account, n)
	for i := range accounts {
		accounts[i] = createRandomAccount(t)
	}
	return accounts
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	ctx := context.Background()

	account := createRandomAccount(t)

	gotAccount, err := testQueries.GetAccount(ctx, account.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, gotAccount)

	assert.Equal(t, gotAccount.ID, account.ID)
	assert.Equal(t, gotAccount.Owner, account.Owner)
	assert.Equal(t, gotAccount.Balance, account.Balance)
	assert.Equal(t, gotAccount.Currency, account.Currency)

	assert.WithinDuration(t, gotAccount.CreatedAt, account.CreatedAt, 0)
}

func TestListAccounts(t *testing.T) {
	ctx := context.Background()

	createRandomAccounts(t, 5)

	args := ListAccountsParams{
		Limit:  3,
		Offset: 1,
	}
	accounts, err := testQueries.ListAccounts(ctx, args)

	assert.NoError(t, err)
	assert.NotEmpty(t, accounts)
	assert.Len(t, accounts, 3)
	for _, account := range accounts {
		assert.NotEmpty(t, account)
		assert.NotZero(t, account.ID)
		assert.NotZero(t, account.CreatedAt)
	}
}

func TestUpdateAccount(t *testing.T) {
	ctx := context.Background()

	account := createRandomAccount(t)

	updateArg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomBalance(),
	}

	updatedAccount, err := testQueries.UpdateAccount(ctx, updateArg)
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedAccount)

	assert.Equal(t, updatedAccount.ID, account.ID)
	assert.Equal(t, updatedAccount.Owner, account.Owner)
	assert.Equal(t, updatedAccount.Balance, updateArg.Balance)
	assert.Equal(t, updatedAccount.Currency, account.Currency)

	assert.WithinDuration(t, updatedAccount.CreatedAt, account.CreatedAt, 0)
}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()

	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(ctx, account.ID)
	assert.NoError(t, err)

	_, err = testQueries.GetAccount(ctx, account.ID)
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}
