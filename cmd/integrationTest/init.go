package main

import "github.com/GoogleCloudPlatform/functions-framework-go/functions"

func init() {
	// Register integration test with the Functions Framework
	functions.CloudEvent("IntegrationTest", IntegrationTest)
}
