package core

type EmailBackend interface {
	ListEmails() ([]EmailMetadata, error)
	GetEmail(id EmailId) (*Email, error)
	SendEmail(email OutgoingEmail) error
}
