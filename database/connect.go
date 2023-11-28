package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Connect() *DB {
	gDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=qooqpll password=qooqpll dbname=blockchain port=5434 sslmode=disable",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}

	db := NewDB(gDB)
	db.Migrate()

	return db
}

func (db *DB) Migrate() {
	db.AutoMigrate(ApiKeys{})
}

func NewDB(db *gorm.DB) *DB {
	return &DB{db}
}
