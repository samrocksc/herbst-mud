package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// 1. Read all equipment rows
	rows, err := db.Query(`
		SELECT id, name, description, slot, level, weight, item_type, color, is_visible, is_immovable,
		       effect_type, effect_value, effect_duration, healing, effect, is_container, container_capacity, is_locked,
		       equipment_template_id
		FROM equipment ORDER BY id`)
	if err != nil {
		log.Fatal("select equipment:", err)
	}
	defer rows.Close()

	type eq struct {
		id       int
		name     string
		desc     string
		slot     string
		level    int
		weight   int
		itemType string
		color    string
		isVis    bool
		isImm    bool
		et       string
		ev       int
		ed       int
		healing  int
		efct    string
		isCont  bool
		contCap int
		isLock  bool
		etID    sql.NullString
	}

	var equipmentRows []eq
	for rows.Next() {
		var r eq
		err := rows.Scan(&r.id, &r.name, &r.desc, &r.slot, &r.level, &r.weight, &r.itemType,
			&r.color, &r.isVis, &r.isImm, &r.et, &r.ev, &r.ed, &r.healing,
			&r.efct, &r.isCont, &r.contCap, &r.isLock, &r.etID,
		)
		if err != nil {
			log.Fatal("scan:", err)
		}
		equipmentRows = append(equipmentRows, r)
	}
	rows.Close()
	fmt.Printf("Found %d equipment rows\n", len(equipmentRows))

	// 2. For each row, create a template (if not exists) and assign
	for _, r := range equipmentRows {
		// Generate deterministic template ID slug
		tmplID := slugify(r.name)
		fmt.Printf("  Equipment #%d -> template '%s'\n", r.id, tmplID)

		// Check if template exists
		var exists int
		if err := db.QueryRow(`SELECT 1 FROM equipment_templates WHERE id = $1`, tmplID).Scan(&exists); err != nil {
			// doesn't exist → create it
			_, err := db.Exec(`
				INSERT INTO equipment_templates
				(id, name, description, slot, level, weight, item_type, color,
				 is_visible, is_immovable, effect_type, effect_value, effect_duration,
				 is_container, container_capacity, is_locked, key_item_id, reveal_condition)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
				tmplID, r.name, r.desc, r.slot, r.level, r.weight, r.itemType, r.color,
				r.isVis, r.isImm, r.et, r.ev, r.ed,
				r.isCont, r.contCap, r.isLock, "", "",
			)
			if err != nil {
				log.Fatal("insert template:", err)
			}
			fmt.Printf("    Created template '%s'\n", tmplID)
		} else {
			fmt.Printf("    Reusing existing template '%s'\n", tmplID)
		}

		// Link equipment to template
		_, err := db.Exec(`UPDATE equipment SET equipment_template_id = $1 WHERE id = $2`, tmplID, r.id)
		if err != nil {
			log.Fatal("update equipment:", err)
		}
		fmt.Printf("    Linked equipment #%d to template '%s'\n", r.id, tmplID)
	}

	fmt.Println("\nMigration completed successfully.")
}

func slugify(s string) string {
	// Lowercase, replace spaces with underscores, strip non-alphanumeric
	s = strings.ToLower(s)
	var b strings.Builder
	for _, ch := range s {
		switch {
		case ch >= 'a' && ch <= 'z':
			b.WriteRune(ch)
		case ch >= '0' && ch <= '9':
			b.WriteRune(ch)
		case ch == ' ' || ch == '-' || ch == '_':
			b.WriteByte('_')
		}
	}
	result := b.String()
	// collapse multiple underscores
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	result = strings.Trim(result, "_")
	return result
}
