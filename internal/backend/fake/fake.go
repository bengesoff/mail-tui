package fake

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/bengesoff/mail-tui/internal/core"
)

type FakeBackend struct {
	emails map[core.EmailId]core.EmailMetadata
}

func NewFakeBackend() *FakeBackend {
	return &FakeBackend{
		emails: map[core.EmailId]core.EmailMetadata{
			"1": {
				Id:      "1",
				From:    "test1@example.com",
				To:      "me@example.com",
				Subject: "First email",
				SentAt:  time.Now(),
				IsRead:  false,
			},
			"2": {
				Id:      "2",
				From:    "test2@example.com",
				To:      "me@example.com",
				Subject: "Second email",
				SentAt:  time.Now().Add(-1 * time.Hour),
				IsRead:  false,
			},
			"3": {
				Id:      "3",
				From:    "test3@example.com",
				To:      "me@example.com",
				Subject: "Third email",
				SentAt:  time.Now().Add(-2 * time.Hour),
				IsRead:  false,
			},
			"4": {
				Id:      "4",
				From:    "test4@example.com",
				To:      "me@example.com",
				Subject: "Fourth email",
				SentAt:  time.Now().Add(-3 * time.Hour),
				IsRead:  false,
			},
		},
	}
}

func (b *FakeBackend) ListEmails() ([]core.EmailMetadata, error) {
	time.Sleep(1 * time.Second)
	// just a fake implementation, otherwise should probably use an ordered data structure
	return slices.SortedFunc(maps.Values(b.emails), func(a, b core.EmailMetadata) int {
		if a.SentAt.Before(b.SentAt) {
			return 1
		}
		if a.SentAt.After(b.SentAt) {
			return -1
		}
		return 0
	}), nil
}

func (b *FakeBackend) GetEmail(id core.EmailId) (*core.Email, error) {
	time.Sleep(1 * time.Second)
	email, ok := b.emails[id]
	if !ok {
		return nil, fmt.Errorf("email not found")
	}
	return &core.Email{
		EmailMetadata: email,
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

func (b *FakeBackend) MarkAsRead(id core.EmailId) error {
	time.Sleep(1 * time.Second)
	email, ok := b.emails[id]
	if !ok {
		return fmt.Errorf("email not found")
	}
	email.IsRead = true
	b.emails[id] = email
	return nil
}
