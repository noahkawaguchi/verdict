//go:build dev

package datastore

import "os"

// init sets up the dbClient before main executes, once per cold start.
func init() {
	if os.Getenv("AWS_SAM_LOCAL") == "true" { // Running the Lambda function locally with SAM
		localClientSetup("http://host.docker.internal:8000")
	} else { // Just running the code in development, not using SAM
		localClientSetup("http://localhost:8000")
	}
	// Create the tables if they don't exist
	if !localTableExists(ballotsTableInfo) {
		createLocalTable(ballotsTableInfo, createBallotsTableInput)
	}
	if !localTableExists(pollsTableInfo) {
		createLocalTable(pollsTableInfo, createPollsTableInput)
	}
	printLocalTables()
}
