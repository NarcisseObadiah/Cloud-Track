package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ProvisionTenantDB(namespace, dbName, password string) error {
	clientset, err := getKubeClient()
	if err != nil {
		return fmt.Errorf("failed to get k8s client: %w", err)
	}

	// Test connectivity by listing namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	fmt.Printf("✅ Connected to cluster. Found %d namespaces.\n", len(namespaces.Items))
	fmt.Printf("📦 Would provision DB %q in namespace %q with password %q\n", dbName, namespace, password)

	// TODO: implement real provisioning (Namespace, Secret, PVC, StatefulSet, etc.)

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
