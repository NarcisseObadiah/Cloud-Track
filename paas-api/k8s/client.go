package k8s

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
type TemplateData struct {
	Namespace string
	DBName    string
	Team      string
	Replicas  int
}


type PodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Age       string `json:"age"`
	Node      string `json:"node"`
	Restarts  int32  `json:"restarts"`
	CPU       string `json:"cpu"`
	Memory    string `json:"memory"`
}

func ListTenantDatabaseClusters(namespace string) ([]DatabaseClusterInfo, error) {
	// Get all PostgreSQL CRDs in the namespace
	cmd := exec.Command("kubectl", "get", "postgresql", "-n", namespace, "-o", "jsonpath={.items[*].metadata.name}")
	output, err := cmd.Output()
	
	if err != nil {
		return nil, fmt.Errorf("failed to list PostgreSQL clusters: %w", err)
	}

	dbNames := strings.Fields(strings.TrimSpace(string(output)))
	var clusters []DatabaseClusterInfo

	for _, dbName := range dbNames {
		cluster, err := GetDatabaseClusterInfo(namespace, dbName)
		if err != nil {
			// If we can't get info for one cluster, still include it with error status
			clusters = append(clusters, DatabaseClusterInfo{
				Name:           dbName,
				Namespace:      namespace,
				Status:         "Error",
				DetailedStatus: fmt.Sprintf("Failed to get info: %v", err),
			})
		} else {
			clusters = append(clusters, *cluster)
		}
	}

	return clusters, nil
}

func ListAllTenantPods() ([]PodInfo, error) {
	clientset, err := getKubeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client: %w", err)
	}

	nsList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var pods []PodInfo

	for _, ns := range nsList.Items {
		if !isTenantNamespace(ns.Name) {
			continue
		}

		podList, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list pods in namespace %s: %w", ns.Name, err)
		}

		for _, pod := range podList.Items {
			age := time.Since(pod.CreationTimestamp.Time).Round(time.Second).String()

			var restarts int32
			for _, cs := range pod.Status.ContainerStatuses {
				restarts += cs.RestartCount
			}

			cpu := "N/A"
			mem := "N/A"
			if len(pod.Spec.Containers) > 0 {
				req := pod.Spec.Containers[0].Resources.Requests
				if cpuQty, ok := req[corev1.ResourceCPU]; ok {
					cpu = cpuQty.String()
				}
				if memQty, ok := req[corev1.ResourceMemory]; ok {
					mem = memQty.String()
				}
			}

			pods = append(pods, PodInfo{
				Name:      pod.Name,
				Namespace: ns.Name,
				Status:    string(pod.Status.Phase),
				Age:       age,
				Node:      pod.Spec.NodeName,
				Restarts:  restarts,
				CPU:       cpu,
				Memory:    mem,
			})
		}
	}

	return pods, nil
}


// Helper to identify tenant namespaces (adjust prefix as needed)
func isTenantNamespace(ns string) bool {
	return len(ns) > 7 && ns[:7] == "tenant-"
}

func ListTenantPodsJSON(namespace string) ([]PodInfo, error) {
	clientset, err := getKubeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client: %w", err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	var result []PodInfo
	for _, pod := range pods.Items {
		age := time.Since(pod.CreationTimestamp.Time).Round(time.Second).String()
		result = append(result, PodInfo{
			Name:      pod.Name,
			Status:    string(pod.Status.Phase),
			Namespace: pod.Namespace,
			Age:       age,
		})
	}
	return result, nil
}



func DeleteTenantDB(namespace, dbName string) error {
	fmt.Printf("Deleting database %s in namespace %s\n", dbName, namespace)

	// Delete the PostgreSQL cluster using kubectl (since it's a CRD)
	cmd := exec.Command("kubectl", "delete", "postgresql", dbName, "-n", namespace)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to delete PostgreSQL cluster: %v, output: %s\n", err, string(output))
		return fmt.Errorf("failed to delete PostgreSQL cluster %s: %w", dbName, err)
	}

	fmt.Printf("Successfully deleted PostgreSQL cluster %s\n", dbName)

	// Wait a moment for the operator to clean up
	time.Sleep(2 * time.Second)

	// The Zalando operator should automatically clean up associated secrets
	// But we can also manually delete them if needed
	clientset, err := getKubeClient()
	if err != nil {
		return fmt.Errorf("failed to get k8s client: %w", err)
	}

	// Delete associated secrets (optional, operator usually handles this)
	secretNames := []string{
		fmt.Sprintf("%s.%s.credentials.postgresql.acid.zalan.do", dbName, dbName),
		fmt.Sprintf("postgres.%s.credentials.postgresql.acid.zalan.do", dbName),
	}

	for _, secretName := range secretNames {
		err := clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Could not delete secret %s (may not exist): %v\n", secretName, err)
		} else {
			fmt.Printf("Deleted secret %s\n", secretName)
		}
	}

	return nil
}


type DatabaseClusterInfo struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Status            string            `json:"status"`
	DetailedStatus    string            `json:"detailed_status"`
	CredentialsReady  bool              `json:"credentials_ready"`
	ConnectionReady   bool              `json:"connection_ready"`
	CreatedAt         string            `json:"created_at"`
	Replicas          int               `json:"replicas"`
	RunningReplicas   int               `json:"running_replicas"`
	ConnectionInfo    map[string]string `json:"connection_info,omitempty"`
	CreationMethod    string            `json:"creation_method"` // "zalando" or "manual"
}

func CheckTenantDBStatus(namespace, dbName string) (string, error) {
	cluster, err := GetDatabaseClusterInfo(namespace, dbName)
	if err != nil {
		return "Unknown", err
	}
	return cluster.Status, nil
}

func GetDatabaseClusterInfo(namespace, dbName string) (*DatabaseClusterInfo, error) {
	clientset, err := getKubeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client: %w", err)
	}

	cluster := &DatabaseClusterInfo{
		Name:      dbName,
		Namespace: namespace,
		Status:    "Unknown",
	}

	// 1. Get PostgreSQL CRD information
	cmd := exec.Command("kubectl", "get", "postgresql", dbName, "-n", namespace, "-o", "json")
	output, err := cmd.Output()
	
	if err != nil {
		cluster.Status = "Not Found"
		cluster.DetailedStatus = "PostgreSQL cluster not found"
		return cluster, nil
	}

	// Parse basic info from CRD
	outputStr := string(output)
	if strings.Contains(outputStr, `"PostgresClusterStatus":"Running"`) {
		cluster.Status = "Running"
		cluster.DetailedStatus = "Database cluster is running"
	} else if strings.Contains(outputStr, `"PostgresClusterStatus":"Creating"`) {
		cluster.Status = "Creating"
		cluster.DetailedStatus = "Database cluster is being created"
	} else if strings.Contains(outputStr, `"PostgresClusterStatus":"CreateFailed"`) {
		cluster.Status = "Failed"
		cluster.DetailedStatus = "Database creation failed"
	} else {
		cluster.Status = "Pending"
		cluster.DetailedStatus = "Database cluster is pending"
	}

	// 2. Get pod information
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("cluster-name=%s", dbName),
	})
	
	if err == nil && len(pods.Items) > 0 {
		runningCount := 0
		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodRunning {
				runningCount++
			}
			// Get creation time from first pod only
			if cluster.CreatedAt == "" {
				cluster.CreatedAt = pod.CreationTimestamp.Format("2006-01-02 15:04:05")
			}
		}
		cluster.Replicas = len(pods.Items)
		cluster.RunningReplicas = runningCount
		
		// Status update based on pods and credentials
		if runningCount > 0 {
			cluster.Status = "Ready"
			cluster.DetailedStatus = "Database is ready"
		}
	}

	// 3. Check for Zalando credentials
	ownerSecretName := fmt.Sprintf("%s.%s.credentials.postgresql.acid.zalan.do", dbName, dbName)
	_, err = clientset.CoreV1().Secrets(namespace).Get(context.TODO(), ownerSecretName, metav1.GetOptions{})
	if err == nil {
		cluster.CredentialsReady = true
		cluster.CreationMethod = "zalando"
		cluster.ConnectionReady = true
		cluster.Status = "Ready"
		cluster.DetailedStatus = "Database is ready with credentials"
	}

	return cluster, nil
}


func ProvisionTenantDB(namespace, dbName string, replicas int) error {
	clientset, err := getKubeClient()
	if err != nil {
		return fmt.Errorf("failed to get k8s client: %w", err)
	}

	// 1. Create Namespace if not exists
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create namespace: %w", err)
		}
	}

	// 2. Render and apply manifest
	data := TemplateData{
		Namespace: namespace,
		DBName:    strings.ToLower(dbName),
		Team:      "paas-team",
		Replicas:  replicas,
	}

	tmplPath := filepath.Join("templates", "postgres-cluster.yaml.tmpl")
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New("postgres").Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = &buf
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply manifest: %w", err)
	}

	// 3. Wait for pod to be ready
	fmt.Printf("Waiting for PostgreSQL pod to be ready...\n")
	
	// Wait up to 2 minutes for pod to be running
	podReady := false
	for attempts := 0; attempts < 60; attempts++ {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("cluster-name=%s", dbName),
		})
		
		if err == nil && len(pods.Items) > 0 {
			for _, pod := range pods.Items {
				if pod.Status.Phase == corev1.PodRunning {
					podReady = true
					break
				}
			}
		}
		
		if podReady {
			break
		}
		
		// Show progress every 15 seconds
		if attempts > 0 && attempts%8 == 0 {
			fmt.Printf("Still waiting for pod... (%d seconds elapsed)\n", attempts*2)
		}
		
		time.Sleep(2 * time.Second)
	}

	if !podReady {
		return fmt.Errorf("pod not ready after 2 minutes")
	}

	fmt.Printf("Pod is ready, checking Zalando credentials...\n")

	// 4. Wait for Zalando credentials
	credentials, err := GetDatabaseCredentials(namespace, dbName, 60*time.Second)
	if err != nil {
		return fmt.Errorf("failed to get database credentials: %v", err)
	}

	fmt.Printf("Database creation successful!\n")
	fmt.Printf("Database: %s\n", credentials["database_name"])
	fmt.Printf("Connection string: %s\n", credentials["connection_string"])
	return nil
}

// ProvisionTenantDBWithCredentials provisions a database and returns credentials immediately
func ProvisionTenantDBWithCredentials(namespace, dbName string, replicas int) (map[string]interface{}, error) {
	clientset, err := getKubeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client: %w", err)
	}

	// 1. Create Namespace if not exists (with retry)
	fmt.Printf("Ensuring namespace %s exists...\n", namespace)
	
	var nsErr error
	for attempts := 0; attempts < 3; attempts++ {
		_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
		if err == nil {
			fmt.Printf("Namespace %s already exists\n", namespace)
			break
		}
		
		// Try to create namespace
		_, nsErr = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}, metav1.CreateOptions{})
		
		if nsErr == nil {
			fmt.Printf("Created namespace %s\n", namespace)
			break
		}
		
		if attempts < 2 {
			fmt.Printf("Namespace creation attempt %d failed, retrying... %v\n", attempts+1, nsErr)
			time.Sleep(2 * time.Second)
		}
	}
	
	if nsErr != nil && err != nil {
		return nil, fmt.Errorf("failed to create namespace after 3 attempts: %w", nsErr)
	}

	// 2. Render and apply manifest
	data := TemplateData{
		Namespace: namespace,
		DBName:    strings.ToLower(dbName),
		Team:      "paas-team",
		Replicas:  replicas,
	}

	tmplPath := filepath.Join("templates", "postgres-cluster.yaml.tmpl")
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New("postgres").Parse(string(tmplBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = &buf
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to apply manifest: %w", err)
	}

	fmt.Printf("Database manifest applied successfully!\n")

	// 3. Return immediate credentials based on Zalando naming conventions
	// The actual credentials will be created by Zalando operator in the background
	result := map[string]interface{}{
		"database_name": dbName,
		"host": fmt.Sprintf("%s.%s.svc.cluster.local", dbName, namespace),
		"port": "5432",
		"status": "provisioning",
		"message": "Database is being created. Credentials will be available shortly.",
		"secret_name": fmt.Sprintf("%s.%s.credentials.postgresql.acid.zalan.do", dbName, dbName),
		"connection_info": map[string]string{
			"host": fmt.Sprintf("%s.%s.svc.cluster.local", dbName, namespace),
			"port": "5432",
			"database": dbName,
			"ssl_mode": "prefer",
			"note": "Username and password will be available in the secret once ready",
		},
		"instructions": map[string]string{
			"check_status": fmt.Sprintf("kubectl get postgresql %s -n %s", dbName, namespace),
			"get_credentials": fmt.Sprintf("kubectl get secret %s.%s.credentials.postgresql.acid.zalan.do -n %s -o yaml", dbName, dbName, namespace),
		},
	}

	fmt.Printf("Database creation initiated for %s in namespace %s\n", dbName, namespace)
	return result, nil
}




func GetDatabaseCredentials(namespace, dbName string, timeout time.Duration) (map[string]interface{}, error) {
    clientset, err := getKubeClient()
    if err != nil {
        return nil, fmt.Errorf("failed to get k8s client: %w", err)
    }

    // The main database owner credentials
    ownerSecretName := fmt.Sprintf("%s.%s.credentials.postgresql.acid.zalan.do", dbName, dbName)
    
    fmt.Printf("Looking for secret: %s in namespace: %s\n", ownerSecretName, namespace)
    
    deadline := time.Now().Add(timeout)
    var ownerCreds map[string]string

    // Wait for owner credentials
    attempt := 0
    for {
        attempt++
        
        secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), ownerSecretName, metav1.GetOptions{})
        if err == nil {
            fmt.Printf("Found owner secret: %s after %d attempts\n", ownerSecretName, attempt)
            ownerCreds = make(map[string]string)
            for key, val := range secret.Data {
                ownerCreds[key] = string(val)
            }
            break
        }
        
        if time.Now().After(deadline) {
            return nil, fmt.Errorf("owner secret %s not found after %v: %w", ownerSecretName, timeout, err)
        }
        
        time.Sleep(2 * time.Second)
    }

    // Prepare response with connection information
    result := map[string]interface{}{
        "database_name": dbName,
        "host": fmt.Sprintf("%s.%s.svc.cluster.local", dbName, namespace),
        "port": "5432",
        "primary_user": ownerCreds,
    }

    // Add connection string for convenience
    if username, ok := ownerCreds["username"]; ok {
        if password, ok := ownerCreds["password"]; ok {
            connectionString := fmt.Sprintf("postgresql://%s:%s@%s.%s.svc.cluster.local:5432/%s", 
                username, password, dbName, namespace, dbName)
            result["connection_string"] = connectionString
            result["connection_string_ssl"] = connectionString + "?sslmode=prefer"
            
            // Add individual components for easier use
            result["connection_info"] = map[string]string{
                "host": fmt.Sprintf("%s.%s.svc.cluster.local", dbName, namespace),
                "port": "5432",
                "database": dbName,
                "username": username,
                "password": password,
                "ssl_mode": "prefer",
            }
        }
    }

    return result, nil
}



func getKubeClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// For local development, always use kubeconfig file first
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")
	if _, err := os.Stat(kubeconfig); err == nil {
		fmt.Printf("Using kubeconfig file: %s\n", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("Failed to load kubeconfig file: %v\n", err)
		}
	}

	// Fallback to in-cluster config if kubeconfig failed
	if config == nil {
		fmt.Printf("Trying in-cluster configuration...\n")
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load both kubeconfig and in-cluster config: %w", err)
		}
	}

	// Set a reasonable timeout
	config.Timeout = 30 * time.Second

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Test the connection
	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kubernetes API server: %w", err)
	}

	return clientset, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // Windows support
}

// Helper function to generate random password
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
