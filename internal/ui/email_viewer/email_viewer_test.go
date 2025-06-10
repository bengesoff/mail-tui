package email_viewer

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bengesoff/mail-tui/internal/core"
	"github.com/bengesoff/mail-tui/internal/ui"
)

type mockBackend struct {
	email *core.Email
	err   error
}

func (m *mockBackend) ListEmails() ([]core.EmailMetadata, error) {
	return nil, nil
}

func (m *mockBackend) GetEmail(id core.EmailId) (*core.Email, error) {
	return m.email, m.err
}

func TestEmailViewerModel_ShowEmailViewerMessage(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)

	updatedModel, cmd := model.Update(ui.ShowEmailViewerMessage{EmailId: "test-id"})

	if !updatedModel.loading {
		t.Error("Expected loading to be true after ShowEmailViewerMessage")
	}

	if updatedModel.error != "" {
		t.Error("Expected error to be cleared")
	}

	if updatedModel.email != nil {
		t.Error("Expected email to be nil when starting to load")
	}

	// Should return a command to load the email
	if cmd == nil {
		t.Error("Expected a command to be returned")
	}
}

func TestEmailViewerModel_EmailLoadedMessage_Success(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)
	model.loading = true

	testEmail := &core.Email{
		EmailMetadata: core.EmailMetadata{
			Id:      "test-123",
			From:    "sender@example.com",
			To:      "receiver@example.com",
			Subject: "Test Subject",
			SentAt:  time.Now(),
		},
		Body: "This is a test email body",
	}

	updatedModel, _ := model.Update(EmailLoadedMessage{
		Email: testEmail,
		Error: nil,
	})

	if updatedModel.loading {
		t.Error("Expected loading to be false after successful EmailLoadedMessage")
	}

	if updatedModel.email == nil {
		t.Error("Expected email to be set after successful load")
	}

	if updatedModel.email.Subject != "Test Subject" {
		t.Errorf("Expected subject 'Test Subject', got '%s'", updatedModel.email.Subject)
	}

	if updatedModel.error != "" {
		t.Errorf("Expected empty error, got '%s'", updatedModel.error)
	}
}

func TestEmailViewerModel_EmailLoadedMessage_Error(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)
	model.loading = true

	testError := errors.New("failed to load email")

	updatedModel, _ := model.Update(EmailLoadedMessage{
		Email: nil,
		Error: testError,
	})

	if updatedModel.loading {
		t.Error("Expected loading to be false after error EmailLoadedMessage")
	}

	if updatedModel.error != "failed to load email" {
		t.Errorf("Expected error 'failed to load email', got '%s'", updatedModel.error)
	}

	if updatedModel.email != nil {
		t.Error("Expected email to remain nil after error")
	}
}

func TestEmailViewerModel_WindowSizeMsg_FirstTime(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)

	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, _ := model.Update(windowMsg)

	if !updatedModel.ready {
		t.Error("Expected ready to be true after first WindowSizeMsg")
	}
}

func TestEmailViewerModel_WindowSizeMsg_SubsequentTimes(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)
	model.ready = true

	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := model.Update(windowMsg)

	if !updatedModel.ready {
		t.Error("Expected ready to remain true after subsequent WindowSizeMsg")
	}
}

func TestEmailViewerModel_KeyMsg_Quit(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(keyMsg)

	// Should return a command to show email list
	if cmd == nil {
		t.Error("Expected a command to be returned for 'q' key")
	}
}

func TestEmailViewerModel_View_NotReady(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)
	model.ready = false

	view := model.View()

	if view != "Loading viewer..." {
		t.Errorf("Expected 'Loading viewer...', got '%s'", view)
	}
}

func TestEmailViewerModel_View_Loading(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)
	model.ready = true
	model.loading = true

	view := model.View()

	if view != "Loading email..." {
		t.Errorf("Expected 'Loading email...', got '%s'", view)
	}
}

func TestEmailViewerModel_View_Error(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)
	model.ready = true
	model.loading = false
	model.error = "Connection timeout"

	view := model.View()

	expected := "Error: Connection timeout"
	if view != expected {
		t.Errorf("Expected '%s', got '%s'", expected, view)
	}
}

func TestEmailViewerModel_View_NormalState(t *testing.T) {
	backend := &mockBackend{}
	model := NewEmailViewerModel(backend)

	// Initialize the viewport by sending a window size message
	model, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	model.loading = false
	model.error = ""

	view := model.View()

	if strings.Contains(view, "Loading") || strings.Contains(view, "Error") {
		t.Error("Expected viewport view, got loading or error message")
	}

	// The view should not be empty when ready
	if view == "" {
		t.Error("Expected non-empty viewport view")
	}
}
