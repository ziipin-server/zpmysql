package zpmysql

import (
	"database/sql"
	"os"
	"runtime"
	"strings"
	"testing"
)

func ensure(err error) {
	if err != nil {
		panic(err)
	}
}
func TestQueryError(t *testing.T) {
	AddListener(func(err error) {
		t.Logf("error! %s", err)
	})
	db, err := sql.Open("zpmysql", os.Getenv("TEST_DSN"))
	if err != nil {
		panic(err)
	}
	db.Query("select * from t10")
	db.Query("select x, y, z from t1")
}

func TestAtomError(t *testing.T) {
	AddListener(func(err error) {
		var stackBs [4096]byte
		n := runtime.Stack(stackBs[:], false)
		t.Logf("error! %s", err)
		for _, line := range strings.Split(string(stackBs[:n]), "\n") {
			line = strings.Trim(line, " \t")
			if strings.HasPrefix(line, "/usr/local/go") || strings.Contains(line, "database/sql.") {
				continue
			}
			t.Log(">>", line)
		}
	})
	db, err := sql.Open("zpmysql", os.Getenv("TEST_DSN"))
	ensure(err)
	tx, err := db.Begin()
	ensure(err)
	tx.Exec("SAVEPOINT SP_1")
}
