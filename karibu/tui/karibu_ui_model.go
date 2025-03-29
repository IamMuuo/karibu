package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type KaribuUiModel struct {
	pages        map[string]tea.Model
	selectedPage string
}

func NewKaribuUiModel() KaribuUiModel {
	kum := KaribuUiModel{}
	kum.pages = make(map[string]tea.Model)
	return kum
}

func (kum *KaribuUiModel) Init() tea.Cmd {
	return nil

}

func (kum *KaribuUiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch event := msg.(type) {
	case tea.KeyMsg:
		switch event.String() {

		case "ctrl+c", "esc", "q":
			return kum, tea.Quit
		}
	case tea.WindowSizeMsg:
		log.Info("Window resize event!")

	}

	_, cmd := kum.CurrentModal().Update(msg)

	return kum, cmd
}

func (kum *KaribuUiModel) View() string {
	current, exists := kum.pages[kum.selectedPage]
	if !exists {
		return (&PageNotFound{}).View()
	}

	return current.View()
}

func (kum *KaribuUiModel) CurrentModal() tea.Model {
	if current, exists := kum.pages[kum.selectedPage]; exists {
		return current
	}

	return &PageNotFound{}
}

func (kum *KaribuUiModel) CurrentModalName() string {
	return kum.selectedPage
}

func (kum *KaribuUiModel) SwitchModal(name string) {
	if _, exists := kum.pages[name]; exists {
		kum.selectedPage = name
	}
}
