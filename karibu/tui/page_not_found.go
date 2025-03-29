package tui

import tea "github.com/charmbracelet/bubbletea"

type PageNotFound struct {
}

func (pnf PageNotFound) Init() tea.Cmd {
	return nil
}

func (pnf PageNotFound) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return pnf, nil
}

func (pnf PageNotFound) View() string {
	return "Page not found"
}
