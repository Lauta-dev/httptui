package logic

import (
	"database/sql"
	"errors"
	"fmt"
	"http_client/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

var sqliteFile = filepath.Join(os.Getenv("HOME"), ".local", "share", "request.db")
var DB *sql.DB

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

func InitDB() error {
	d, err := sql.Open("sqlite", sqliteFile)
	if err != nil {
		return err
	}

	if err := d.Ping(); err != nil {
		return err
	}

	DB = d
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func CreateDatabase() (bool, error) {
	q := `
	CREATE TABLE request_history (
  id TEXT PRIMARY KEY,
  url TEXT,
  method TEXT,
  status_code INTEGER,
  content_type TEXT,
  response_body TEXT,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);
	`

	_, err := DB.Exec(q)

	if err != nil {
		return false, err
	}

	return true, nil
}

func DelItems(id string) (string, error) {
	_, err := DB.Exec("DELETE FROM request_history WHERE id = ?", id)

	if err != nil {
		return "", err
	}

	return "Eliminado", nil
}

func GetAllItems() ([]history, error) {
	// Leer y mostrar
	rows, err := DB.Query("SELECT id, url, method, status_code FROM request_history ORDER BY created_at asc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // CERRAR SIEMPRE

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

	if len(v) == 0 {
		return nil, errors.New("No hay elementos")
	}

	return v, nil
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

	_, err := DB.Exec(insert, id, url, method, code, contentType, body.ToDisplay, createAt)

	if err != nil {
		return err
	}

	return nil
}

func GetItemById(iid string) (history, error) {

	var f history
	var id, url, method, body, statusCode, contentType, createAt string

	row := DB.QueryRow("SELECT id, url, method, status_code, content_type, response_body, created_at FROM request_history WHERE id = ?", iid)

	if err := row.Scan(&id, &url, &method, &statusCode, &contentType, &body, &createAt); err != nil {
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
