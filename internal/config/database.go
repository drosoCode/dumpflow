package config

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/drosocode/dumpflow/internal/database"
)

//go:generate cp ../../sql/config.sql ./

//go:embed config.sql
var embedFS embed.FS

var DB *Queries
var dbms database.DBMSConn

func ConfigDB(conn database.DBMSConn) error {
	dbms = conn
	var err error
	DB, err = connect(conn)
	return err
}

func connect(conn database.DBMSConn) (*Queries, error) {
	// connect to db
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=dumpflow sslmode=disable", conn.Host, conn.Port, conn.User, conn.Password))
	if err != nil {
		return nil, err
	}

	// try to count tables
	rows, err := db.Query("SELECT COUNT(*) FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';")
	if err != nil {
		// if query failed, the db likely does not exists, try to create it
		db.Close()

		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", conn.Host, conn.Port, conn.User, conn.Password))
		if err != nil {
			return nil, err
		}
		_, err = db.Exec("CREATE DATABASE dumpflow")
		db.Close()
		if err != nil {
			return nil, err
		}
		return connect(conn)
	}
	// read the result of the query
	var cnt int
	rows.Next()
	err = rows.Scan(&cnt)
	rows.Close()
	if err != nil {
		return nil, err
	}

	// read the schema file
	schema, err := embedFS.ReadFile("config.sql")
	if err != nil {
		return nil, err
	}

	// if the count doesn't match, execute the table creation script
	if cnt != strings.Count(string(schema), "CREATE TABLE") {
		_, err = db.Exec(string(schema))
		if err != nil {
			return nil, err
		}
	}

	return New(db), nil
}
