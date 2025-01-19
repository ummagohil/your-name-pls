package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	nameScreen screen = iota
	menuScreen
	todoScreen
	jokeScreen
	quoteScreen
)

type model struct {
	textInput    textinput.Model
	spinner      spinner.Model
	name         string
	currentTodo  string
	todos        []string
	currentScreen screen
	err          error
}

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	danger    = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF4040"}

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(special).
			Bold(true).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(highlight).
			SetString("► ")

	errorStyle = lipgloss.NewStyle().
			Foreground(danger).
			Bold(true)
)

// Sample data
var (
	jokes = []string{
		"Why do programmers prefer dark mode? Because light attracts bugs!",
		"What's a programmer's favorite place? The foo bar!",
		"Why did the programmer quit his job? Because he didn't get arrays!",
	}
	
	quotes = []string{
		"First, solve the problem. Then, write the code. - John Johnson",
		"Code is like humor. When you have to explain it, it's bad. - Cory House",
		"The best error message is the one that never shows up. - Thomas Fuchs",
	}
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter your name"
	ti.Focus()
	ti.CharLimit = 30
	ti.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(highlight)

	return model{
		textInput:    ti,
		spinner:      s,
		currentScreen: nameScreen,
		todos:        make([]string, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		
		case "enter":
			switch m.currentScreen {
			case nameScreen:
				m.name = m.textInput.Value()
				m.currentScreen = menuScreen
				return m, nil
			
			case todoScreen:
				if m.textInput.Value() != "" {
					m.todos = append(m.todos, m.textInput.Value())
					m.textInput.Reset()
				}
			}
		
		case "esc":
			if m.currentScreen != nameScreen {
				m.currentScreen = menuScreen
				return m, nil
			}
		
		case "1":
			if m.currentScreen == menuScreen {
				m.currentScreen = todoScreen
				m.textInput.Placeholder = "Add a todo item"
				m.textInput.Focus()
				return m, textinput.Blink
			}
		
		case "2":
			if m.currentScreen == menuScreen {
				m.currentScreen = jokeScreen
				return m, nil
			}
		
		case "3":
			if m.currentScreen == menuScreen {
				m.currentScreen = quoteScreen
				return m, nil
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.currentScreen {
	case nameScreen:
		return dialogBoxStyle.Render(
			titleStyle.Render("Welcome! What's your name?") + "\n\n" +
				m.textInput.View(),
		) + "\n"

	case menuScreen:
		return dialogBoxStyle.Render(
			titleStyle.Render(fmt.Sprintf("Hello, %s! What would you like to do?", m.name)) + "\n\n" +
				itemStyle.Render("1. Manage Todo List") + "\n" +
				itemStyle.Render("2. Get a Random Joke") + "\n" +
				itemStyle.Render("3. Get an Inspirational Quote") + "\n\n" +
				itemStyle.Render("(Press ESC to return to this menu)"),
		) + "\n"

	case todoScreen:
		var todoList string
		for i, todo := range m.todos {
			todoList += fmt.Sprintf("%d. %s\n", i+1, todo)
		}
		return dialogBoxStyle.Render(
			titleStyle.Render("Todo List") + "\n" +
				todoList + "\n" +
				m.textInput.View() + "\n\n" +
				itemStyle.Render("(Press ESC to go back)"),
		) + "\n"

	case jokeScreen:
		joke := jokes[rand.Intn(len(jokes))]
		return dialogBoxStyle.Render(
			titleStyle.Render("Here's a joke for you:") + "\n\n" +
				itemStyle.Render(joke) + "\n\n" +
				itemStyle.Render("(Press ESC to go back)"),
		) + "\n"

	case quoteScreen:
		quote := quotes[rand.Intn(len(quotes))]
		return dialogBoxStyle.Render(
			titleStyle.Render("Your inspirational quote:") + "\n\n" +
				itemStyle.Render(quote) + "\n\n" +
				itemStyle.Render("(Press ESC to go back)"),
		) + "\n"

	default:
		return "Something went wrong!"
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}