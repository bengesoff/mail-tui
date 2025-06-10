package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bengesoff/mail-tui/internal/backend/fake"
	"github.com/bengesoff/mail-tui/internal/ui"
)

func TestModel_InitialState(t *testing.T) {
	m := NewAppModel(&fake.FakeBackend{})

	if m.activeView != ListViewName {
		t.Errorf("Expected initial view to be '%s', got '%s'", ListViewName, m.activeView)
	}

	if m.emailList == nil {
		t.Error("Expected emailList to be initialized")
	}

	if m.emailViewer == nil {
		t.Error("Expected emailViewer to be initialized")
	}

	if m.emailComposer == nil {
		t.Error("Expected emailComposer to be initialized")
	}
}

func TestModel_Init(t *testing.T) {
	m := NewAppModel(&fake.FakeBackend{})
	cmd := m.Init()

	if cmd == nil {
		t.Error("Expected Init to return a command")
	}

	// Execute the command to verify it returns ShowEmailListMessage
	msg := cmd()
	_, ok := msg.(ui.ShowEmailListMessage)
	if !ok {
		t.Error("Expected Init command to return ShowEmailListMessage")
	}
}

func TestModel_Update_ShowEmailListMessage(t *testing.T) {
	m := NewAppModel(&fake.FakeBackend{})
	m.activeView = ViewerViewName // Start with viewer view

	updatedModel, _ := m.Update(ui.ShowEmailListMessage{})
	updated := updatedModel.(AppModel)

	if updated.activeView != ListViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailListMessage, got '%s'", ListViewName, updated.activeView)
	}
}

func TestModel_Update_ShowEmailViewerMessage(t *testing.T) {
	m := NewAppModel(&fake.FakeBackend{})

	updatedModel, _ := m.Update(ui.ShowEmailViewerMessage{EmailId: "test-id"})
	updated := updatedModel.(AppModel)

	if updated.activeView != ViewerViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailViewerMessage, got '%s'", ViewerViewName, updated.activeView)
	}
}

func TestModel_Update_ShowEmailComposerMessage(t *testing.T) {
	m := NewAppModel(&fake.FakeBackend{})

	updatedModel, _ := m.Update(ui.ShowEmailComposerMessage{})
	updated := updatedModel.(AppModel)

	if updated.activeView != ComposerViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailComposerMessage, got '%s'", ComposerViewName, updated.activeView)
	}
}

func TestModel_Update_CtrlC(t *testing.T) {
	m := NewAppModel(&fake.FakeBackend{})

	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := m.Update(keyMsg)

	// Should return quit command
	if cmd == nil {
		t.Error("Expected Ctrl+C to return a quit command")
	}
}

func TestModel_ViewSwitching_Sequence(t *testing.T) {
	m := NewAppModel(&fake.FakeBackend{})

	// Should start with list view
	if m.activeView != ListViewName {
		t.Errorf("Expected initial view to be '%s', got '%s'", ListViewName, m.activeView)
	}

	// Switch to viewer
	um, _ := m.Update(ui.ShowEmailViewerMessage{EmailId: "test"})
	updatedModel := um.(AppModel)

	if updatedModel.activeView != ViewerViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailViewerMessage, got '%s'", ViewerViewName, updatedModel.activeView)
	}

	// Switch back to list
	um, _ = m.Update(ui.ShowEmailListMessage{})
	updatedModel = um.(AppModel)

	if updatedModel.activeView != ListViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailListMessage, got '%s'", ListViewName, updatedModel.activeView)
	}

	// Switch to viewer again
	um, _ = m.Update(ui.ShowEmailViewerMessage{EmailId: "another-test"})
	updatedModel = um.(AppModel)

	if updatedModel.activeView != ViewerViewName {
		t.Errorf("Expected view to be '%s' after second ShowEmailViewerMessage, got '%s'", ViewerViewName, updatedModel.activeView)
	}

	// Switch to composer
	um, _ = m.Update(ui.ShowEmailComposerMessage{})
	updatedModel = um.(AppModel)

	if updatedModel.activeView != ComposerViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailComposerMessage, got '%s'", ComposerViewName, updatedModel.activeView)
	}

	// Switch back to list from composer
	um, _ = m.Update(ui.ShowEmailListMessage{})
	updatedModel = um.(AppModel)

	if updatedModel.activeView != ListViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailListMessage from composer, got '%s'", ListViewName, updatedModel.activeView)
	}
}
