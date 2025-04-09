package src

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func GetResourceEvents(conn *pgx.Conn, resourceType, resourceName string) ([]v1.Event, error) {
	query := `
		SELECT resource, resource_name, reason, message, event_time
		FROM kubernetes_events
		WHERE LOWER(resource) = LOWER($1)
		AND LOWER(resource_name) = LOWER($2)
		AND event_time >= NOW() - INTERVAL '10 minutes'
		ORDER BY event_time DESC
	`

	rows, err := conn.Query(context.Background(), query, resourceType, resourceName)
	if err != nil {
		return nil, fmt.Errorf("error querying events: %v", err)
	}
	defer rows.Close()

	var events []v1.Event
	for rows.Next() {
		var event v1.Event
		var eventTime time.Time
		err := rows.Scan(
			&event.InvolvedObject.Kind,
			&event.InvolvedObject.Name,
			&event.Reason,
			&event.Message,
			&eventTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning event: %v", err)
		}
		event.EventTime.Time = metav1.NewTime(eventTime).Time
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %v", err)
	}

	return events, nil
}
