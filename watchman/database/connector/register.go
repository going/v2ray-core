package connector

import (
	"github.com/jmoiron/sqlx"
)

var Connectors = struct {
	Default ConnectorName
}{
	Default: "default",
}

type ConnectorName string

var cacheDBs = make(map[ConnectorName]*sqlx.DB, 0)

func RegisterConnector(name ConnectorName) Connector {
	conn := &Connect{name: name}
	conn.getDB = func() *sqlx.DB {
		return cacheDBs[name]
	}
	return conn
}

func RegisterDB(name ConnectorName, db *sqlx.DB) {
	if cacheDBs[name] != nil {
		// TODO
		// panic("this api not support override db connection")
	}
	cacheDBs[name] = db
}
