package imap

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/bengesoff/mail-tui/internal/core"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type ImapBackend struct {
	client *imapclient.Client
}

func NewImapBackend(address, username, password string) (*ImapBackend, error) {
	// should really use TLS
	client, err := imapclient.DialInsecure(address, nil)
	if err != nil {
		return nil, err
	}
	err = client.Login(username, password).Wait()
	if err != nil {
		return nil, err
	}

	_, err = client.Select("INBOX", nil).Wait()
	if err != nil {
		return nil, err
	}

	return &ImapBackend{client}, nil
}

// ListEmails fetches all messages.
// It does not do any pagination, but it should do for large mailboxes.
func (b *ImapBackend) ListEmails() ([]core.EmailMetadata, error) {
	sequenceSet := imap.SeqSet{}
	// fetch 1:* (for fetching all)
	sequenceSet.AddRange(1, 0)
	messages, err := b.client.Fetch(sequenceSet, &imap.FetchOptions{
		Envelope: true,
	}).Collect()
	if err != nil {
		return nil, err
	}

	result := []core.EmailMetadata{}
	for _, message := range messages {
		result = append(result, fetchMessageBufferToEmailMetadata(message))
	}
	return result, nil
}

// GetEmail fetches a single email by its sequence number.
func (b *ImapBackend) GetEmail(id core.EmailId) (*core.Email, error) {
	sequenceNumber, err := strconv.Atoi(string(id))
	if err != nil {
		return nil, err
	}

	messages, err := b.client.Fetch(imap.SeqSetNum(uint32(sequenceNumber)), &imap.FetchOptions{
		Envelope:    true,
		BodySection: []*imap.FetchItemBodySection{{Specifier: imap.PartSpecifierText}},
	}).Collect()
	if err != nil {
		return nil, err
	}
	if len(messages) != 1 {
		return nil, fmt.Errorf("expected 1 message, got %d", len(messages))
	}

	message := messages[0]
	// fairly rudimentary body decoding - will probably break for more complicated multipart messages
	body := string(message.BodySection[0].Bytes)

	return &core.Email{
		EmailMetadata: fetchMessageBufferToEmailMetadata(message),
		Body:          body,
	}, nil
}

// SendEmail does nothing currently.
func (b *ImapBackend) SendEmail(email core.OutgoingEmail) error {
	time.Sleep(1 * time.Second)
	// should use net/smtp to send the message somewhere
	return nil
}

// MarkAsRead uses the STORE command to add the SEEN flag to an email with a given sequence number.
func (b *ImapBackend) MarkAsRead(id core.EmailId) error {
	sequenceNumber, err := strconv.Atoi(string(id))
	if err != nil {
		return err
	}

	return b.client.Store(
		imap.SeqSetNum(uint32(sequenceNumber)),
		&imap.StoreFlags{
			Op:     imap.StoreFlagsAdd,
			Flags:  []imap.Flag{imap.FlagSeen},
			Silent: true,
		},
		nil).Close()
}

func (b *ImapBackend) Close() error {
	err := b.client.Logout().Wait()
	if err != nil {
		return err
	}
	return b.client.Close()
}

func fetchMessageBufferToEmailMetadata(message *imapclient.FetchMessageBuffer) core.EmailMetadata {
	return core.EmailMetadata{
		Id:      core.EmailId(strconv.Itoa(int(message.SeqNum))),
		Subject: message.Envelope.Subject,
		From:    message.Envelope.From[0].Addr(),
		To:      message.Envelope.To[0].Addr(),
		SentAt:  message.Envelope.Date,
		IsRead:  slices.Contains(message.Flags, imap.FlagSeen),
	}
}
