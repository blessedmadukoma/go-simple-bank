package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

// TaskDistributor
type TaskDistributor interface {
	// DistributeTask()
	// SendEmail()
	// VerifyEmail()
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

// NewRedisTaskDistributor creates a new RedisTaskDistributor
func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)

	return &RedisTaskDistributor{
		client: client,
	}
}
