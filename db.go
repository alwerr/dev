package dev

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type DB struct {
	Sql *sql.DB
}

func StartDb(path, structure string) DB {
	s := strings.Split(path, `db.db`)[0]
	// if set.Dev {
	// 	os.RemoveAll(s)
	// }
	if _, err := os.Stat(s); os.IsNotExist(err) {
		os.MkdirAll(s, 0700)
	}
	// Best method for concurrent write
	db, err := sql.Open(`sqlite3`, path+`?cache=shared&_journal_mode=WAL&pooling=true&_busy_timeout=5000`)
	// this must have!but it not work when there is more then  one connection
	/*
			db := dbs.Start(`data/db.db`, DBStructure())
		dba := dbs.Start(`dataa/db.db`, DBStructure())
		dba just return timeout
	*/
	db.SetMaxOpenConns(1)

	if err != nil {
		panic(err)
	}
	// Create new tables
	_, err = db.Exec(structure)
	if err != nil {
		panic(err)
	}

	//
	fmt.Println(`Database connected`)
	return DB{Sql: db}
}

func (db DB) FetchRows(q string, args ...any) []byte {
	return FetchRows(db.Sql, q, args...)
}
func (db DB) SendRows(w io.Writer, q string, args ...any) {
	io.Writer.Write(w, FetchRows(db.Sql, q, args...))
}
func (db DB) FetchRow(q string, args ...any) []byte {
	return FetchRow(db.Sql, q, args...)
}
func (db DB) SendRow(w io.Writer, q string, args ...any) {
	io.Writer.Write(w, FetchRow(db.Sql, q, args...))
}
func (db DB) FetchRowMap(q string, args ...any) map[string]interface{} {
	return aaa(db.Sql, q, args...)
}

func (db DB) Exec(q string, args ...any) (string, error) {
	res, e := db.Sql.Exec(q, args...)
	if e != nil {
		fmt.Print(`Sql error: ` + e.Error() + q)
		return "", e
	}
	jsn, e := json.Marshal(args)
	if e != nil {
		fmt.Println(e)
	}
	_, e = db.Sql.Exec(`insert into log(q,f,t)values(?,?,?)`, q, string(jsn), time.Now().Unix())
	if e != nil {
		fmt.Println(e)
	}
	var id int64
	id, e = res.LastInsertId()
	return strconv.FormatInt(id, 10), nil
}

func FetchRows(dbs *sql.DB, q string, args ...any) []byte {
	rows, _ := dbs.Query(q, args...)
	cols, _ := rows.Columns()
	var data []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return []byte(``)
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		data = append(data, m)

	}
	bolB, _ := json.Marshal(data)
	if len(data) == 0 {
		bolB = []byte(`[]`)
	}
	fmt.Print()
	return bolB

}
func FetchRow(dbs *sql.DB, q string, args ...any) []byte {
	rows, _ := dbs.Query(q, args...) // Note: Ignoring errors for brevity
	cols, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		bolB, _ := json.Marshal(m)
		return bolB
	}
	return nil
}
func FetchRowMap(dbs *sql.DB, q string, args ...any) map[string]interface{} {
	rows, _ := dbs.Query(q, args...) // Note: Ignoring errors for brevity
	cols, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		json.Marshal(m)
	}
	return nil
}
func aaa(db *sql.DB, q string, args ...any) map[string]interface{} {
	rows, err := db.Query(q, args...)
	if err != nil {
		// Handle error
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		// Handle error
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range values {
		valuePtrs[i] = &values[i]
	}
	var m map[string]interface{}
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			// Handle error
		}

		m = make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Handle null values if needed
			m[col] = val
		}

		// Use the map
		// fmt.Println(m)
	}
	// no results = m is nil
	// fmt.Println(m == nil)
	return m

	// return nil
}

/*
func FetchRows(q string) []byte {
	rows, _ := dbs.Query(q)
	cols, _ := rows.Columns()
	var data []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return []byte(``)
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		data = append(data, m)

	}
	bolB, _ := json.Marshal(data)
	return bolB

}
func FetchRow(q string) []byte {
	rows, _ := dbs.Query(q) // Note: Ignoring errors for brevity
	cols, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		bolB, _ := json.Marshal(m)
		return bolB
	}
	return nil
}
*/
