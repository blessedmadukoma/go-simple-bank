package worker

import (
	"context"

	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

// NewRedisTaskProcessor creates a new RedisTaskProcessor
func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

// Start starts the RedisTaskProcessor
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVeryEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
