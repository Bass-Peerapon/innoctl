package filepicker

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	filepicker   filepicker.Model
	title        string
	selectedFile *string
	quitting     *bool
	exit         *bool
	err          error
}

func (m model) IsExit() bool {
	return *m.exit
}

func (m model) GetSelectedFile() string {
	return *m.selectedFile
}

type clearErrorMsg struct{}

var (
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	titleStyle        = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
)

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func InitialFilepickerModel() model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".go"}
	fp.AutoHeight = false
	fp.Height = 6
	return model{
		filepicker:   fp,
		selectedFile: new(string),
		quitting:     new(bool),
		exit:         new(bool),
		err:          nil,
	}
}

func (m *model) SetTitle(title string) {
	m.title = title
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			*m.quitting = true
			*m.exit = true
			return m, tea.Quit
		case "y":
			*m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		*m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		*m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m model) View() string {
	if *m.quitting {
		return ""
	}
	var s strings.Builder
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	}
	s.WriteString(fmt.Sprintf(titleStyle.Render("Pick a %s file"), m.title) + "\n\n")
	s.WriteString("Selected file: " + selectedItemStyle.Render(*m.selectedFile) + "\n\n")
	s.WriteString(m.filepicker.View() + "\n\n")
	s.WriteString(fmt.Sprintf("Press %s to confirm selected.", focusedStyle.Render("y")))

	return s.String()
}
