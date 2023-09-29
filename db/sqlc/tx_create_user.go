package db

import "context"

// CreateUserTxParams contains the input parameters of the create user transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserTxResult is the result of the create user transaction
type CreateUserTxResult struct {
	User User
}

// CreateUserTx performs a
func (s *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User) // call the callback function to allow for retry if there is an error
	})

	return result, err
}
