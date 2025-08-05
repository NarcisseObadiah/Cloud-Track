package handlers

import (
	"fmt"
	"time"
	"net/http"
	"paas-api/k8s"

	"github.com/gin-gonic/gin"
)

type DBRequest struct {
	Username string `json:"username" binding:"required"`
	DBName   string `json:"db_name"`   // Optional: will auto-generate if not provided
	Replicas int    `json:"replicas"`  // Optional: defaults to 1
}

type DeleteDBRequest struct {
	Username string `json:"username" binding:"required"`
	DBName   string `json:"db_name"  binding:"required"`
}

func ListAllTenantPodsHandler(c *gin.Context) {
	podGroups, err := k8s.ListAllTenantPods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, podGroups)
}

func ListDatabaseClusters(c *gin.Context) {
	username := c.Param("username")
	namespace := "tenant-" + username

	clusters, err := k8s.ListTenantDatabaseClusters(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"namespace": namespace,
		"clusters": clusters,
		"total_clusters": len(clusters),
		"summary": generateClusterSummary(clusters),
	})
}

func GetDatabaseClusterDetails(c *gin.Context) {
	username := c.Param("username")
	dbName := c.Param("db_name")
	namespace := "tenant-" + username

	cluster, err := k8s.GetDatabaseClusterInfo(namespace, dbName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"cluster": cluster,
	})
}

func generateClusterSummary(clusters []k8s.DatabaseClusterInfo) map[string]interface{} {
	summary := map[string]interface{}{
		"total": len(clusters),
		"ready": 0,
		"creating": 0,
		"failed": 0,
		"credentials_ready": 0,
		"connection_ready": 0,
		"manual_created": 0,
		"zalando_created": 0,
	}

	for _, cluster := range clusters {
		switch cluster.Status {
		case "Ready":
			summary["ready"] = summary["ready"].(int) + 1
		case "Creating", "Pending", "Pod Running", "Credentials Ready":
			summary["creating"] = summary["creating"].(int) + 1
		case "Failed", "Error":
			summary["failed"] = summary["failed"].(int) + 1
		}

		if cluster.CredentialsReady {
			summary["credentials_ready"] = summary["credentials_ready"].(int) + 1
		}

		if cluster.ConnectionReady {
			summary["connection_ready"] = summary["connection_ready"].(int) + 1
		}

		if cluster.CreationMethod == "manual" {
			summary["manual_created"] = summary["manual_created"].(int) + 1
		} else if cluster.CreationMethod == "zalando" {
			summary["zalando_created"] = summary["zalando_created"].(int) + 1
		}
	}

	return summary
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
		fmt.Printf("Invalid delete request: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Delete request received - Username: %s, DBName: %s\n", req.Username, req.DBName)

	namespace := "tenant-" + req.Username
	err := k8s.DeleteTenantDB(namespace, req.DBName)
	if err != nil {
		fmt.Printf("Failed to delete database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Database %s deleted successfully from namespace %s\n", req.DBName, namespace)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Database deleted successfully",
		"namespace": namespace,
		"db_name":   req.DBName,
	})
}

func GetDatabaseCredentials(c *gin.Context) {
	username := c.Param("username")
	dbName := c.Param("db_name")
	namespace := "tenant-" + username

	// Try to get credentials without waiting
	credentials, err := k8s.GetDatabaseCredentials(namespace, dbName, 5*time.Second)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Credentials not yet available",
			"message": "Database may still be initializing",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"db_name":     dbName,
		"username":    username,
		"credentials": credentials,
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

	// Auto-generate database name if not provided
	if req.DBName == "" {
		req.DBName = fmt.Sprintf("%s-db", req.Username)
	}

	// Set default replicas
	if req.Replicas <= 0 {
		req.Replicas = 1
	}

	namespace := "tenant-" + req.Username

	// Use the simpler provision function that handles everything
	credentials, err := k8s.ProvisionTenantDBWithCredentials(namespace, req.DBName, req.Replicas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Database provisioned successfully",
		"namespace":   namespace,
		"db_name":     req.DBName,
		"credentials": credentials,
	})
}
