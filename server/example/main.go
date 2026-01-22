package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	godebugbar "github.com/user/godebugbar"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User is an example model
type User struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"size:255"`
	Email string `gorm:"size:255;uniqueIndex"`
}

func main() {
	// Initialize the debug bar
	debugBar := godebugbar.New(godebugbar.Config{
		Enabled:            true,
		WebSocketPath:      "/_debugbar/ws",
		MaxRequests:        100,
		CaptureRequestBody: true,
		MaxBodySize:        64 * 1024,
		AllowedOrigins:     []string{"*"},
	})

	// Initialize GORM with SQLite
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Register the debug bar GORM plugin
	if err := db.Use(debugBar.GormPlugin()); err != nil {
		log.Fatal("Failed to register debug bar GORM plugin:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&User{})

	// Initialize Gin
	r := gin.Default()

	// Add debug bar middleware
	r.Use(debugBar.Middleware())

	// Optional: Add recovery middleware that logs panics to debug bar
	r.Use(debugBar.RecoveryMiddleware())

	// Optional: Add error handler to capture Gin errors
	r.Use(debugBar.ErrorHandler())

	// Register debug bar WebSocket endpoint
	debugBar.RegisterRoutes(r)

	// Example routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/users", func(c *gin.Context) {
		var users []User
		// Use the request context for GORM queries
		// This associates queries with the current request in debug bar
		result := db.WithContext(c.Request.Context()).Find(&users)
		if result.Error != nil {
			debugBar.LogError(c, result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			debugBar.LogWarning(c, "Invalid user data: "+err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := db.WithContext(c.Request.Context()).Create(&user)
		if result.Error != nil {
			debugBar.LogError(c, result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Add custom data to the debug bar
		debugBar.AddCustomData(c, "created_user_id", user.ID)

		c.JSON(http.StatusCreated, user)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")

		result := db.WithContext(c.Request.Context()).First(&user, id)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				debugBar.LogNotice(c, "User not found: "+id)
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				debugBar.LogError(c, result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			}
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.GET("/error", func(c *gin.Context) {
		// Example of logging an error
		err := errors.New("this is a test error")
		debugBar.LogErrorWithContext(c, err, map[string]any{
			"custom_field": "custom_value",
			"user_id":      123,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	})

	r.GET("/panic", func(c *gin.Context) {
		// This will be caught by the recovery middleware
		panic("intentional panic for testing")
	})

	log.Println("Server starting on :8080")
	log.Println("Debug bar WebSocket available at ws://localhost:8080/_debugbar/ws")
	r.Run(":8080")
}
