package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/pkg/config"
)

// InitDB 初始化 SQLite 数据库并应用推荐参数。
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	if cfg.Path == "" {
		return nil, fmt.Errorf("database path is required")
	}

	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := configureSQLite(sqlDB); err != nil {
		return nil, err
	}

	return db, nil
}

func configureSQLite(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = -2000",
		"PRAGMA foreign_keys = ON",
	}
	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("apply pragma %s: %w", pragma, err)
		}
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)
	return nil
}

// RunMigrations 执行文件目录中的所有未应用迁移。
func RunMigrations(dbPath, migrationsDir string) error {
	if migrationsDir == "" {
		return fmt.Errorf("migrations directory is required")
	}
	if _, err := os.Stat(migrationsDir); err != nil {
		return fmt.Errorf("stat migrations dir: %w", err)
	}

	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("open sqlite db for migrate: %w", err)
	}
	defer sqlDB.Close()

	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("create sqlite driver: %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s", migrationsDir)
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}
	return nil
}

// GetMigrationVersion 返回当前版本号及是否 dirty。
func GetMigrationVersion(dbPath, migrationsDir string) (uint, bool, error) {
	if migrationsDir == "" {
		return 0, false, fmt.Errorf("migrations directory is required")
	}

	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, false, err
	}
	defer sqlDB.Close()

	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return 0, false, err
	}

	sourceURL := fmt.Sprintf("file://%s", migrationsDir)
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "sqlite3", driver)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil
		}
		return 0, false, err
	}

	return version, dirty, nil
}
