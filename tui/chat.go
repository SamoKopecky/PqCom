package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	smallStep    = 1
	bigStep      = 10
	editorHeight = 3
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f2f2f2"))

	displayBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#666666"))

	textAreaBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#ffa31a"))
)

type model struct {
	width    int
	height   int
	input    textarea.Model
	display  viewport.Model
	messages []string
	send     chan []byte
	recv     chan []byte
	help     help.Model
}

func NewModel(send, recv chan []byte) model {
	m := model{
		input: textarea.Model{},
		send:  send,
		recv:  recv,
		help:  help.New(),
	}
	// w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	// m.display = viewport.New(0, 0)
	// m.display.Style = displayBorderStyle
	// a := displayBorderStyle.Copy().Width(w).Height(h - editorHeight - 5).Render("")
	// m.display.SetContent(a)

	m.input = textarea.New()
	m.input.Prompt = ""
	m.input.Placeholder = "Send a message..."
	m.input.ShowLineNumbers = false
	m.input.Cursor.Style = cursorStyle
	m.input.FocusedStyle.CursorLine = lipgloss.NewStyle()
	m.input.FocusedStyle.Base = textAreaBorderStyle
	m.input.KeyMap.DeleteWordBackward.SetEnabled(false)
	m.input.KeyMap.InsertNewline.SetEnabled(false)
	m.input.Focus()
	return m
}

func (m model) Init() tea.Cmd {
	// m.display = viewport.New(0, 0)
	// m.display.Style = displayBorderStyle
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	select {
	case <-m.send:
		m.messages = append(m.messages, string(<-m.send))
		a := displayBorderStyle.Copy().Width(m.width - 2).Height(m.height - m.input.Height() - 5).Render(strings.Join(m.messages, "\n"))
		m.display.SetContent(a)
	default:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, ("You: ")+m.input.Value())
			m.input.Reset()
			// for i, msg := range m.messages {

			// 	if len(msg) > m.display.Width-2 {
			// 		// fmt.Printf("##\n")
			// 		one := m.messages[i][:m.display.Width-2]
			// 		two := fmt.Sprintf("\n%s", m.messages[i][m.display.Width-2:])
			// 		m.messages[i] = fmt.Sprintf("%s%s", one, two)
			// 	}
			// }
			a := displayBorderStyle.Copy().Width(m.width - 2).Height(m.height - m.input.Height() - 5).Render(strings.Join(m.messages, "\n"))
			m.display.SetContent(a)

			// m.display.SetContent(strings.Join(m.messages, "\n"))
			m.display.GotoBottom()
		case tea.KeyUp:
			m.display.LineUp(smallStep)
		case tea.KeyDown:
			m.display.LineDown(smallStep)
		case tea.KeyShiftUp:
			m.display.LineUp(bigStep)
		case tea.KeyShiftDown:
			m.display.LineDown(bigStep)
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}
	m.sizeInputs()

	newModel, cmd := m.input.Update(msg)
	m.input = newModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) sizeInputs() {
	m.input.SetWidth(m.width)
	m.input.SetHeight(editorHeight)
	// a := displayBorderStyle.Copy().Width(m.width).Height(m.height - m.input.Height() - 3).Render("")
	m.display.Width = m.width
	m.display.Height = m.height - m.input.Height() - 3
	// m.display.SetContent(a)

}

func (m model) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	})

	return m.display.View() + "\n" + m.input.View() + "\n" + help
}
