package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bengesoff/mail-tui/internal/ui"
)

func TestModel_InitialState(t *testing.T) {
	m := initialModel()

	if m.activeView != ListViewName {
		t.Errorf("Expected initial view to be '%s', got '%s'", ListViewName, m.activeView)
	}

	if m.emailList == nil {
		t.Error("Expected emailList to be initialized")
	}

	if m.emailViewer == nil {
		t.Error("Expected emailViewer to be initialized")
	}
}

func TestModel_Init(t *testing.T) {
	m := initialModel()
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
	m := initialModel()
	m.activeView = ViewerViewName // Start with viewer view

	updatedModel, _ := m.Update(ui.ShowEmailListMessage{})
	updated := updatedModel.(model)

	if updated.activeView != ListViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailListMessage, got '%s'", ListViewName, updated.activeView)
	}
}

func TestModel_Update_ShowEmailViewerMessage(t *testing.T) {
	m := initialModel()

	updatedModel, _ := m.Update(ui.ShowEmailViewerMessage{EmailId: "test-id"})
	updated := updatedModel.(model)

	if updated.activeView != ViewerViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailViewerMessage, got '%s'", ViewerViewName, updated.activeView)
	}
}

func TestModel_Update_CtrlC(t *testing.T) {
	m := initialModel()

	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := m.Update(keyMsg)

	// Should return quit command
	if cmd == nil {
		t.Error("Expected Ctrl+C to return a quit command")
	}
}

func TestModel_ViewSwitching_Sequence(t *testing.T) {
	m := initialModel()

	// Should start with list view
	if m.activeView != ListViewName {
		t.Errorf("Expected initial view to be '%s', got '%s'", ListViewName, m.activeView)
	}

	// Switch to viewer
	updatedModel, _ := m.Update(ui.ShowEmailViewerMessage{EmailId: "test"})
	m = updatedModel.(model)

	if m.activeView != ViewerViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailViewerMessage, got '%s'", ViewerViewName, m.activeView)
	}

	// Switch back to list
	updatedModel, _ = m.Update(ui.ShowEmailListMessage{})
	m = updatedModel.(model)

	if m.activeView != ListViewName {
		t.Errorf("Expected view to be '%s' after ShowEmailListMessage, got '%s'", ListViewName, m.activeView)
	}

	// Switch to viewer again
	updatedModel, _ = m.Update(ui.ShowEmailViewerMessage{EmailId: "another-test"})
	m = updatedModel.(model)

	if m.activeView != ViewerViewName {
		t.Errorf("Expected view to be '%s' after second ShowEmailViewerMessage, got '%s'", ViewerViewName, m.activeView)
	}
}
