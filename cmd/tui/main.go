package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bengesoff/mail-tui/internal/backend/fake"
	"github.com/bengesoff/mail-tui/internal/backend/imap"
	"github.com/bengesoff/mail-tui/internal/core"
	"github.com/bengesoff/mail-tui/internal/ui/app"
)

var flags struct {
	useImap     bool
	imapAddress string
	username    string
	password    string
}

func main() {
	flag.BoolVar(&flags.useImap, "use-imap", true, "Use IMAP backend")
	flag.StringVar(&flags.imapAddress, "imap-address", "localhost:1143", "IMAP server address (hostname:port)")
	// Not good, shouldn't be passed in plaintext. Ideally would be an environment variable
	flag.StringVar(&flags.username, "username", "bob", "IMAP username")
	flag.StringVar(&flags.password, "password", "pass", "IMAP password")

	flag.Parse()

	var backend core.EmailBackend
	if flags.useImap {
		// could also be initialised inside the bubbletea program in order to display a loading spinner
		imapBackend, err := imap.NewImapBackend(flags.imapAddress, flags.username, flags.password)
		if err != nil {
			fmt.Printf("failed to create IMAP backend: %v\n", err)
			os.Exit(1)
		}
		backend = imapBackend
		defer func() { _ = imapBackend.Close() }()
	} else {
		backend = fake.NewFakeBackend()
	}

	appModel := app.NewAppModel(backend)
	program := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion())

	_, err := program.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
