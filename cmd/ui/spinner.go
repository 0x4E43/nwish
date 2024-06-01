package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinerModel struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

type errMsg error

func SpinnerModel() SpinerModel {
	spn := spinner.New()
	spn.Spinner = spinner.Dot
	spn.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return SpinerModel{spinner: spn}
}

func (spn SpinerModel) Init() tea.Cmd {
	return spn.spinner.Tick
}

func (spn SpinerModel) View() string {
	if spn.err != nil {
		return spn.err.Error()
	}

	str := fmt.Sprintf("\n\n%s Content loaded...\n\npress c to continue \npress q to quit\n\n", spn.spinner.View())
	if spn.quitting {
		return str + "\n"
	}

	return str
}

func (spn SpinerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			spn.quitting = true
			return spn, tea.Quit
		case "c", "C":
			return RenderMain(), tea.Quit
		default:
			return spn, nil
		}

	case errMsg:
		spn.err = msg
		return spn, nil

	default:
		var cmd tea.Cmd
		spn.spinner, cmd = spn.spinner.Update(msg)
		return spn, cmd
	}
}
