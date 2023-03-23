package hiSpeed

import (
	"database/sql"

	_ "github.com/jackc/pgconn" //Import the libraries needed...after in the terminal using go get each one...
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func (h *HiSpeed) OpenDB(dbType, dsn string) (*sql.DB, error) {
	if dbType == "postgres" || dbType == "postgresql" {
		dbType = "pgx"
	}

	db, err := sql.Open(dbType, dsn) //using go sql driver get a db connection...
	if err != nil {
		return nil, err
	}

	err = db.Ping() //Then ping the database to ensure connection
	if err != nil {
		return nil, err
	}

	return db, nil //If all above works return pool of connections
}
