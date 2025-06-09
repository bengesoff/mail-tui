package email_viewer

import (
	"mail-tui/internal/core"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

var (
	metadataHeadingStyle = lipgloss.NewStyle().
				Bold(true)

	bodyStyle = func(width int) lipgloss.Style {
		return lipgloss.NewStyle().
			Width(width).
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1)
	}
)

func RenderEmail(email *core.Email, windowWidth int) (string, error) {
	sent, err := email.SentAt.MarshalText()
	if err != nil {
		return "", err
	}

	rows := [][]string{
		{"From", email.From},
		{"To", email.To},
		{"Sent", string(sent)},
		{"Subject", email.Subject},
	}
	metadata := table.New().
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return metadataHeadingStyle
			}
			return lipgloss.NewStyle()
		}).
		Rows(rows...)
	output := metadata.Render() + "\n"

	output += bodyStyle(windowWidth - 2).Render(
		wrap.String(
			wordwrap.String(
				email.Body,
				windowWidth-4),
			windowWidth-4),
	)

	return output, nil
}
