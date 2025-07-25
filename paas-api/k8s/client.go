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
	// No need to call getKubeClient since we're using kubectl directly
	cmd := exec.Command("kubectl", "get", "postgresql", dbName, "-n", namespace, "-o", "jsonpath={.status.PostgresClusterStatus}")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}
	return string(output), nil
}



func ProvisionTenantDB(namespace, dbName, password string) error {
	clientset, err := getKubeClient()
	if err != nil {
		return fmt.Errorf("failed to get k8s client: %w", err)
	}

	// 1. Create Namespace
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}, metav1.CreateOptions{})

		if err != nil {
			return fmt.Errorf("failed to create namespace: %w", err)
		}
	}

	// 2. Create credentials secret
	secretName := fmt.Sprintf("postgres.%s.credentials.postgresql.acid.zalan.do", dbName)
	_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		StringData: map[string]string{
			"username": dbName,
			"password": password,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	// 3. Render and apply manifest
	data := TemplateData{
		Namespace: namespace,
		DBName:    strings.ToLower(dbName),
		DBUser:    strings.ToLower(dbName),
		Password:  password,
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
