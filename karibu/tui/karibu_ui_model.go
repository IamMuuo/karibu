package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

type KaribuUiModel struct {
	pages        map[string]tea.Model
	selectedPage string
	session      *ssh.Session
	pty          *ssh.Pty
}

func NewKaribuUiModel(session *ssh.Session, pty *ssh.Pty) KaribuUiModel {
	kum := KaribuUiModel{}
	kum.session = session
	kum.pty = pty
	kum.pages = make(map[string]tea.Model)

	// add the necessary page routes
	kum.pages["/splashscreen"] = SplashScreen{session: session, pty: pty}
	kum.pages["/not-found"] = PageNotFound{session: session, pty: pty}

	return kum
}

func (kum *KaribuUiModel) Init() tea.Cmd {
	return kum.SwitchModal("/splashscreen")
}

func (kum *KaribuUiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch event := msg.(type) {
	case tea.KeyMsg:
		switch event.String() {

		case "ctrl+c", "esc", "q":
			return kum, tea.Quit
		}
	case tea.WindowSizeMsg:
		// log.Info("Window resize event!")

	case SwitchPageMsg:
		log.Infof("Switching to page %s", event.TargetPage)
		return kum, kum.SwitchModal(event.TargetPage)
	}

	_, cmd := kum.CurrentModal().Update(msg)

	return kum, cmd
}

func (kum *KaribuUiModel) View() string {
	current, exists := kum.pages[kum.selectedPage]
	if !exists {

		return PageNotFound{session: kum.session, pty: kum.pty}.View()
	}

	return current.View()
}

func (kum *KaribuUiModel) CurrentModal() tea.Model {
	if current, exists := kum.pages[kum.selectedPage]; exists {
		return current
	}

	return &PageNotFound{session: kum.session, pty: kum.pty}

}

func (kum *KaribuUiModel) CurrentModalName() string {
	return kum.selectedPage
}

func (kum *KaribuUiModel) SwitchModal(name string) tea.Cmd {
	if page, exists := kum.pages[name]; exists {
		kum.selectedPage = name
		return page.Init()
	}
	kum.selectedPage = name
	return nil
}
