package ui

import "github.com/bengesoff/mail-tui/internal/core"

type ShowEmailListMessage struct{}

type ShowEmailViewerMessage struct {
	EmailId core.EmailId
}
