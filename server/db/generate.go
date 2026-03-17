//go:build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"herbst-server/db/schema"
)

func main() {
	ex, err := entc.GenerateGraph("./schema", &gen.Config{
		Schema:  schema.Schema,
	})
	if err != nil {
		log.Fatal("running entc generate:", err)
	}
	_ = ex
}