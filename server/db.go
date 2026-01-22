package godebugbar

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const dbContextKey = "debugbar_db"

// DBMiddleware creates a middleware that injects a context-aware DB into each request.
// This allows handlers to use GetDB(c) instead of db.WithContext(c.Request.Context())
func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a DB instance with the request context
		contextDB := db.WithContext(c.Request.Context())
		c.Set(dbContextKey, contextDB)
		c.Next()
	}
}

// GetDB retrieves the context-aware DB from the Gin context.
// Must be used with DBMiddleware.
func GetDB(c *gin.Context) *gorm.DB {
	if db, exists := c.Get(dbContextKey); exists {
		return db.(*gorm.DB)
	}
	return nil
}
