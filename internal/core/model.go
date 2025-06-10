package core

import "time"

type EmailId string

type EmailMetadata struct {
	Id      EmailId
	From    string
	To      string
	Subject string
	SentAt  time.Time
}

type Email struct {
	EmailMetadata
	Body string
}

type OutgoingEmail struct {
	To      string
	Subject string
	Body    string
}
