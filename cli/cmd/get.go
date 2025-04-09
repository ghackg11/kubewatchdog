package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"golearn/src"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <resource-type> <resource-name>",
	Short: "Get events for a specific resource",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		resourceName := args[1]

		// Connect to database
		conn := src.ConnectToDB()
		defer conn.Close(context.Background())

		// Get events
		events, err := src.GetResourceEvents(conn, resourceType, resourceName)
		if err != nil {
			log.Fatalf("Error getting events: %v", err)
		}

		if len(events) == 0 {
			fmt.Printf("No events found for %s/%s in the last 10 minutes\n", resourceType, resourceName)
			return
		}

		fmt.Printf("Events for %s/%s in the last 10 minutes:\n", resourceType, resourceName)
		for _, event := range events {
			fmt.Printf("[%s] %s - %s\n",
				event.EventTime.Time.Format(time.RFC3339),
				event.Reason,
				event.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
