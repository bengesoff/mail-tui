package core

import "time"

type Email struct {
	From    string
	To      string
	Subject string
	SentAt  time.Time
	Body    string
}
