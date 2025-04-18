package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Prints a greeting",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, I am Kubewatch")
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
}
