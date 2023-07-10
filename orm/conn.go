package orm

import (
	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open returns a gorm.DB instance for the given sql.DB sqlite instance.
func Open(sqlDB *sql.DB) (*gorm.DB, error) {
	return gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
}
