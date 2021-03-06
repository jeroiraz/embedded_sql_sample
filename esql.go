package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codenotary/immudb/embedded/store"
	"github.com/codenotary/immudb/pkg/sql"
	"github.com/olekukonko/tablewriter"
)

func main() {
	catalogStore, err := store.Open("catalog", store.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	dataStore, err := store.Open("sqldata", store.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	engine, err := sql.NewEngine(catalogStore, dataStore, []byte("sql"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = engine.ExecStmt("CREATE DATABASE db1")
	if err != nil {
		log.Fatal(err)
	}

	_, err = engine.ExecStmt("USE DATABASE db1")
	if err != nil {
		log.Fatal(err)
	}

	_, err = engine.ExecStmt("CREATE TABLE table1 (id INTEGER, title STRING, PRIMARY KEY id)")
	if err != nil {
		log.Fatal(err)
	}

	rowCount := 10

	for i := 0; i < rowCount; i++ {
		_, err = engine.ExecStmt(fmt.Sprintf("UPSERT INTO table1 (id, title) VALUES (%d, 'title%d')", i, i))
		if err != nil {
			log.Fatal(err)
		}

	}

	time.Sleep(time.Duration(rowCount) * time.Millisecond)

	r, err := engine.QueryStmt("SELECT id, title FROM table1")
	if err != nil {
		log.Fatal(err)
	}

	consoleTable := tablewriter.NewWriter(os.Stdout)
	consoleTable.SetHeader([]string{"Id", "Title"})

	for i := 0; i < rowCount; i++ {
		row, err := r.Read()
		if err != nil {
			log.Fatal(err)
		}

		if uint64(i) != row.Values[0] {
			log.Fatalf("expected %d, actual %d", uint64(i), row.Values[0])
		}

		if fmt.Sprintf("title%d", i) != row.Values[1] {
			log.Fatalf("expected %s, actual %s", fmt.Sprintf("title%d", i), row.Values[1])
		}

		consoleRow := make([]string, 0)
		for _, consoleCol := range row.Values {
			consoleRow = append(consoleRow, fmt.Sprintf("%v", consoleCol))
		}
		consoleTable.Append(consoleRow)

	}
	consoleTable.Render()
}
