package ui

import (
	"http_client/cmd"
	"http_client/internal/clipboard"
	component "http_client/ui/components"
	"http_client/ui/events"
	"http_client/ui/layout"
	"http_client/ui/shotcust"
	"http_client/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AppSetup maneja la configuración inicial de la aplicación
type AppSetup struct {
	appState   *AppState
	tabManager *TabManager
}

// NewAppSetup crea un nuevo configurador de aplicación
func NewAppSetup(appState *AppState) *AppSetup {
	return &AppSetup{
		appState:   appState,
		tabManager: NewTabManager(appState),
	}
}

// SetupPages configura todas las páginas de la aplicación
func (setup *AppSetup) SetupPages(main *layout.Layout, cli *cmd.CliOptions) {
	setup.appState.SetMainPage(tview.NewPages())
	tabList := component.TabsList()
	setup.appState.SetTabList(tabList)
	setup.appState.tabList = tabList.List

	// Configurar pestaña inicial
	setup.tabManager.CreateInitialTab(setup.appState.currentTab)
	tabList.List.AddItem("Not Found", setup.appState.currentTab, 0, nil)
	setup.tabManager.currentTabId = &setup.appState.currentTab

	// ColorBackground tcell.Color = tcell.GetColor("#1e1e2e") // Base
	setup.appState.mainPage.SetBackgroundColor(tcell.GetColor("#1e1e2e"))

	// Agregar páginas
	setup.appState.mainPage.AddPage("main", main.Main, true, true)
	setup.appState.mainPage.AddPage("help", component.Help(), true, false)
	setup.appState.mainPage.AddPage("tab", tabList.Container, true, false)

	if cli.ActHistory {
		setup.appState.mainPage.AddPage("history", component.History(setup.appState.app), true, false)
	}
}

// SetupRequestSender configura el enviador de peticiones
func (setup *AppSetup) SetupRequestSender() *events.RequestSender {
	return &events.RequestSender{
		App:           setup.appState.app,
		ResponseView:  setup.appState.responseView,
		ResponseInfo:  setup.appState.responseInfo,
		ContentToCopy: &setup.appState.contentToCopy,
	}
}

// SetupShortcuts configura los atajos de teclado
func (setup *AppSetup) SetupShortcuts(main *layout.Layout, cli *cmd.CliOptions) {
	requestSender := setup.SetupRequestSender()

	shortcuts := &shotcust.Shortcuts{
		App:                    setup.appState.app,
		MainPage:               setup.appState.mainPage,
		SwitchPage:             main.EditorPanel.Editor,
		ResponseFlex:           main.RightPanel.Container,
		ShowRequestHistoryPage: setup.appState.showHelpPage,
		ShowHelpPage:           setup.appState.showHelpPage,
		RightPanel:             main.RightPanel.Container,
		ChangeToFullScreen:     setup.appState.fullScreen,
		ShowTabPage:            &setup.appState.showTabPage,
		ExitsTab: func() bool {
			_, exists := setup.appState.tabs[setup.appState.currentTab]
			utils.LogDebug("Hi")
			return exists
		},

		CopyFn: func() {
			clipboard.Copy(setup.appState.contentToCopy)
		},
		ResponseFn: func() {
			requestSender.SendRequest(
				main.LeftPanel.Input,
				main.LeftPanel.DropDown,
				main.EditorPanel.Body,
				main.EditorPanel.Header,
				main.EditorPanel.QueryParam,
				main.EditorPanel.PathParam,
				main.EditorPanel.Variable,
				cli.ActHistory,
			)

		},
		FocusForm: func() {
			setup.appState.app.SetFocus(main.Form.Container)
		},
		SaveStateFn: func() {
			setup.tabManager.SaveCurrentTabState(main)
		},
	}

	shortcuts.RegisterKeys()
}

// SetupEventHandlers configura todos los manejadores de eventos
func (setup *AppSetup) SetupEventHandlers(main *layout.Layout) {
	setup.tabManager.SetupTabListHandlers(main)
}
