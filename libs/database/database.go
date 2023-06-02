package database

import (
	//"database/sql"
	//"database/sql/driver"
	"errors"
	"fmt"
	"simplerest/libs/settings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	_ "github.com/sijms/go-ora/v2"
)

const (
	DBDriverSQLite3 = "sqlite3"
	DBDRiverPSQL    = "postres"
	DBDRiverPGX     = "pgx"
	DBDRiverOracle  = "oracle"
	DBDRiverOCI8    = "oci"
	DBDRiverMysql   = "mysql"
)

type Database struct {
	dbsettings settings.Database
	DB         *sqlx.DB
}

func New(s settings.Database) Database {
	return Database{
		dbsettings: s,
		DB:         nil,
	}
}

func goodToGo(driver string) bool {
	switch driver {
	case DBDRiverMysql, DBDRiverOCI8, DBDRiverOracle, DBDRiverPGX, DBDRiverPSQL, DBDriverSQLite3:
		return true
	default:
		return false
	}
}

func (d *Database) Initialize() error {
	var err error

	if !goodToGo(d.dbsettings.Driver) {
		return errors.New("Unknown driver " + d.dbsettings.Driver)
	}

	if d.DB, err = sqlx.Open(d.dbsettings.Driver, d.dbsettings.Location); err != nil {
		return err
	}

	if err = d.DB.Ping(); err != nil {
		return err
	}

	if err = d.migrate(); err != nil {
		return err
	}
	return nil
}

func (d *Database) dialect() string {
	switch d.dbsettings.Driver {
	case "pgx":
		return "postgres"
	case "oracle":
		return "oci8"
	default:
		return d.dbsettings.Driver
	}
}

func (d *Database) migrate() error {
	fmt.Printf("Applying migrations: ")
	if d.dbsettings.Migrations != "" {
		migrations := &migrate.FileMigrationSource{
			Dir: d.dbsettings.Migrations,
		}
		n, err := migrate.Exec(d.DB.DB, d.dialect(), migrations, migrate.Up)
		if err != nil {
			fmt.Println("error", n)
			return err
		}
		fmt.Printf("%d done\n", n)
	} else {
		fmt.Println("none set up")
	}
	return nil
}
