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

	// 1. Create Namespace if not exists
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create namespace: %w", err)
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

	// 3. Wait for pod to be ready
	fmt.Printf("Waiting for PostgreSQL pod to be ready...\n")
	
	// Wait up to 3 minutes for pod to be running
	podReady := false
	var lastPodStatus string
	for attempts := 0; attempts < 90; attempts++ {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("cluster-name=%s", dbName),
		})
		
		if err == nil && len(pods.Items) > 0 {
			for _, pod := range pods.Items {
				lastPodStatus = string(pod.Status.Phase)
				if pod.Status.Phase == corev1.PodRunning {
					// Check if containers are ready
					allContainersReady := true
					for _, status := range pod.Status.ContainerStatuses {
						if !status.Ready {
							allContainersReady = false
							break
						}
					}
					if allContainersReady {
						podReady = true
						break
					}
				}
				
				// Log any error conditions
				if pod.Status.Phase == corev1.PodFailed {
					return nil, fmt.Errorf("pod failed to start: %s", pod.Status.Message)
				}
			}
		}
		
		if podReady {
			break
		}
		
		// Show progress every 30 seconds
		if attempts > 0 && attempts%15 == 0 {
			fmt.Printf("Still waiting for pod... (%d seconds elapsed, current status: %s)\n", attempts*2, lastPodStatus)
		}
		
		time.Sleep(2 * time.Second)
	}

	if !podReady {
		return nil, fmt.Errorf("pod not ready after 3 minutes, last status: %s", lastPodStatus)
	}

	fmt.Printf("Pod is ready, waiting for Zalando credentials...\n")

	// 4. Wait for Zalando credentials (2 minutes should be enough)
	credentials, err := GetDatabaseCredentials(namespace, dbName, 2*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get database credentials: %v", err)
	}

	fmt.Printf("Database creation successful!\n")
	return credentials, nil
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
	// First, try in-cluster config (Kubernetes runtime)
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig file (local development)
		kubeconfig := filepath.Join(homeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
		}
	}

	return kubernetes.NewForConfig(config)
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
