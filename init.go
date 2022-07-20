package zpmysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func init() {
	sql.Register("zpmysql", &ZpMySQLDriver{driver: &mysql.MySQLDriver{}})
}
