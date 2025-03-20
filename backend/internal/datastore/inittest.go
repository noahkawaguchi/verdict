//go:build test

package datastore

// init sets up the dbClient before main executes, once per cold start.
func init() {
	// Set up AWS config for local development without SAM
	localClientSetup("http://localhost:8000")
}
