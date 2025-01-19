package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textInput textinput.Model
	name      string
	done      bool
}

// Style definitions
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	questionStyle = lipgloss.NewStyle().
			Foreground(special).
			Bold(true)

	responseStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true)
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter your name"
	ti.Focus()
	ti.CharLimit = 30
	ti.Width = 20

	return model{
		textInput: ti,
		done:      false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.name = m.textInput.Value()
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.done {
		return dialogBoxStyle.Render(
			responseStyle.Render(fmt.Sprintf("Nice to meet you, %s! 👋", m.name)),
		) + "\n"
	}

	return dialogBoxStyle.Render(
		questionStyle.Render("What's your name?") + "\n\n" +
			m.textInput.View(),
	) + "\n"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}