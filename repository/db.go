package repository

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

type RepositoryFactory struct {
	dialect     Dialect
	autoMigrate bool
}

type Dialect uint

const (
	DialectSQLite = iota
	DialectMySQL
	DialectPostgres
	DialectSQLServer
)

func NewRepositoryFactory(dialect Dialect, autoMigrate bool) *RepositoryFactory {
	return &RepositoryFactory{
		dialect:     dialect,
		autoMigrate: autoMigrate,
	}
}

func (r RepositoryFactory) connect() *gorm.DB {
	switch r.dialect {
	case DialectSQLite:
		return r.connectSQLite()
	default:
		message := "Unsupported dialect"
		log.Fatal(message)
		return nil
	}
}

func (r RepositoryFactory) connectSQLite() *gorm.DB {
	// Open an in-memory SQLite database
	once.Do(func() {
		_db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			message := "Failed to open SQLite database"
			log.Fatal(message)
		}
		db = _db
	})
	return db
}
