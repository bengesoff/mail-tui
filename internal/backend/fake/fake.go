package fake

import (
	"fmt"
	"time"

	"github.com/bengesoff/mail-tui/internal/core"
)

type FakeBackend struct{}

func (b *FakeBackend) ListEmails() ([]core.EmailMetadata, error) {
	time.Sleep(1 * time.Second)
	return []core.EmailMetadata{
		{
			Id:      "1",
			From:    "test1@example.com",
			To:      "me@example.com",
			Subject: "First email",
			SentAt:  time.Now(),
		},
		{
			Id:      "2",
			From:    "test2@example.com",
			To:      "me@example.com",
			Subject: "Second email",
			SentAt:  time.Now().Add(-1 * time.Hour),
		},
		{
			Id:      "3",
			From:    "test3@example.com",
			To:      "me@example.com",
			Subject: "Third email",
			SentAt:  time.Now().Add(-2 * time.Hour),
		},
		{
			Id:      "4",
			From:    "test4@example.com",
			To:      "me@example.com",
			Subject: "Fourth email",
			SentAt:  time.Now().Add(-3 * time.Hour),
		},
	}, nil
}

func (b *FakeBackend) GetEmail(id core.EmailId) (*core.Email, error) {
	time.Sleep(1 * time.Second)
	return &core.Email{
		EmailMetadata: core.EmailMetadata{
			Id:      id,
			From:    "from@example.com",
			To:      "to@example.com",
			Subject: "Test email",
			SentAt:  time.Now(),
		},
		Body: fmt.Sprintf("To whom it may concern,\n\n"+
			"This is a test email with ID %s.\n\n"+
			"Yours sincerely,\n\n"+
			"Tester", id),
	}, nil
}

func (b *FakeBackend) SendEmail(email core.OutgoingEmail) error {
	time.Sleep(1 * time.Second)
	return nil
}
