package logic

import (
	"database/sql"
	"fmt"
	"http_client/utils"
	"log"
	"os"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

var (
	sqliteFile string = os.Getenv("HOME") + ".local/share/request.db"
)

// [GET, 200] URL
// ID
type history struct {
	ID           string
	URL          string
	Method       string
	StatusCode   int
	ContentType  string
	ResponseBody string
	CreatedAt    string
}

func DelItems(id string) (string, error) {
	db, err := sql.Open("sqlite", sqliteFile)

	if err != nil {
		return "", err
	}

	_, err = db.Exec("DELETE FROM request_history WHERE id = ?", id)

	if err != nil {
		return "", err
	}

	return "Eliminado", nil
}

func GetAllItems() []history {
	db, err := sql.Open("sqlite", sqliteFile)
	if err != nil {
		log.Fatal(err)
	}

	// Leer y mostrar
	rows, err := db.Query("SELECT id, url, method, status_code FROM request_history ORDER BY created_at asc")
	if err != nil {
		log.Fatal(err)
	}

	var v []history

	for rows.Next() {
		var id, url, method string
		var code int
		if err := rows.Scan(&id, &url, &method, &code); err != nil {
			log.Fatal(err)
		}

		v = append(v, history{
			ID:         id,
			URL:        url,
			Method:     method,
			StatusCode: code,
		})
	}

	return v
}

func SaveItems(url string, code string, contentType string, responseBody string, method string) error {
	createAt := time.Now().Format(time.RFC3339)
	id := createAt
	body := utils.IdentText([]byte(responseBody), contentType)

	insert := `
	INSERT INTO request_history (
		id, url, method, status_code, content_type, response_body, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	db, err := sql.Open("sqlite", sqliteFile)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(insert, id, url, method, code, contentType, body.ToDisplay, createAt)

	if err != nil {
		return err
	}

	return nil
}

func GetItemById(iid string) (history, error) {
	db, err := sql.Open("sqlite", sqliteFile)

	var f history
	var id, url, method, body, statusCode, contentType, createAt string

	if err != nil {
		return f, fmt.Errorf("error al abrir base de datos: %w", err)
	}

	row := db.QueryRow("SELECT id, url, method, status_code, content_type, response_body, created_at FROM request_history WHERE id = ?", iid)

	err = row.Scan(&id, &url, &method, &statusCode, &contentType, &body, &createAt)

	if err != nil {
		return f, fmt.Errorf("error al escanear fila: %w", err)
	}

	d, err := strconv.Atoi(statusCode)
	if err != nil {
		return f, fmt.Errorf("error al pasar de string a int: %w", err)
	}

	f.ID = id
	f.Method = method
	f.URL = url
	f.ResponseBody = body
	f.StatusCode = d
	f.ContentType = contentType
	f.CreatedAt = createAt

	return f, nil
}
