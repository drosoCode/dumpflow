package database

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate cp ../../sql/schema.sql ./

//go:embed schema.sql
var embedFS embed.FS

var db map[string]*pgxpool.Pool
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
	db = map[string]*pgxpool.Pool{}
}

func DeleteDB(name string) error {
	l := len(dbms.Prefix)
	if name[0:l] != dbms.Prefix {
		name = dbms.Prefix + name
	}

	db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=prefer", dbms.User, dbms.Password, dbms.Host, dbms.Port))
	if err != nil {
		return err
	}
	_, err = db.Exec(context.Background(), "DROP DATABASE "+name)
	db.Close(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func PreImport(name string) error {
	l := len(dbms.Prefix)
	if name[0:l] != dbms.Prefix {
		name = dbms.Prefix + name
	}

	db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=prefer", dbms.User, dbms.Password, dbms.Host, dbms.Port, name))
	if err != nil {
		return err
	}
	defer db.Close(context.Background())

	_, err = db.Exec(context.Background(), `
		ALTER TABLE badges SET UNLOGGED;
		ALTER TABLE comments SET UNLOGGED;
		ALTER TABLE post_history SET UNLOGGED;
		ALTER TABLE post_links SET UNLOGGED;
		ALTER TABLE posts SET UNLOGGED;
		ALTER TABLE tags SET UNLOGGED;
		ALTER TABLE users SET UNLOGGED;
		ALTER TABLE votes SET UNLOGGED;
	`)
	if err != nil {
		return err
	}
	return nil
}

func PostImport(name string) error {
	l := len(dbms.Prefix)
	if name[0:l] != dbms.Prefix {
		name = dbms.Prefix + name
	}

	db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=prefer", dbms.User, dbms.Password, dbms.Host, dbms.Port, name))
	if err != nil {
		return err
	}
	defer db.Close(context.Background())

	_, err = db.Exec(context.Background(), `
		ALTER TABLE badges SET LOGGED, ADD PRIMARY KEY (id);
		ALTER TABLE comments SET LOGGED, ADD PRIMARY KEY (id);
		ALTER TABLE post_history SET LOGGED, ADD PRIMARY KEY (id);
		ALTER TABLE post_links SET LOGGED, ADD PRIMARY KEY (id);
		ALTER TABLE posts SET LOGGED, ADD PRIMARY KEY (id);
		ALTER TABLE tags SET LOGGED, ADD PRIMARY KEY (id);
		ALTER TABLE users SET LOGGED, ADD PRIMARY KEY (id);
		ALTER TABLE votes SET LOGGED, ADD PRIMARY KEY (id);

		DROP TABLE IF EXISTS "posts_idx";
		CREATE INDEX posts_idx ON posts USING GIN ("body_index");

		DROP TABLE IF EXISTS "comments_idx";
		CREATE INDEX comments_idx ON comments USING GIN ("text_index");

		DROP TABLE IF EXISTS "post_history_idx";
		CREATE INDEX post_history_idx ON post_history USING GIN ("text_index");
	`)
	if err != nil {
		return err
	}
	return nil
}

func GetDB(name string) (*Queries, error) {
	db, err := GetRawDB(name)
	if err != nil {
		return nil, err
	}
	return New(db), nil
}

func GetRawDB(name string) (*pgxpool.Pool, error) {
	dbname := name
	l := len(dbms.Prefix)
	if name[0:l] != dbms.Prefix {
		dbname = dbms.Prefix + name
	}

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

func connect(conn DBMSConn, database string) (*pgxpool.Pool, error) {
	// connect to db
	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=prefer", conn.User, conn.Password, conn.Host, conn.Port, database))

	var rows pgx.Rows
	if err == nil {
		// try to count tables
		rows, err = db.Query(context.Background(), "SELECT COUNT(*) FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';")
		if err != nil {
			db.Close()
		}
	}

	if err != nil {
		// if query failed, the db likely does not exists, try to create it

		db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=prefer", conn.User, conn.Password, conn.Host, conn.Port))
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(context.Background(), "CREATE DATABASE "+database)
		db.Close(context.Background())
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
		_, err = db.Exec(context.Background(), string(schema))
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
