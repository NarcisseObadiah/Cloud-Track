package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"paas-api/k8s"
)

type DBRequest struct {
	Username string `json:"username" binding:"required"`
	DBName   string `json:"db_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateDatabase(c *gin.Context) {
	var req DBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespace := "tenant-" + req.Username

	err := k8s.ProvisionTenantDB(namespace, req.DBName, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Database provisioned",
		"namespace": namespace,
		"db_name":   req.DBName,
	})
}
