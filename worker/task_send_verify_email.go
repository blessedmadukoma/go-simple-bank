package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/blessedmadukoma/go-simple-bank/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// PayloadSendVerifyEmail is the payload for sending a verification email.
type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

const TaskSendVeryEmail = "task:send_verify_email"

// DistributeTaskSendVerifyEmail distributes a task, and send a verification email.
func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(

	ctx context.Context,
	payload PayloadSendVerifyEmail,
	opts ...asynq.Option,

) error {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVeryEmail, jsonPayload)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retries", info.MaxRetry).Msg("enqueued task")

	return nil
}

// ProcessTaskSendVerifyEmail processes a task, and send a verification email.
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	// send email to user
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	verifyURL := fmt.Sprintf("http://localhost:3000/verify-email?id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s, <br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyURL)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("processed task")

	return nil
}
