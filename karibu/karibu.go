package karibu

import tea "github.com/charmbracelet/bubbletea"

type Karibu struct {
}

func (k Karibu) Init() tea.Cmd {
	return nil
}

func (k Karibu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return k, tea.Quit
		}
	}
	return k, nil
}

func (k Karibu) View() string {
	return "Hello, Bubble Tea! Press 'q' to quit."
}
