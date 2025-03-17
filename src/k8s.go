package src

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"time"
)

func LoadK8sConfig() *rest.Config {
	var config *rest.Config
	var err error
	if os.Getenv("ENV") == "kubernetes" {
		log.Println("In k8s env")
		config, err = rest.InClusterConfig()
	} else if os.Getenv("ENV") == "local" {
		log.Println("In local env")
		config, err = clientcmd.BuildConfigFromFlags("", "/Users/gbehl/.kube/config")
	}
	if err != nil {
		log.Fatal("Failed to load in-cluster config:", err)
	}
	return config
}

func CreateK8sClient(config *rest.Config) *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Failed to create Kubernetes client:", err)
	}
	return clientset
}

func WatchEvents(clientset *kubernetes.Clientset, conn *pgx.Conn) {
	// Create a watch request for Kubernetes events
	watcher, err := clientset.CoreV1().Events("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to set up event watcher:", err)
	}

	// Process events as they occur
	fmt.Println("Watching Kubernetes events...")

	for event := range watcher.ResultChan() {
		ev, ok := event.Object.(*v1.Event)
		if !ok {
			log.Println("Unexpected event type received")
			continue
		}

		// Print event details
		fmt.Printf("[%s] Event: %s - Reason: %s - Message: %s\n",
			ev.ObjectMeta.CreationTimestamp.Time.Format(time.RFC3339), ev.InvolvedObject.Kind, ev.Reason, ev.Message)
		// save events to db
		saveEventToDB(conn, ev)
	}
	defer conn.Close(context.Background())
}
