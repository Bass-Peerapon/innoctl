package viewport

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bass-Peerapon/innoctl/cmd/create/dbmlstruct"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type Model struct {
	DbmlTextArea textarea.Model
	GoView       viewport.Model
	Output       *string
	err          error
}

type (
	TickMsg struct{}
	errMsg  error
)

var (
	result      string
	resultError error
)

func InitModel() Model {
	ti := textarea.New()
	ti.Placeholder = "Enter your DBML here"
	ti.Focus()
	ti.CharLimit = 0 // Unlimited
	ti.MaxHeight = 0

	return Model{
		DbmlTextArea: ti,
		Output:       new(string),
		err:          nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.DbmlTextArea.Focused() {
				m.DbmlTextArea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlD:
			m.DbmlTextArea.SetValue("")
		case tea.KeyCtrlB:
			// copy to clipboard
			if resultError != nil {
				clipboard.WriteAll(resultError.Error())
			} else {
				clipboard.WriteAll(result)
			}

			m.GoView.Style = m.GoView.Style.Background(lipgloss.Color("#808080"))
			return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
				return TickMsg{}
			})

		default:
			if !m.DbmlTextArea.Focused() {
				cmd = m.DbmlTextArea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case TickMsg:
		m.GoView.Style = m.GoView.Style.Background(lipgloss.NoColor{})
		return m, nil

	case tea.WindowSizeMsg:
		headerGoHeight := lipgloss.Height(m.headerGoView())
		footerGoHeight := lipgloss.Height(m.footerGoView())
		verticalMarginGoHeight := headerGoHeight + footerGoHeight
		headerDbmlHeigt := lipgloss.Height(m.headerDbmlView())
		footerDbmlHeight := lipgloss.Height(m.footerDbmlView())

		m.DbmlTextArea.SetHeight(msg.Height - headerDbmlHeigt - footerDbmlHeight)
		m.DbmlTextArea.SetWidth(msg.Width / 2)

		dbmlW := lipgloss.Width(m.DbmlTextArea.View())
		m.GoView = viewport.New(msg.Width-dbmlW, msg.Height-verticalMarginGoHeight)
		m.GoView.YPosition = headerGoHeight

	case errMsg:
		m.err = msg
		return m, nil
	}
	result, resultError = dbmlstruct.GenerateFormString(m.DbmlTextArea.Value())
	if resultError != nil {
		m.GoView.SetContent(resultError.Error())
	} else {
		m.GoView.SetContent(result)
	}
	m.DbmlTextArea, cmd = m.DbmlTextArea.Update(msg)
	cmds = append(cmds, cmd)

	// Handle keyboard and mouse events in the viewport
	m.GoView, cmd = m.GoView.Update(msg)
	cmds = append(cmds, cmd)

	m.GoView.Style = m.GoView.Style.Background(lipgloss.NoColor{})
	return m, tea.Batch(cmds...)
}

func (m Model) headerGoView() string {
	title := titleStyle.Render("Go Struct Generator")
	line := strings.Repeat("─", max(0, m.GoView.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerGoView() string {
	info := infoStyle.Render("Press Ctrl+B to copy to clipboard")
	line := strings.Repeat("─", max(0, m.GoView.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Model) headerDbmlView() string {
	title := titleStyle.Render("Dbml Text Area")
	line := strings.Repeat("─", max(0, m.GoView.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerDbmlView() string {
	info := infoStyle.Render("Press Ctrl+C to exit, Ctrl+D to clear")
	line := strings.Repeat("─", max(0, m.DbmlTextArea.Width()-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Model) View() string {
	var sb strings.Builder
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%s\n%s\n%s", m.headerDbmlView(), m.DbmlTextArea.View(), m.footerDbmlView()), fmt.Sprintf("%s\n%s\n%s", m.headerGoView(), m.GoView.View(), m.footerGoView())))
	return sb.String()
}
