package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"
)

//go:generate cp ../../sql/schema.sql ./

//go:embed schema.sql
var embedFS embed.FS

var db map[string]*Queries
var dbms DBMSConn

type DBMSConn struct {
	Host     string
	Port     int
	User     string
	Password string
	Prefix   string
}

func ConfigDB(conn DBMSConn) {
	dbms = conn
	db = map[string]*Queries{}
}

func DeleteDB(name string) error {
	l := len(dbms.Prefix)
	if name[0:l] != dbms.Prefix {
		return errors.New("invalid prefix")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", dbms.Host, dbms.Port, dbms.User, dbms.Password))
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP DATABASE " + name)
	db.Close()
	if err != nil {
		return err
	}
	return nil
}

func GetDB(name string) (*Queries, error) {
	dbname := dbms.Prefix + name

	if dbconn, ok := db[dbname]; ok {
		return dbconn, nil
	}

	d, err := connect(dbms, dbname)
	if err != nil {
		return nil, err
	}
	db[dbname] = d
	return d, nil
}

func connect(conn DBMSConn, database string) (*Queries, error) {
	// connect to db
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conn.Host, conn.Port, conn.User, conn.Password, database))
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
		_, err = db.Exec("CREATE DATABASE " + database)
		db.Close()
		if err != nil {
			return nil, err
		}
		return connect(conn, database)
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
	schema, err := embedFS.ReadFile("schema.sql")
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
