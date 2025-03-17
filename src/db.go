package src

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	v1 "k8s.io/api/core/v1"
	"log"
	"os"
)

func ConnectToDB() *pgx.Conn {
	var dbUrl string
	if os.Getenv("ENV") == "local" {
		dbUrl = "postgres://admin:mysecurepassword@127.0.0.1:5432/kubewatchdog?sslmode=disable"
	} else {
		dbUrl = "postgres://admin:mysecurepassword@timescaledb:5432/kubewatchdog?sslmode=disable"
	}
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Failed to connect to TimescaleDB: %v", err)
	}
	return conn
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
