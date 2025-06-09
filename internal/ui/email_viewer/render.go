package email_viewer

import (
	"fmt"

	"mail-tui/internal/core"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

func RenderEmail(email *core.Email, windowWidth int) (string, error) {
	var output string
	output += fmt.Sprintf("From: %s\n", email.From)
	output += fmt.Sprintf("To: %s\n", email.To)

	sent, err := email.SentAt.MarshalText()
	if err != nil {
		return "", err
	}
	output += fmt.Sprintf("Sent: %s\n", sent)

	output += fmt.Sprintf("Subject: %s\n\n", email.Subject)

	output += indent.String(
		wrap.String(
			wordwrap.String(
				email.Body,
				windowWidth-2),
			windowWidth-2),
		2)

	return output, nil
}
