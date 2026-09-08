package db

import (
	"context"
	"testing"

	"simple_bank/util"

	"github.com/stretchr/testify/assert"
)

func createRandomEntry(t *testing.T) Entry {
	ctx := context.Background()
	arg := CreateEntryParams{
		AccountID: createRandomAccount(t).ID,
		Amount:    util.RandomBalance(),
	}

	entry, err := testQueries.CreateEntry(ctx, arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, entry)

	assert.Equal(t, entry.AccountID, arg.AccountID)
	assert.Equal(t, entry.Amount, arg.Amount)

	assert.NotZero(t, entry.ID)
	assert.NotZero(t, entry.CreatedAt)

	return entry
}

func createRandomEntries(t *testing.T, n int) []Entry {
	entries := make([]Entry, n)
	for i := range entries {
		entries[i] = createRandomEntry(t)
	}
	return entries
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	ctx := context.Background()

	entry := createRandomEntry(t)

	gotEntry, err := testQueries.GetEntry(ctx, entry.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, gotEntry)

	assert.Equal(t, gotEntry.ID, entry.ID)
	assert.Equal(t, gotEntry.AccountID, entry.AccountID)
	assert.Equal(t, gotEntry.Amount, entry.Amount)

	assert.WithinDuration(t, gotEntry.CreatedAt, entry.CreatedAt, 0)
}

func TestListEntries(t *testing.T) {
	ctx := context.Background()

	createRandomEntries(t, 5)

	args := ListEntriesParams{
		Limit:  3,
		Offset: 1,
	}
	entries, err := testQueries.ListEntries(ctx, args)

	assert.NoError(t, err)
	assert.NotEmpty(t, entries)
	assert.Len(t, entries, 3)
	for _, entry := range entries {
		assert.NotEmpty(t, entry)
		assert.NotZero(t, entry.ID)
		assert.NotZero(t, entry.CreatedAt)
	}
}
