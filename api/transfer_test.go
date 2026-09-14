import (
	"testing"

	"simple_bank/api"
)

func TestCreateTransfer(t *testing.T) {
	fromAccount := api.RandomAccount()
	toAccount := api.RandomAccount()

	// Here you would typically call the function to create a transfer
	// and then check the result using assertions.
	transfer, err := api.CreateTransfer(fromAccount.ID, toAccount.ID, 100)
	if err != nil {
		t.Fatalf("failed to create transfer: %v", err)
	}
	if transfer.FromAccountID != fromAccount.ID || transfer.ToAccountID != toAccount.ID || transfer.Amount != 100 {
		t.Errorf("unexpected transfer result: %+v", transfer)
	}
}