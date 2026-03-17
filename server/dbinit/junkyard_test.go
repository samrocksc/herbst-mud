package dbinit

import (
	"context"
	"herbst-server/db"
	"herbst-server/db/room"
	"testing"

	_ "herbst-server/db/runtime"
)

func TestInitJunkyard(t *testing.T) {
	client, err := db.Open("sqlite3", "file:junkyard_test?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer client.Close()

	// Run migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// First initialize the cross-way and fountain (required for Junkyard to connect)
	if err := InitCrossWay(client); err != nil {
		t.Fatalf("failed to init crossway: %v", err)
	}
	if err := InitFountain(client); err != nil {
		t.Fatalf("failed to init fountain: %v", err)
	}

	// Run InitJunkyard
	if err := InitJunkyard(client); err != nil {
		t.Fatalf("failed to init junkyard: %v", err)
	}

	// Verify the entrance room was created
	ctx := context.Background()
	entrance, err := client.Room.Query().Where(room.NameEQ("Junkyard Entrance")).Only(ctx)
	if err != nil {
		t.Fatalf("failed to find Junkyard Entrance: %v", err)
	}
	if entrance == nil {
		t.Fatal("Junkyard Entrance room is nil")
	}

	// Verify exits from entrance
	if entrance.Exits == nil {
		t.Fatal("Entrance room has no exits")
	}

	// Check that it connects to Fountain Plaza (west exit)
	if entrance.Exits["west"] == 0 {
		t.Error("Junkyard Entrance should have west exit to Fountain Plaza")
	}

	// Verify we have 25 rooms total (5x5 grid)
	allRooms, err := client.Room.Query().All(ctx)
	if err != nil {
		t.Fatalf("failed to query rooms: %v", err)
	}

	// We should have: 5 (crossway) + 2 (fountain + fountain plaza) + 25 (junkyard) = 32 rooms
	if len(allRooms) < 25 {
		t.Errorf("Expected at least 25 junkyard rooms, got %d", len(allRooms))
	}

	// Verify the exit room was created
	exit, err := client.Room.Query().Where(room.NameEQ("Junkyard Exit")).Only(ctx)
	if err != nil {
		t.Fatalf("failed to find Junkyard Exit: %v", err)
	}
	if exit == nil {
		t.Fatal("Junkyard Exit room is nil")
	}

	t.Logf("Junkyard test passed: created %d rooms total", len(allRooms))
}

func TestInitJunkyardIdempotent(t *testing.T) {
	client, err := db.Open("sqlite3", "file:junkyard_idempotent_test?mode=memory&_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer client.Close()

	// Run migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// First initialize the cross-way and fountain
	if err := InitCrossWay(client); err != nil {
		t.Fatalf("failed to init crossway: %v", err)
	}
	if err := InitFountain(client); err != nil {
		t.Fatalf("failed to init fountain: %v", err)
	}

	// Run InitJunkyard twice
	if err := InitJunkyard(client); err != nil {
		t.Fatalf("failed to init junkyard first time: %v", err)
	}

	// Count rooms before second run
	ctx := context.Background()
	roomsBefore, _ := client.Room.Query().Count(ctx)

	// Run again - should be idempotent (skip)
	if err := InitJunkyard(client); err != nil {
		t.Fatalf("failed to init junkyard second time: %v", err)
	}

	roomsAfter, _ := client.Room.Query().Count(ctx)

	if roomsBefore != roomsAfter {
		t.Errorf("InitJunkyard is not idempotent: %d rooms before, %d after", roomsBefore, roomsAfter)
	}

	t.Logf("Idempotency test passed: %d rooms (unchanged)", roomsAfter)
}