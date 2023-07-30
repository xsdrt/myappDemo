package hiSpeed

import (
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//These funcs should be able to be run from the cli...
//write some func that will run an migrationUp ; load a file from the migrations folder in the myappDemo and execute the sql inside the file
//write a migrationDown func; reverses the migrationUP
//write a func that also forces a migration; whew!!  lets see if I can do it....
//going to use a package called golang-migrate/migrate from github designed to use with a cli...
//grab the drivers also for mysql-postgres-sqlserver ie: go get github.com/golang-migrate/migrate/v4/database/postgres replace the db name(s) with the ones you want
//and also  go get github.com/golang-migrate/migrate/v4/source/file

func (h *HiSpeed) MigrateUp(dsn string) error { //need to change the postgres driver to the one used by golang-migrate !!!
	rootPath := filepath.ToSlash(h.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn) //for windows path...hmmm but if using wsl might not need to use; leave it for now as end user might mot be using wsl
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration: ", err)
		return err
	}

	return nil
}

func (h *HiSpeed) MigrateDownAll(dsn string) error {
	rootPath := filepath.ToSlash(h.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		return err
	}

	return nil
}

func (h *HiSpeed) Steps(n int, dsn string) error { //if the int n is positive; will run up migrations; if a negative run down migrations...
	rootPath := filepath.ToSlash(h.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		return err
	}

	return nil
}

// Sometime when running migrations might have an error in the file itself, sounds like golang-migrate will mark this as dirty in the db and not allow , well force a fix.
func (h *HiSpeed) MigrateForce(dsn string) error {
	rootPath := filepath.ToSlash(h.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(-1); err != nil { //force the removal of last migration...
		return err
	}

	return nil
}

//need to write some tests for these funcs...
