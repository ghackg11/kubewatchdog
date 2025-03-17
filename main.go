package main

import (
	"golearn/src"
)

func main() {
	// Create a Kubernetes client
	config := src.LoadK8sConfig()
	clientSet := src.CreateK8sClient(config)
	conn := src.ConnectToDB()

	// Watch for Events
	src.WatchEvents(clientSet, conn)
	select {}
}
