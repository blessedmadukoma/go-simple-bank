package db

import (
	"context"
	"testing"
	"time"

	"github.com/blessedmadukoma/go-simple-bank/util"
	"github.com/stretchr/testify/require"
)

// createRandomUser randomly creates a user
func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:    util.RandomOwner(),
		HashedPassword:  "",
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

// TestCreateAccount performs test operation to create a new account using random values
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// TestGetAccount tests if the GetAccount function retrieves an existing account
func TestGetUserByUsername(t *testing.T) {
	// create account
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
