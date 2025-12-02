package ui

import (
	"http_client/cmd"
	"http_client/logic"
	"http_client/ui/layout"
	"log"
)

// StartApp inicia la aplicación principal
func StartApp() {
	cli := cmd.Launch()
	if cli.Help {
		return
	}

	if err := logic.InitDB(); err != nil {
		log.Fatal(err)
	}

	defer logic.Close()
	logic.CreateDatabase()

	// Crear estado de la aplicación
	appState := NewAppState()

	// Configurar layout principal
	main := layout.MainLayout()
	appState.SetResponseViews(main.RightPanel.ResponseView, main.RightPanel.ResponseInfo)

	// Crear configurador de aplicación
	appSetup := NewAppSetup(appState)

	// Configurar páginas
	appSetup.SetupPages(main, &cli)

	// Configurar eventos y atajos
	appSetup.SetupShortcuts(main, &cli)
	appSetup.SetupEventHandlers(main)

	// Configurar variables de entorno
	main.EditorPanel.Variable.SetText(logic.ReadEnvFile(cli.EnvFilePath), false)

	// Configurar estilos y ejecutar
	SetupStyles()

	if err := appState.app.SetRoot(appState.mainPage, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
