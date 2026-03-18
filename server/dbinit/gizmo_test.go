package dbinit

import (
	"context"
	"testing"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/npctemplate"
	"herbst-server/db/room"
)

// TestInitGizmoNPC tests the Gizmo NPC initialization
func TestInitGizmoNPC(t *testing.T) {
	client, err := db.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Create schema
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create a fountain room first
	fountainRoom, err := client.Room.Create().
		SetName("The Fountain").
		SetDescription("A murky fountain").
		SetExits(map[string]int{}).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create fountain room: %v", err)
	}

	// Test InitGizmoNPC
	err = InitGizmoNPC(client)
	if err != nil {
		t.Fatalf("failed to init Gizmo NPC: %v", err)
	}

	// Verify NPC template was created
	template, err := client.NPCTemplate.Get(ctx, "gizmo")
	if err != nil {
		t.Fatalf("failed to get gizmo template: %v", err)
	}
	if template.Name != "Gizmo" {
		t.Errorf("expected template name 'Gizmo', got %q", template.Name)
	}
	if template.Race != "half-dog" {
		t.Errorf("expected template race 'half-dog', got %q", template.Race)
	}
	if template.Disposition != npctemplate.DispositionFriendly {
		t.Errorf("expected disposition 'friendly', got %v", template.Disposition)
	}
	if template.Greeting == "" {
		t.Error("expected greeting to be set")
	}

	// Verify character was created
	gizmoChars, err := client.Character.Query().
		Where(character.NameEQ("Gizmo")).
		All(ctx)
	if err != nil {
		t.Fatalf("failed to query Gizmo character: %v", err)
	}
	if len(gizmoChars) != 1 {
		t.Errorf("expected 1 Gizmo character, got %d", len(gizmoChars))
	}
	gizmo := gizmoChars[0]
	if !gizmo.IsNPC {
		t.Error("expected Gizmo to be NPC")
	}
	if gizmo.CurrentRoomId != fountainRoom.ID {
		t.Errorf("expected Gizmo in fountain room, got room %d", gizmo.CurrentRoomId)
	}
	if gizmo.Race != "half-dog" {
		t.Errorf("expected race 'half-dog', got %q", gizmo.Race)
	}

	// Test idempotency - running again should not error
	err = InitGizmoNPC(client)
	if err != nil {
		t.Fatalf("second init should be idempotent: %v", err)
	}

	// Should still only have one Gizmo
	gizmoChars, err = client.Character.Query().
		Where(character.NameEQ("Gizmo")).
		All(ctx)
	if err != nil {
		t.Fatalf("failed to query Gizmo character: %v", err)
	}
	if len(gizmoChars) != 1 {
		t.Errorf("expected 1 Gizmo after second init (idempotent), got %d", len(gizmoChars))
	}

	t.Log("All Gizmo NPC tests passed!")
}

// TestGizmoNPCInFountainRoom tests that Gizmo spawns in the fountain room
func TestGizmoNPCInFountainRoom(t *testing.T) {
	client, err := db.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Create schema
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Create rooms
	fountainRoom, err := client.Room.Create().
		SetName("The Fountain").
		SetDescription("A murky fountain").
		SetExits(map[string]int{}).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create fountain room: %v", err)
	}

	_, err = client.Room.Create().
		SetName("Fountain Plaza").
		SetDescription("A dusty plaza").
		SetIsStartingRoom(true).
		SetExits(map[string]int{}).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create fountain plaza: %v", err)
	}

	// Init Gizmo
	err = InitGizmoNPC(client)
	if err != nil {
		t.Fatalf("failed to init Gizmo NPC: %v", err)
	}

	// Verify Gizmo is in the fountain room
	roomChars, err := client.Room.Query().Where(room.NameEQ("The Fountain")).Only(ctx)
	if err != nil {
		t.Fatalf("failed to get fountain room: %v", err)
	}

	characters, err := roomChars.QueryCharacters().All(ctx)
	if err != nil {
		t.Fatalf("failed to query characters in fountain: %v", err)
	}

	foundGizmo := false
	for _, c := range characters {
		if c.Name == "Gizmo" {
			foundGizmo = true
			break
		}
	}
	if !foundGizmo {
		t.Error("Gizmo should be in the fountain room")
	}

	// Verify Gizmo has the NPC template edge
	gizmoChars, err := client.Character.Query().
		Where(character.NameEQ("Gizmo")).
		All(ctx)
	if err != nil {
		t.Fatalf("failed to query Gizmo: %v", err)
	}
	if len(gizmoChars) > 0 {
		template, err := gizmoChars[0].QueryNPCTemplate().Only(ctx)
		if err != nil {
			t.Fatalf("failed to query NPC template: %v", err)
		}
		if template.ID != "gizmo" {
			t.Errorf("expected template ID 'gizmo', got %q", template.ID)
		}
	}

	t.Log("Fountain room test passed!")
}