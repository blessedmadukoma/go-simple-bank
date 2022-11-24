package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransferByID(context.Background(), transfer.ID)
		require.NoError(t, err)

		sourceAccountTransaction := result.FromEntry
		require.NotEmpty(t, sourceAccountTransaction)
		require.Equal(t, account1.ID, sourceAccountTransaction.AccountID)
		require.Equal(t, amount, sourceAccountTransaction.Amount)
		require.NotZero(t, sourceAccountTransaction.ID)

		_, err = store.GetTransferByID(context.Background(), sourceAccountTransaction.ID)
		require.NoError(t, err)

		destinationAccountTransaction := result.ToEntry
		require.NotEmpty(t, destinationAccountTransaction)
		require.Equal(t, account2.ID, destinationAccountTransaction.AccountID)
		require.Equal(t, amount, destinationAccountTransaction.Amount)
		require.NotZero(t, destinationAccountTransaction.ID)

		_, err = store.GetEntryByID(context.Background(), destinationAccountTransaction.ID)
		require.NoError(t, err)

		//go:TODO check account balance
	}

}
