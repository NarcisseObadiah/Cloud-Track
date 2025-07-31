package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "paas-api/auth"
    "paas-api/handlers"
    "github.com/gin-contrib/cors"

)


func main() {
    if err := middleware.InitJWT(); err != nil {
        log.Fatalf("Failed to initialize JWKS: %v", err)
    }

    r := gin.Default()

    // 👇 Add CORS configuration here
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, // your frontend URL
        AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

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
}
