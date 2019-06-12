package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

const driverName = "postgres"

type Config struct {
	DSN          string
	DatabaseName string
	SourceURL    string
	Logger       *log.Logger
}

type Migration struct {
	config Config
}

func NewMigration(c Config) *Migration {
	return &Migration{config: c}
}

func (m *Migration) Do() error {
	m.config.Logger.Printf("database migration started: %s", m.config.DSN)

	db, err := m.connectDB()
	if err != nil {
		return fmt.Errorf("migration failed: unable to open database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: m.config.DatabaseName,
	})
	if err != nil {
		return fmt.Errorf("migration failed: unable to create driver: %v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(m.config.SourceURL, m.config.DatabaseName, driver)
	if err != nil {
		return fmt.Errorf("migration failed: unable to create migration: %v", err)
	}

	err = migration.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("unable to migrate: %v", err)
		}
	}

	m.config.Logger.Print("database successfully migrated")

	return nil
}

func (m *Migration) connectDB() (*sql.DB, error) {
	db, err := sql.Open(driverName, m.config.DSN)
	if err != nil {
		return nil, err
	}
	m.config.Logger.Print("database opened")
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			return db, nil
		}
		m.config.Logger.Printf("database ping failed: retry %d: %v", i, err)
		time.Sleep(time.Second)
	}
	return nil, err
}
