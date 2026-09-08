package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store instance.
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

var txKey = struct{}{}

// execTx executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters for the TransferTx function.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of the TransferTx function.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to another.
// It creates a transfer record, updates account balances, and ensures all operations are executed
// within a single database transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create Entry Records for both accounts involved in the transfer
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Get account details for both accounts involved in the transfer
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = transferMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = transferMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	amount int64) (Account, Account, error) {
	fromAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: -amount,
		ID:     fromAccountID,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	toAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount,
		ID:     toAccountID,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	// Return the updated fromAccount as a representative result
	return fromAccount, toAccount, nil
}
