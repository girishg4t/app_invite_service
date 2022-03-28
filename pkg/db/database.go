package db

import (
	"database/sql"
	"fmt"

	"github.com/girishg4t/app_invite_service/pkg/model"
	"github.com/jinzhu/gorm"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"

	// side load for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// ConnectToDb Connects to the database
func ConnectToDb(dbCfg model.Database) (*gorm.DB, error) {
	dbConn, err := gorm.Open(dbCfg.Type, dbCfg.DataSource)
	if err != nil {
		return nil, err
	}
	dbConn.LogMode(dbCfg.Debug)
	dbConn.SingularTable(true)
	log, _ := zap.NewDevelopment()
	err = migrateSchema(dbCfg, log)
	return dbConn, err
}

func migrateSchema(dbCfg model.Database, logger *zap.Logger) error {
	logger.Debug("Migrating database schema", zap.String("type", dbCfg.Type), zap.String("schema", dbCfg.Schema))
	m := migrate.FileMigrationSource{Dir: dbCfg.Schema}
	db, err := sql.Open(dbCfg.Type, dbCfg.DataSource)
	if err != nil {
		return fmt.Errorf("error opening db connection: %w", err)
	}
	defer db.Close()
	n, err := migrate.Exec(db, dbCfg.Type, m, migrate.Up)
	if err != nil {
		return fmt.Errorf("error applying db migration: %w", err)
	}
	logger.Debug("db migration applied!", zap.Int("applied_migration", n))
	return nil
}
