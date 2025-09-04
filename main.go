package main

import (
	"database/sql"
	"http_client/ui"
	"http_client/utils"
	"os"
)

func createTableIfNotExist() error {
	table := `
CREATE TABLE IF NOT EXISTS history (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    status_code INTEGER,
    content_type TEXT,
    response_body TEXT,
    created_at TEXT NOT NULL
);`

	db, err := sql.Open("sqlite", os.Getenv("HOME")+"/.local/share/request.db")

	if err != nil {
		return err
	}

	_, err = db.Exec(table)

	if err != nil {
		return err
	}

	return nil
}
func main() {

	if err := createTableIfNotExist(); err != nil {
		utils.LogDebug("Error al crear la base de datos" + " " + err.Error())
	}

	ui.StartApp()
}
