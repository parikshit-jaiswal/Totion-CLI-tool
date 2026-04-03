package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/30, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type model struct {
	msg      string
	cursor   int
	choices  []string
	selected map[int]struct{}

	fullTitle     string
	visibleRunes  int
	typingDone    bool
	frames        int
	blinkOn       bool
	cursorTarget  float64
	cursorOpacity float64
	cursorVel     float64
	spring        harmonica.Spring
}

func initializeModel() model {
	return model{
		msg:      "Hello, World!",
		cursor:   0,
		choices:  []string{"Add task", "View tasks", "Quit"},
		selected: make(map[int]struct{}),

		fullTitle:    "Welcome to the Totion",
		visibleRunes: 0,
		typingDone:   false,
		frames:       0,
		blinkOn:      true,
		cursorTarget: 1,
		spring:       harmonica.NewSpring(harmonica.FPS(30), 8.0, 0.6),
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch keyMsg := msg.(type) {
	case tickMsg:
		titleLen := len([]rune(m.fullTitle))
		if !m.typingDone && m.frames%2 == 0 {
			m.visibleRunes++
			if m.visibleRunes >= titleLen {
				m.visibleRunes = titleLen
				m.typingDone = true
			}
		}

		if m.frames%12 == 0 {
			m.blinkOn = !m.blinkOn
		}

		if m.blinkOn {
			m.cursorTarget = 1
		} else {
			m.cursorTarget = 0
		}

		m.cursorOpacity, m.cursorVel = m.spring.Update(m.cursorOpacity, m.cursorVel, m.cursorTarget)
		m.frames++
		return m, tickCmd()

	case tea.KeyPressMsg:
		switch keyMsg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", "space":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil

}

func (m model) View() tea.View {
	titleRunes := []rune(m.fullTitle)
	if m.visibleRunes > len(titleRunes) {
		m.visibleRunes = len(titleRunes)
	}

	visibleTitle := string(titleRunes[:m.visibleRunes])

	cursorColor := lipgloss.Color("#2F3E46")
	if m.cursorOpacity > 0.65 {
		cursorColor = lipgloss.Color("#84A98C")
	} else if m.cursorOpacity > 0.35 {
		cursorColor = lipgloss.Color("#52796F")
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FF69B4"))
	cursorStyle := lipgloss.NewStyle().Bold(true).Foreground(cursorColor)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#52796F"))
	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#84A98C")).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#CAD2C5"))

	var b strings.Builder
	b.WriteString(titleStyle.Render(visibleTitle + cursorStyle.Render("▌")))
	b.WriteString("\n\n")

	for i, choice := range m.choices {
		pointer := " "
		if i == m.cursor {
			pointer = ">"
		}

		selected := " "
		if _, ok := m.selected[i]; ok {
			selected = "x"
		}

		line := fmt.Sprintf("%s [%s] %s", pointer, selected, choice)
		if i == m.cursor {
			b.WriteString(focusedStyle.Render(line))
		} else {
			b.WriteString(normalStyle.Render(line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(hintStyle.Render("j/k or arrows to move • space/enter to toggle • q to quit"))

	return tea.NewView(b.String())
}

func main() {
	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
	}
	os.Exit(0)
}
