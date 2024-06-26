package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf(">>before -> Acct1 bal: %d, acct2 bal: %d\n", account1.Balance, account2.Balance)
	//run n concurrent transfer transaction
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)
	amount := int64(10)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		transfer := result.Transfer
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		_, err = testStore.GetTransferByID(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)
		require.NotZero(t, fromEntry.ID)

		_, err = testStore.GetEntryByID(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)

		_, err = testStore.GetEntryByID(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		fmt.Println(fromAccount.Balance, toAccount.Balance, amount)

		//check balance on account
		fmt.Printf(">>transaction tx -> fromAccct bal: %d, toAcct bal: %d\n", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final updated balance
	updatedAccount1, err := testStore.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testStore.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">>after tansaction -> updatedAcct1 bal: %d, updatedAcct2 bal: %d \n", updatedAccount1.Balance, updatedAccount1.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

// func TestTransferTxDeadlock(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	account2 := createRandomAccount(t)

// 	fmt.Println(">> before:", account1.Balance, account2.Balance)

// 	n := 10
// 	amount := int64(10)
// 	errs := make(chan error)

// 	for i := 0; i < n; i++ {
// 		FromAccountID := account1.ID
// 		ToAccountID := account2.ID

// 		if i%2 == 1 {
// 			FromAccountID = account2.ID
// 			ToAccountID = account1.ID
// 		}
// 	}
// }
