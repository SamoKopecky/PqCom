package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SamoKopecky/pqcom/main/network"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	smallStep    = 1
	newLines     = 2
	helpHeight   = 1
	bigStep      = 10
	editorHeight = 3
)

var (
	displayBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#666666"))

	textAreaBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#ffa31a"))
)

type tickMsg int
type sendFunc func([]byte, network.Type)

type model struct {
	width    int
	height   int
	help     help.Model
	input    textarea.Model
	display  viewport.Model
	messages []string
	send     sendFunc
	recv     chan []byte
	peerAddr string
}

func newModel(recv chan []byte, send sendFunc, peerAddr string) model {
	m := model{
		input:    textarea.Model{},
		send:     send,
		recv:     recv,
		help:     help.New(),
		peerAddr: peerAddr,
	}
	m.input = textarea.New()
	m.input.Prompt = ""
	m.input.Placeholder = "Send a message..."
	m.input.ShowLineNumbers = false
	m.input.FocusedStyle.CursorLine = lipgloss.NewStyle()
	m.input.FocusedStyle.Base = textAreaBorderStyle
	m.input.KeyMap.DeleteWordBackward.SetEnabled(false)
	m.input.KeyMap.InsertNewline.SetEnabled(false)
	m.input.Focus()
	return m
}

func (m *model) setDisplay() {
	var rowLen int

	for i, msg := range m.messages {
		m.messages[i] = strings.ReplaceAll(m.messages[i], "\n", "")
		// -2 cause of the borders
		rowLen = m.display.Width - 2
		if len(msg) > rowLen {
			for j := rowLen; j < len(msg); j += rowLen {
				m.messages[i] = m.messages[i][:j] + "\n" + m.messages[i][j:]
			}
		}
	}
	content := displayBorderStyle.Copy().
		Width(m.width - 2).
		Height(m.height - editorHeight - newLines - helpHeight - 2).
		Render(strings.Join(m.messages, "\n"))
	m.display.SetContent(content)
}

func tick() tea.Msg {
	time.Sleep(time.Second / 16)
	return tickMsg(1)
}

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		select {
		case data := <-m.recv:
			m.messages = append(m.messages, "["+m.peerAddr+"]: "+string(data))
			m.setDisplay()
		default:
		}
		return m, tick
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			m.send([]byte(m.input.Value()), network.ContentT)
			m.messages = append(m.messages, "[you]: "+m.input.Value())
			m.input.Reset()
			m.setDisplay()
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
		m.setDisplay()
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
	m.display.Width = m.width
	m.display.Height = m.height - editorHeight - newLines - helpHeight
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

func NewChatTui(stream network.Stream, send sendFunc) {
	recv := make(chan []byte)
	go func() {
		for {
			msg := <-stream.Msg
			recv <- msg.Data
		}
	}()
	model := newModel(recv, send, stream.Conn.RemoteAddr().String())
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
