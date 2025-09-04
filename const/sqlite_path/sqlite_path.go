package sqlitepath

import "os"

var (
	pathSql string = os.Getenv("HOME") + ".local/share/request.db"
)
