package cmd

import (
	"flag"

	"github.com/drosocode/dumpflow/internal/database"
	"github.com/jamiealquiza/envy"
)

type ConfigData struct {
	DeleteOnFinish bool
	ServeAddr      string
	DataDir        string
	DBConn         database.DBMSConn
}

var Config ConfigData

func ParseConfig() {
	serve := flag.String("serve", "0.0.0.0:3002", "bind address")
	host := flag.String("db_host", "127.0.0.1", "host of the DBMS server")
	port := flag.Int("db_port", 5432, "port of the DBMS server")
	username := flag.String("db_username", "postgres", "username for the DBMS server")
	password := flag.String("db_password", "", "password for the DBMS server")
	prefix := flag.String("prefix", "so_", "prefix for stackexchange databases")
	deleteOnFinish := flag.Bool("delete_onfinish", false, "delete the files when an import is finished")
	datadir := flag.String("data_dir", ".", "path to data dir")

	envy.Parse("DF")
	flag.Parse()

	Config = ConfigData{
		DeleteOnFinish: *deleteOnFinish,
		ServeAddr:      *serve,
		DataDir:        *datadir,
		DBConn: database.DBMSConn{
			Host:     *host,
			Port:     *port,
			User:     *username,
			Password: *password,
			Prefix:   *prefix,
		},
	}
}
