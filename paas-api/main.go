package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"paas-api/handlers"
	"paas-api/auth"
)

func main() {
	if err := middleware.InitJWT(); err != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	r := gin.Default()

	// Public / tenant routes
	r.POST("/databases", middleware.AuthMiddleware("tenant"), handlers.CreateDatabase)
	r.DELETE("/databases", middleware.AuthMiddleware("tenant"), handlers.DeleteDatabase)
	r.GET("/databases/:username/:db_name/status", middleware.AuthMiddleware("tenant"), handlers.GetDatabaseStatus)
	r.GET("/pods/:namespace", middleware.AuthMiddleware("tenant"), handlers.ListTenantPodsHandler)

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware("admin"))
	{
		admin.GET("/tenants/pods", handlers.ListAllTenantPodsHandler)
	}

	log.Println("API listening on port 8080")
	r.Run(":8080")
}¬