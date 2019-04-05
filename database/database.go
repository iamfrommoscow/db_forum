package database

import (
	"io/ioutil"
	"log"

	"github.com/jackc/pgx"
)

var connectionConfig = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "postgres",
	User:     "postgres",
	Password: "postgres",
}

var connectionPoolConfig = pgx.ConnPoolConfig{
	ConnConfig:     connectionConfig,
	MaxConnections: 8,
}

func Connect() *pgx.ConnPool {

	connectionPool, err := pgx.NewConnPool(connectionPoolConfig)

	if err != nil {
		log.Fatal(err)
	}
	if query, err := ioutil.ReadFile("createDB.sql"); err != nil {
		log.Println(err)
		return nil
	} else {
		if _, err := connectionPool.Exec(string(query)); err != nil {
			log.Println(err)
			return nil
		}
	}
	return connectionPool
}

var Connection = Connect()

func StartTransaction() *pgx.Tx {
	if transaction, err := Connection.Begin(); err != nil {
		log.Fatal(err)

		return transaction
	} else {

		return transaction
	}
}

func Exec(queryStr string) error {
	_, err := Connection.Exec(queryStr)
	return err

}
