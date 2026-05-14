package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database connection
var db *gorm.DB

// Models
type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	SKU         string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	Quantity    int       `gorm:"default:0" json:"quantity"`
	UnitPrice   float64   `gorm:"not null" json:"unit_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Health check response
type HealthResponse struct {
	Status   string            `json:"status"`
	Service  string            `json:"service"`
	Database DatabaseStatus    `json:"database"`
	Time     string            `json:"time"`
}

type DatabaseStatus struct {
	Type      string `json:"type"`
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get configuration
	port := getEnv("PORT", "7000")
	mysqlHost := getEnv("MYSQL_HOST", "localhost")
	mysqlPort := getEnv("MYSQL_PORT", "3306")
	mysqlUser := getEnv("MYSQL_USER", "root")
	mysqlPassword := getEnv("MYSQL_PASSWORD", "password")
	mysqlDatabase := getEnv("MYSQL_DATABASE", "warehouse_db")

	// Try to connect to MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("⚠️  Failed to connect to MySQL: %v", err)
		log.Println("💡 Make sure MySQL is running on port 3306")
		log.Println("💡 Run: docker run -d --name mysql-warehouse -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=warehouse_db -p 3306:3306 mysql:8.0")
	} else {
		log.Println("✅ Database connected successfully!")
		
		// Auto migrate
		if err := db.AutoMigrate(&Item{}); err != nil {
			log.Printf("⚠️  Failed to migrate database: %v", err)
		} else {
			log.Println("✅ Database migration completed!")
		}
	}

	// Setup Gin router
	router := gin.Default()

	// CORS middleware — allow frontend origin
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	router.GET("/health", healthCheck)

	// Item endpoints
	router.POST("/api/items", createItem)
	router.GET("/api/items", getItems)
	router.GET("/api/items/:id", getItem)

	// Start server
	log.Printf("🚀 Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// Health check handler
func healthCheck(c *gin.Context) {
	dbStatus := DatabaseStatus{
		Type:      "mysql",
		Connected: false,
		Message:   "Not connected",
	}

	// Check database connection
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Ping(); err == nil {
				dbStatus.Connected = true
				dbStatus.Message = "MySQL connected"
			} else {
				dbStatus.Message = fmt.Sprintf("MySQL ping failed: %v", err)
			}
		} else {
			dbStatus.Message = fmt.Sprintf("Failed to get DB instance: %v", err)
		}
	} else {
		dbStatus.Message = "Database not initialized"
	}

	status := "ok"
	if !dbStatus.Connected {
		status = "degraded"
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status:   status,
		Service:  "go-warehouse-api",
		Database: dbStatus,
		Time:     time.Now().Format(time.RFC3339),
	})
}

// Create item handler
func createItem(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}

	if err := db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// Get all items handler
func getItems(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}

	var items []Item
	if err := db.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// Get single item handler
func getItem(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}

	id := c.Param("id")
	var item Item
	if err := db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
