package main

import (
	"context"
	"fmt"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

func main() {
	// Define the containerd socket path
	socketPath := "/run/containerd/containerd.sock"

	// Create a new client connected to the default socket path
	client, err := containerd.New(socketPath)
	if err != nil {
		log.Fatalf("Failed to connect to containerd: %v", err)
	}
	defer client.Close()

	// List namespaces
	nsService := client.NamespaceService()
	namespaceList, err := nsService.List(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch namespaces: %v", err)
	}

	if len(namespaceList) == 0 {
		fmt.Println("No namespaces found.")
		return
	}

	fmt.Println("Namespaces and their containers:")
	for _, ns := range namespaceList {
		fmt.Printf("Namespace: %s\n", ns)

		// Set the current namespace in context
		ctx := namespaces.WithNamespace(context.Background(), ns)

		// List containers in the current namespace
		containers, err := client.Containers(ctx)
		if err != nil {
			log.Printf("Failed to fetch containers for namespace %s: %v", ns, err)
			continue
		}

		if len(containers) == 0 {
			fmt.Println("  No containers in this namespace.")
			continue
		}

		for _, container := range containers {
			fmt.Printf("  - %s\n", container.ID())
		}
	}
}
