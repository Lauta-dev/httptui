package cmd

import "fmt"

// Mostrar ayuda por terminal
func Help() {

	dialog := `
	Programa desarrollado por: https://github.com/lauta-dev

--help | -help: Mostrar esta ayuda
--env-file | -env-file: Cargar archivo .env
--activate-history | -activate-history: Activar historial de request
	`

	fmt.Print(dialog)

}
