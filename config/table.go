package config

import (
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#666666"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func NewTable(kems []string, signs []string) {
	columns := []table.Column{
		{Title: "Key Encapsulation (kem_alg)aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Width: 30},
		{Title: "Digital Signature (sign_alg)", Width: 30},
	}
	rows := []table.Row{}

	rowsC := 0
	if len(kems) < len(signs) {
		rowsC = len(signs)
		for i := 0; i < len(signs)-len(kems); i++ {
			kems = append(kems, "")
		}
	} else {
		rowsC = len(kems)
		for i := 0; i < len(kems)-len(signs); i++ {
			signs = append(signs, "")
		}
	}
	for i := 0; i < rowsC; i++ {
		rows = append(rows, table.Row{kems[i], signs[i]})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#666666")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#666666")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Error().Str("error", err.Error()).Msg("Error while listing algorithms")
		os.Exit(1)
	}
}
