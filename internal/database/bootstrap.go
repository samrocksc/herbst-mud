package database

import (
	"log"
)

// BootstrapGlobalState is a convenience function to load all existing data into global state
// This should be called once during server startup after database initialization
func BootstrapGlobalState(dbAdapter *DBAdapter) error {
	log.Println("Initializing global state tracking...")
	
	err := dbAdapter.InitializeGlobalState()
	if err != nil {
		log.Printf("Failed to initialize global state: %v", err)
		return err
	}
	
	log.Println("Global state tracking initialized successfully")
	return nil
}

// BootstrapGlobalStateForNewGame is a convenience function to initialize global state for a fresh game
// This creates default states for all existing rooms and characters
func BootstrapGlobalStateForNewGame(dbAdapter *DBAdapter) error {
	log.Println("Bootstrapping global state for new game...")
	
	err := dbAdapter.InitializeGlobalState()
	if err != nil {
		log.Printf("Failed to bootstrap global state: %v", err)
		return err
	}
	
	log.Println("Global state bootstrapped successfully")
	return nil
}

// RefreshGlobalState is a convenience function to refresh all global state data
// This can be used to re-sync global state if it gets out of sync
func RefreshGlobalState(dbAdapter *DBAdapter) error {
	log.Println("Refreshing global state tracking...")
	
	// Re-initialize all states
	err := dbAdapter.InitializeGlobalState()
	if err != nil {
		log.Printf("Failed to refresh global state: %v", err)
		return err
	}
	
	log.Println("Global state refreshed successfully")
	return nil
}