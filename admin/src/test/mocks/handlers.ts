import { http, HttpResponse } from "msw";
import {
  makeAbility,
  makeEquipmentTemplate,
  makeNPCTemplate,
  makeNPCInstance,
  makeQuest,
  makeWorld,
  makeCharacter,
  makeRace,
  makeTag,
  makeRoom,
} from "../factories";

export const handlers = [
  // Abilities
  http.get("/api/abilities", () => {
    return HttpResponse.json({ abilities: [makeAbility({ id: 1 }), makeAbility({ id: 2, name: "Icebolt" })] });
  }),
  http.post("/api/abilities", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeAbility({ id: 3, ...body }));
  }),
  http.get("/api/abilities/:id", ({ params }) => {
    return HttpResponse.json(makeAbility({ id: Number(params.id) }));
  }),
  http.put("/api/abilities/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeAbility({ id: Number(params.id), ...body }));
  }),
  http.delete("/api/abilities/:id", () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // Equipment templates
  http.get("/api/equipment-templates", () => {
    return HttpResponse.json([makeEquipmentTemplate({ id: 1 }), makeEquipmentTemplate({ id: 2, name: "Steel Sword" })]);
  }),
  http.get("/api/equipment-templates/:id", ({ params }) => {
    return HttpResponse.json(makeEquipmentTemplate({ id: Number(params.id) }));
  }),
  http.post("/api/equipment-templates", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeEquipmentTemplate({ id: 3, ...body }));
  }),
  http.put("/api/equipment-templates/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeEquipmentTemplate({ id: Number(params.id), ...body }));
  }),
  http.delete("/api/equipment-templates/:id", () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // NPC templates
  http.get("/api/npc-templates", () => {
    return HttpResponse.json([makeNPCTemplate({ id: 1 }), makeNPCTemplate({ id: 2, name: "Goblin" })]);
  }),
  http.post("/api/npc-templates", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeNPCTemplate({ id: 3, ...body }));
  }),
  http.get("/api/npc-templates/:id", ({ params }) => {
    return HttpResponse.json(makeNPCTemplate({ id: Number(params.id) }));
  }),

  // NPC instances
  http.get("/api/npc-instances", () => {
    return HttpResponse.json([makeNPCInstance({ id: 1 }), makeNPCInstance({ id: 2 })]);
  }),

  // Item instances
  http.get("/api/item-instances", () => {
    return HttpResponse.json([
      { id: 1, equipment_template_id: 1 },
      { id: 2, equipment_template_id: 1 },
      { id: 3, equipment_template_id: 2 },
    ]);
  }),

  // Quests
  http.get("/api/quests", () => {
    return HttpResponse.json([makeQuest({ id: 1 }), makeQuest({ id: 2, name: "Slay the Dragon" })]);
  }),
  http.get("/api/quests/lookups", () => {
    return HttpResponse.json({
      quest_types: [{ id: "main", name: "Main" }, { id: "side", name: "Side" }],
      npcs: [],
      rooms: [],
      items: [],
      effects: [],
      tags: [],
      achievements: [],
      prerequisite_quests: [],
    });
  }),
  http.get("/api/quests/:id", ({ params }) => {
    return HttpResponse.json(makeQuest({ id: Number(params.id) }));
  }),
  http.post("/api/quests", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeQuest({ id: 3, ...body }));
  }),
  http.put("/api/quests/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeQuest({ id: Number(params.id), ...body }));
  }),
  http.delete("/api/quests/:id", () => {
    return new HttpResponse(null, { status: 204 });
  }),
  http.get("/api/quests/lookups", () => {
    return HttpResponse.json({
      quest_types: [{ id: "main", name: "Main" }, { id: "side", name: "Side" }],
      npcs: [],
      rooms: [],
      items: [],
      effects: [],
      tags: [],
      achievements: [],
      prerequisite_quests: [],
    });
  }),

  // Races
  http.get("/api/races", () => {
    return HttpResponse.json([makeRace({ id: 1 }), makeRace({ id: 2, name: "Elf" })]);
  }),

  // Tags
  http.get("/api/tags", () => {
    return HttpResponse.json([makeTag({ id: 1 }), makeTag({ id: 2, name: "rare" })]);
  }),
  http.post("/api/tags", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeTag({ id: 3, ...body }));
  }),
  http.put("/api/tags/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeTag({ id: Number(params.id), ...body }));
  }),
  http.delete("/api/tags/:id", () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // Worlds
  http.get("/api/worlds/db", () => {
    return HttpResponse.json({ worlds: [makeWorld({ id: 1 }), makeWorld({ id: 2, name: "Dark Realm" })] });
  }),
  http.get("/api/worlds/:id", ({ params }) => {
    return HttpResponse.json(makeWorld({ id: Number(params.id) }));
  }),
  http.post("/api/worlds", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeWorld({ id: 3, ...body }));
  }),
  http.put("/api/worlds/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeWorld({ id: Number(params.id), ...body }));
  }),
  http.delete("/api/worlds/:id", () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // Characters
  http.get("/characters", () => {
    return HttpResponse.json([makeCharacter({ id: 1 }), makeCharacter({ id: 2, name: "Alice" })]);
  }),
  http.get("/characters/:id", ({ params }) => {
    return HttpResponse.json(makeCharacter({ id: Number(params.id) }));
  }),
  http.put("/characters/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json(makeCharacter({ id: Number(params.id), ...body }));
  }),

  // Rooms
  http.get("/api/rooms", () => {
    return HttpResponse.json([makeRoom({ id: 1 }), makeRoom({ id: 2, name: "Dungeon" })]);
  }),
];
