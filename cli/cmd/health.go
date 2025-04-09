package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"golearn/src"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health <resource-type> <resource-name>",
	Short: "Get health summary for a specific resource",
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
		eventString := ""
		for _, event := range events {
			eventString += fmt.Sprintf("[%s] %s - %s\n",
				event.EventTime.Time.Format(time.RFC3339),
				event.Reason,
				event.Message)
		}
		prompt := "You are kubewatch, a observability tool for kubernetes. You need to access the events for the kubernetes object given to you and check if you see an issue. if there is an issue, summarise that, and tell the possible cause, otherwise, report that the object is healthy. These are the events for " + resourceType + "/" + resourceName + " in the last 10 minutes, Events: " + eventString
		response, err := src.GenerateLlmResponse(prompt)
		if err != nil {
			log.Fatalf("Error generating response from LLM: %v", err)
		}
		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
