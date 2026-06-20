package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func handleHealth(c *gin.Context) {
	ok(c, gin.H{"service": "qingtan-library", "status": "ok"})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func main() {
	dbPath := os.Getenv("LIBRARY_DB")
	if dbPath == "" {
		dbPath = "library.db"
	}
	initDB(dbPath)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(corsMiddleware())

	api := r.Group("/api")
	api.GET("/health", handleHealth)
	api.POST("/auth/login", handleLogin)

	auth := api.Group("")
	auth.Use(AuthRequired())
	auth.GET("/auth/me", handleMe)
	auth.PUT("/auth/password", handleChangePassword)

	auth.GET("/books", handleListBooks)
	auth.GET("/books/categories", handleCategories)
	auth.GET("/books/:id", handleGetBook)

	admin := auth.Group("")
	admin.Use(AdminRequired())
	admin.POST("/books", handleCreateBook)
	admin.PUT("/books/:id", handleUpdateBook)
	admin.DELETE("/books/:id", handleDeleteBook)

	auth.GET("/borrows", handleListBorrows)
	auth.POST("/borrows", handleBorrow)
	auth.POST("/borrows/:id/return", handleReturn)

	admin.GET("/users", handleListUsers)
	admin.POST("/users", handleCreateUser)
	admin.PUT("/users/:id", handleUpdateUser)
	admin.DELETE("/users/:id", handleDeleteUser)

	admin.GET("/stats", handleStats)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("青檀图书馆服务启动 :%s (db=%s)", port, dbPath)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
