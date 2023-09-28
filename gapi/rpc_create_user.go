package gapi

import (
	"context"
	"time"

	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/blessedmadukoma/go-simple-bank/gvalidator"
	"github.com/blessedmadukoma/go-simple-bank/pb"
	"github.com/blessedmadukoma/go-simple-bank/util"
	"github.com/blessedmadukoma/go-simple-bank/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := srv.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Unimplemented, "failed to create user: %s", err)
	}

	// TODO: use db transaction

	// send verify email to user
	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}

	// set task to be processed or retried
	opts := []asynq.Option{
		asynq.MaxRetry(10),                // to be retried at most 10 times
		asynq.ProcessIn(10 * time.Second), // to be processed in 10 seconds
		asynq.Queue(worker.QueueCritical), // to be processed in the "critical" queue
	}

	err = srv.taskDistributor.DistributeTaskSendVerifyEmail(ctx, *taskPayload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "failed to distribute task to send verify email: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

// validateCreateUserRequest validates the request body of CreateUserRequest.
func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := gvalidator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := gvalidator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := gvalidator.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := gvalidator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
