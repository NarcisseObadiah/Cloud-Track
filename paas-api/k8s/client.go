package k8s

import "fmt"


func ProvisionTenantDB(namespace, dbName, password string) error {
	// TODO: Use client-go to:
	// - Create namespace
	// - Apply resource quota
	// - Create secrets
	// - Deploy StatefulSet + PVC + Service
	fmt.Printf("Would provision DB %s in namespace %s with password %s\n", dbName, namespace, password)
	return nil
}