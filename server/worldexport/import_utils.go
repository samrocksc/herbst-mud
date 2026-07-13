package worldexport

import (
	"fmt"
	"strconv"
)

type idMaps struct {
	Rooms             map[int]int
	Characters        map[int]int
	Abilities         map[int]int
	Races             map[int]int
	Genders           map[int]int
	Factions          map[int]int
	FactionCategories map[int]int
	Tags              map[int]int
	Equipment         map[int]int
	NPCTemplates      map[string]string
	Zones             map[string]string
	Triggers          map[int]int
	Quests            map[int]int
	Recipes           map[int]int
	DialogNodes       map[string]string
	SocialCommands    map[int]int
	ShopTemplates     map[int]int
	EffectHooks       map[int]int
	Effects           map[int]int // standalone Effects are global, ID stays same
	Skills            map[int]int
}

func newIDMaps() *idMaps {
	return &idMaps{
		Rooms:             make(map[int]int),
		Characters:        make(map[int]int),
		Abilities:         make(map[int]int),
		Races:             make(map[int]int),
		Genders:           make(map[int]int),
		Factions:          make(map[int]int),
		FactionCategories: make(map[int]int),
		Tags:              make(map[int]int),
		Equipment:         make(map[int]int),
		NPCTemplates:      make(map[string]string),
		Zones:             make(map[string]string),
		Triggers:          make(map[int]int),
		Quests:            make(map[int]int),
		Recipes:           make(map[int]int),
		DialogNodes:       make(map[string]string),
		SocialCommands:    make(map[int]int),
		ShopTemplates:     make(map[int]int),
		EffectHooks:       make(map[int]int),
		Effects:           make(map[int]int),
		Skills:            make(map[int]int),
	}
}

func intVal(v interface{}) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}

func intValOr(v interface{}, def int) int {
	if v == nil {
		return def
	}
	r := intVal(v)
	if r == 0 {
		return def
	}
	return r
}

func strVal(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case nil:
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func strValOr(v interface{}, def string) string {
	if v == nil || strVal(v) == "" {
		return def
	}
	return strVal(v)
}

func strPtr(v interface{}) *string {
	if v == nil {
		return nil
	}
	s := strVal(v)
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(v int) *int {
	return &v
}

func intPtrVal(v interface{}) *int {
	if v == nil {
		return nil
	}
	n := intVal(v)
	return &n
}

func boolVal(v interface{}) bool {
	switch b := v.(type) {
	case bool:
		return b
	case string:
		return b == "true"
	}
	return false
}

func floatValOr(v interface{}, def float64) float64 {
	switch f := v.(type) {
	case float64:
		return f
	case int:
		return float64(f)
	}
	return def
}

func mapVal(v interface{}) map[string]interface{} {
	switch m := v.(type) {
	case map[string]interface{}:
		return m
	}
	return nil
}

func mapValInt(v interface{}, def map[string]int) map[string]int {
	m := mapVal(v)
	if m == nil {
		return def
	}
	out := make(map[string]int, len(m))
	for k, val := range m {
		out[k] = intVal(val)
	}
	return out
}

func mapPtrVal(v interface{}) *map[string]interface{} {
	if v == nil {
		return nil
	}
	m := mapVal(v)
	if m == nil {
		return nil
	}
	return &m
}

// edgeID extracts an int ID from a nested edge object in the exported JSON.
// e.g. edgeID(m["edges"], "ability") returns the ability ID from edges.ability.id
func edgeID(edges interface{}, key string) int {
	em := mapVal(edges)
	if em == nil {
		return 0
	}
	node := mapVal(em[key])
	if node == nil {
		return 0
	}
	return intVal(node["id"])
}

// edgeStrID extracts a string ID from a nested edge object.
func edgeStrID(edges interface{}, key string) string {
	em := mapVal(edges)
	if em == nil {
		return ""
	}
	node := mapVal(em[key])
	if node == nil {
		return ""
	}
	return strVal(node["id"])
}

func strSliceVal(v interface{}, def []string) []string {
	switch s := v.(type) {
	case []interface{}:
		out := make([]string, 0, len(s))
		for _, item := range s {
			out = append(out, strVal(item))
		}
		return out
	case []string:
		return s
	}
	return def
}

func strSlicePtrVal(v interface{}) *[]string {
	if v == nil {
		return nil
	}
	s := strSliceVal(v, nil)
	if s == nil {
		return nil
	}
	return &s
}
