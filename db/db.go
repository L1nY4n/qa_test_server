package db

import (
	"context"
	"fmt"
	"time"

	"qa_test_server/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dsn string, autoMigrate bool) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("connect db failed: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("open db pool failed: %w", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping db failed: %w", err)
	}

	if autoMigrate {
		if err := DB.AutoMigrate(model.Nano_Dev_capture_packed{}); err != nil {
			return fmt.Errorf("auto migrate failed: %w", err)
		}
	}
	return nil
}

func Health(timeout time.Duration) (bool, string) {
	if DB == nil {
		return false, "db is not initialized"
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return false, err.Error()
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return false, err.Error()
	}
	return true, "ok"
}
