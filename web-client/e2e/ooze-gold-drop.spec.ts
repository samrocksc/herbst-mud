import { test, expect } from "@playwright/test";

const API_BASE = "http://localhost:8080";

async function getAuthToken(): Promise<string> {
  const res = await fetch(`${API_BASE}/users/auth`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: "sma", password: "sma" }),
  });
  const data = await res.json();
  return data.token;
}

test.describe("Ooze gold drop in Junkyard", () => {
  test("kill an ooze and verify gold drop", async () => {
    const token = await getAuthToken();

    // Get NPCs in room 92 (Junkyard Center) — use the npc-instances endpoint
    // which is admin-authenticated
    const npcsRes = await fetch(
      `${API_BASE}/api/npc-instances?world_id=2&active=true`,
      {
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    const npcsData = await npcsRes.json();

    // Find an ooze in room 92
    const oozes = (npcsData || []).filter(
      (n: { name: string; room_id: number }) =>
        n.name.includes("Ooze") && n.room_id === 92,
    );
    console.log(`Oozes in room 92: ${oozes.length}`);

    // If no ooze in room 92, check all Junkyard rooms
    let ooze = oozes[0];
    if (!ooze) {
      const allOozes = (npcsData || []).filter((n: { name: string }) =>
        n.name.includes("Ooze"),
      );
      console.log(
        `All active oozes: ${allOozes.length} — rooms: ${allOozes.map((o: { room_id: number }) => o.room_id).join(", ")}`,
      );
      expect(allOozes.length).toBeGreaterThan(0);
      // Move smack to the ooze's room
      ooze = allOozes[0];
      await fetch(`${API_BASE}/characters/5`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ current_room_id: ooze.room_id }),
      });
    } else {
      // Move smack to room 92
      await fetch(`${API_BASE}/characters/5`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ current_room_id: 92 }),
      });
    }

    console.log(
      `Targeting ${ooze.name} (id=${ooze.id}) in room ${ooze.room_id}`,
    );

    // Kill the ooze via damage API
    const killRes = await fetch(
      `${API_BASE}/characters/${ooze.id}/damage`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ damage: 9999, attacker_id: 5 }),
      },
    );
    const killData = await killRes.json();
    console.log(`Kill response: ${JSON.stringify(killData)}`);

    // Verify the NPC was defeated
    expect(killData.defeated).toBe(true);

    // Verify gold was dropped
    expect(killData.gold_dropped).toBeDefined();
    expect(killData.gold_dropped).toBeGreaterThan(0);
    console.log(`Gold dropped: ${killData.gold_dropped}`);

    console.log(
      `SUCCESS: ${ooze.name} (id=${ooze.id}) defeated, dropped ${killData.gold_dropped} gold!`,
    );
  });

  test("verify multiple oozes are alive in Junkyard", async () => {
    const token = await getAuthToken();

    const npcsRes = await fetch(
      `${API_BASE}/api/npc-instances?world_id=2&active=true`,
      {
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    const npcsData = await npcsRes.json();

    const oozes = (npcsData || []).filter((n: { name: string }) =>
      n.name.includes("Ooze"),
    );
    console.log(
      `Active oozes: ${oozes.length} — ${oozes.map((o: { id: number; name: string; room_id: number }) => `${o.name}(${o.id}) in room ${o.room_id}`).join(", ")}`,
    );

    expect(oozes.length).toBeGreaterThan(0);
  });
});