package hiSpeed

import "database/sql"

type initPaths struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}

//Below referenced in hiSpeed.go file...
type databaseConfig struct { //Not exporting of course
	dsn      string //The connection string...
	database string //The database connecting too...
}

type Database struct { //Exporting this one as it would be useful to the end users of HiSpeed...
	DataType string  //What type of Db : Postgres:MSSQL, Mariadb, MySql ect...
	Pool     *sql.DB //Connection pool from drive.go with a pointer too sql.DB
}
