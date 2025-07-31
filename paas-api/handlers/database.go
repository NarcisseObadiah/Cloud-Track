package handlers

import (
	"net/http"
	"paas-api/k8s"

	"github.com/gin-gonic/gin"
)

type DBRequest struct {
	Username string `json:"username" binding:"required"`
	DBName   string `json:"db_name"  binding:"required"`
	Password string `json:"password" `
}

type DeleteDBRequest struct {
	Username string `json:"username" binding:"required"`
	DBName   string `json:"db_name"  binding:"required"`
}


// ListAllTenantPodsHandler returns pods grouped by tenant namespace
func ListAllTenantPodsHandler(c *gin.Context) {
	podGroups, err := k8s.ListAllTenantPods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, podGroups)
}


func ListTenantPodsHandler(c *gin.Context) {
	namespace := c.Param("namespace")

	pods, err := k8s.ListTenantPodsJSON(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pods)
}

func DeleteDatabase(c *gin.Context) {
	var req DeleteDBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespace := "tenant-" + req.Username
	err := k8s.DeleteTenantDB(namespace, req.DBName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Database deleted",
		"namespace": namespace,
		"db_name":   req.DBName,
	})
}



func GetDatabaseStatus(c *gin.Context) {
	username := c.Param("username")
	dbName := c.Param("db_name")
	namespace := "tenant-" + username

	status, err := k8s.CheckTenantDBStatus(namespace, dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"db_name":  dbName,
		"status":   status,
		"username": username,
	})
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

