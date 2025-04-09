package cmd

import (
	"fmt"
	"log"

	"golearn/src"

	"github.com/spf13/cobra"
)

var showAll bool

// helloCmd represents the hello command
var listResourcesCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes resources",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing Resources...")

		// Load Kubernetes config
		config := src.LoadK8sConfig()

		// Create Kubernetes client
		clientset := src.CreateK8sClient(config)

		// List resource types
		resourceTypes, err := src.ListK8sResourceTypes(clientset, showAll)
		if err != nil {
			log.Fatalf("Error listing resources: %v", err)
		}

		if showAll {
			fmt.Println("All Available Resource Types:")
		} else {
			fmt.Println("Core Resource Types:")
		}
		for _, resourceType := range resourceTypes {
			fmt.Printf("- %s\n", resourceType)
		}
	},
}

func init() {
	rootCmd.AddCommand(listResourcesCmd)
	listResourcesCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all resource types including API groups")
}
