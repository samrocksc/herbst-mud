package adapters

import (
	"log"
	"os"
	"strings"
)

// debugMode checks if debug mode is enabled
var debugMode = strings.ToLower(os.Getenv("DEBUG")) == "true"

// debugLog logs debug messages only when DEBUG=true
func debugLog(format string, v ...interface{}) {
	if debugMode {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// infoLog always logs informational messages
func infoLog(format string, v ...interface{}) {
	log.Printf(format, v...)
}
