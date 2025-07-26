package k8s

import (
	"bytes"
	"context"
	"fmt"
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
	DBUser    string
	Password  string
	Team      string
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
	clientset, err := getKubeClient()
	if err != nil {
		return fmt.Errorf("failed to get k8s client: %w", err)
	}

	// Option 1: Delete entire namespace (DB + secrets + CRs)
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}
	// Option 2: If keeping namespace, delete DB + secret:
	// _ = clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	// _ = exec.Command("kubectl", "delete", "postgresql", dbName, "-n", namespace).Run()

	return nil
}


func CheckTenantDBStatus(namespace, dbName string) (string, error) {
	// 1. Get PostgresClusterStatus from the CRD
	cmd := exec.Command("kubectl", "get", "postgresql", dbName, "-n", namespace, "-o", "jsonpath={.status.PostgresClusterStatus}")
	output, err := cmd.Output()
	status := strings.TrimSpace(string(output))

	// If the CRD status is missing or CreateFailed, check pods directly
	if err != nil || status == "" || status == "CreateFailed" {
		// 2. Check if the DB pod is actually running
		podCheck := exec.Command("kubectl", "get", "pods", "-n", namespace, "-l", fmt.Sprintf("cluster-name=%s", dbName), "-o", "jsonpath={.items[*].status.phase}")
		podStatusBytes, podErr := podCheck.Output()
		podStatus := strings.TrimSpace(string(podStatusBytes))

		if podErr == nil && strings.Contains(podStatus, "Running") {
			return "Running (fallback from pod)", nil
		}

		// If pod check fails, return original error
		if err != nil {
			return "", fmt.Errorf("failed to get CRD status: %w", err)
		}
		return status, nil
	}

	return status, nil
}


func ProvisionTenantDB(namespace, dbName, password string) error {
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

	// 2. Create tenant user secret (separate from operator secrets)
	tenantSecretName := fmt.Sprintf("tenant-%s-credentials", dbName)
	_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: tenantSecretName,
		},
		StringData: map[string]string{
			"username": dbName,
			"password": password,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("failed to create tenant secret: %w", err)
	}

	// 3. Render and apply manifest (without embedding password or secret references)
	data := TemplateData{
		Namespace: namespace,
		DBName:    strings.ToLower(dbName),
		DBUser:    strings.ToLower(dbName),
		Password:  "",  //no pass password here, operator handles secrets internally
		Team:      "paas-team",
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

	return nil
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
