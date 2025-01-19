package main

import (
	"fmt"
)

func internalKubernetesEndpoint(serviceName, namespace string, port int) string {
	return fmt.Sprintf("https://%s.%s.svc.cluster.local:%d", serviceName, namespace, port)
}
