package mail

import (
	"testing"

	"github.com/blessedmadukoma/go-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Hello from Go!"
	content := `
	<h1>This is a test email from Go!</h1>
	<p>Feel free to delete this email.</p>
	`

	to := []string{"bmadukoma@gmail.com"}
	files := []string{"../go.mod"}

	err = sender.SendEmail(subject, content, to, nil, nil, files)
	require.NoError(t, err)

}
