//go:build ignore

package main

import (
	"log"

	"entgo.io/ent/cmd/ent"
	"github.com/agext/levenshtein"
)

func main() {
	// Generate ent code
	cmd := ent.NewGenerateCmd()
	if err := cmd.Run(); err != nil {
		log.Fatalf("ent generate failed: %v", err)
	}
}
