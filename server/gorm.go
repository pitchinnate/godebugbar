package godebugbar

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	callbackPrefix   = "debugbar"
	startTimeKey     = "debugbar:start_time"
	queryIDKey       = "debugbar:query_id"
)

// GormDebugBarPlugin is the GORM plugin for tracking database queries
type GormDebugBarPlugin struct {
	debugBar *DebugBar
}

// Name returns the plugin name
func (p *GormDebugBarPlugin) Name() string {
	return "GormDebugBarPlugin"
}

// Initialize sets up the GORM callbacks
func (p *GormDebugBarPlugin) Initialize(db *gorm.DB) error {
	// Register callbacks for all operations
	callbacks := []struct {
		name     string
		register func(*gorm.DB, string, func(*gorm.DB)) error
	}{
		{"create", registerBeforeCreate},
		{"query", registerBeforeQuery},
		{"update", registerBeforeUpdate},
		{"delete", registerBeforeDelete},
		{"row", registerBeforeRow},
		{"raw", registerBeforeRaw},
	}

	for _, cb := range callbacks {
		beforeName := fmt.Sprintf("%s:before_%s", callbackPrefix, cb.name)
		if err := cb.register(db, beforeName, p.beforeCallback); err != nil {
			return err
		}
	}

	// Register after callbacks
	afterCallbacks := []struct {
		name     string
		register func(*gorm.DB, string, func(*gorm.DB)) error
	}{
		{"create", registerAfterCreate},
		{"query", registerAfterQuery},
		{"update", registerAfterUpdate},
		{"delete", registerAfterDelete},
		{"row", registerAfterRow},
		{"raw", registerAfterRaw},
	}

	for _, cb := range afterCallbacks {
		afterName := fmt.Sprintf("%s:after_%s", callbackPrefix, cb.name)
		if err := cb.register(db, afterName, p.afterCallback); err != nil {
			return err
		}
	}

	return nil
}

func registerBeforeCreate(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Create().Before("gorm:create").Register(name, fn)
}

func registerAfterCreate(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Create().After("gorm:create").Register(name, fn)
}

func registerBeforeQuery(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Query().Before("gorm:query").Register(name, fn)
}

func registerAfterQuery(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Query().After("gorm:query").Register(name, fn)
}

func registerBeforeUpdate(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Update().Before("gorm:update").Register(name, fn)
}

func registerAfterUpdate(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Update().After("gorm:update").Register(name, fn)
}

func registerBeforeDelete(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Delete().Before("gorm:delete").Register(name, fn)
}

func registerAfterDelete(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Delete().After("gorm:delete").Register(name, fn)
}

func registerBeforeRow(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Row().Before("gorm:row").Register(name, fn)
}

func registerAfterRow(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Row().After("gorm:row").Register(name, fn)
}

func registerBeforeRaw(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Raw().Before("gorm:raw").Register(name, fn)
}

func registerAfterRaw(db *gorm.DB, name string, fn func(*gorm.DB)) error {
	return db.Callback().Raw().After("gorm:raw").Register(name, fn)
}

// beforeCallback is called before each database operation
func (p *GormDebugBarPlugin) beforeCallback(db *gorm.DB) {
	if !p.debugBar.config.Enabled {
		return
	}

	db.InstanceSet(startTimeKey, time.Now())
	db.InstanceSet(queryIDKey, uuid.New().String())
}

// afterCallback is called after each database operation
func (p *GormDebugBarPlugin) afterCallback(db *gorm.DB) {
	if !p.debugBar.config.Enabled {
		return
	}

	// Get start time
	startTimeVal, ok := db.InstanceGet(startTimeKey)
	if !ok {
		return
	}
	startTime, ok := startTimeVal.(time.Time)
	if !ok {
		return
	}

	// Get query ID
	queryIDVal, ok := db.InstanceGet(queryIDKey)
	if !ok {
		return
	}
	queryID, ok := queryIDVal.(string)
	if !ok {
		return
	}

	duration := time.Since(startTime)

	// Build query info
	queryInfo := QueryInfo{
		ID:           queryID,
		Query:        db.Statement.SQL.String(),
		Args:         db.Statement.Vars,
		Duration:     duration,
		DurationMs:   float64(duration.Nanoseconds()) / 1e6,
		RowsAffected: db.RowsAffected,
		StartTime:    startTime,
		Source:       getCallerInfo(),
	}

	// Capture error if any
	if db.Error != nil {
		queryInfo.Error = db.Error.Error()
	}

	// Add query to the request context
	p.debugBar.addQuery(db.Statement.Context, queryInfo)
}

// getCallerInfo returns the file and line number of the caller
func getCallerInfo() string {
	// Skip frames to get to the actual caller
	// Skip: getCallerInfo, afterCallback, gorm internals
	for i := 4; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// Skip gorm internal files
		if strings.Contains(file, "gorm.io") {
			continue
		}
		// Skip this package
		if strings.Contains(file, "godebugbar") {
			continue
		}

		return fmt.Sprintf("%s:%d", file, line)
	}

	return ""
}

