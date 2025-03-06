package main

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

func connectToDB() *pgx.Conn {
	dbUrl := "postgres://admin:mysecurepassword@timescaledb:5432/kubewatchdog?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Failed to connect to TimescaleDB: %v", err)
	}
	return conn
}

func loadK8sConfig() *rest.Config {
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

func createK8sClient(config *rest.Config) *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Failed to create Kubernetes client:", err)
	}
	return clientset
}

func main() {
	// Create a Kubernetes client
	config := loadK8sConfig()
	clientSet := createK8sClient(config)
	conn := connectToDB()

	// Watch for Events
	watchEvents(clientSet, conn)

	select {}
}

func saveEventToDB(conn *pgx.Conn, ev *v1.Event) {
	_, err := conn.Exec(
		context.Background(),
		`INSERT INTO kubernetes_events (
                               id, 
                               event_time, 
                               event_type, 
                               reason, message, 
                               namespace, 
                               resource, 
                               resource_name
            )
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (id, event_time) DO NOTHING;`,
		ev.UID,
		ev.ObjectMeta.CreationTimestamp.Time,
		ev.Type,
		ev.Reason,
		ev.Message,
		ev.Namespace,
		ev.InvolvedObject.Kind,
		ev.InvolvedObject.Name,
	)

	if err != nil {
		log.Printf("Failed to insert event into TimescaleDB: %v", err)
	} else {
		fmt.Printf("Saved event: %s - %s\n", ev.Reason, ev.Message)
	}

}

func watchEvents(clientset *kubernetes.Clientset, conn *pgx.Conn) {
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
