🟣 feat(zones): zone scaffold + admin UI (v0.39.11)

Zones group rooms into named geographic/quest areas. Adds the data
layer, REST API, and admin UI in the map editor.

Server (zone scaffold + REST API)
  schema/zone.go                       new — Zone entity (string ID, room_ids
                                       []int as persistent membership list,
                                       parent_zone_id for sub-zones, color
                                       for map tinting)
  db/zone*.go                          new — ent-generated
  repository/zone_repo.go              new — ZoneRepository interface + impl
  service/zone_service.go              new — ZoneService (name conflict guard,
                                       parent-zone same-world check, hex
                                       color validation, min_level guard)
  routes/zone_routes.go                new — CRUD + GET /api/zones/:id/rooms
                                       (returns room list with `exists:
                                       false` for rooms removed from the
                                       world — the "red chip" signal)
  worldexport/                         new — full world export/import with
                                       zones; idMaps[Races|Factions] are
                                       map[int]string (Polymorphic IDs are
                                       strings; legacy int PKs preserved on
                                       the source side)
  schema/npc_template.go               +7 fields: roam_pattern (enum:
                                       static/wander/patrol/return_home),
                                       roam_zone_ids, roam_interval_seconds,
                                       roam_pause_min/max_seconds,
                                       last_moved_at, notify_on_enter
  schema/trigger.go                    +examine_weight (pre-existing
                                       half-feature completed); target_id
                                       made Optional (the DB column was
                                       character varying; polymorphic FK)
  schema/room.go                       +zone_ids []string (denormalized)
  repository/{interface,room_repo,
  faction_repo,character_faction_repo} Faction/Race string-UUID migration
                                       REVERTED — int PKs ARE unique IDs,
                                       Sam's "always use IDs" principle is
                                       satisfied. Service+repo type
                                       signatures reverted to int
  service/{interface,room_service_impl,
  container}                           UpdateRoomInput.ZoneIDs wired;
                                       roomService syncs Zone.room_ids
                                       both ways on create/update of a
                                       room's zoneIds
  service/ability_eligibility.go,
  character.go, npc_service_impl.go,
  ability_service_impl.go              Reverted to int IDs
  routes/room_crud.go                  create+update accept zoneIds
  main.go                              RegisterZoneRoutes registered

Admin (UI in the map editor)
  hooks/useZones.ts                    new — useZones() (list/create/update/
                                       delete) and useZoneRooms(zoneId)
                                       (room list with exists flag, add
                                       rooms, remove room). Uses the
                                       shared utils/apiFetch helpers.
  hooks/useRooms.ts                    +RoomInput.zoneIds
  components/map/types.ts              +Room.zoneIds
  components/map/ZonesPanel.tsx        new — accordion in MapSidebar: create
                                       form, zone rows with color swatch
                                       and delete, ZoneRoomsEditor (chip UI
                                       for room membership with
                                       search-to-add; only existing rooms
                                       are addable; ghost rooms render in
                                       red; × to remove)
  components/map/RoomZonesField.tsx    new — multi-select checkboxes for
                                       the room editor: selected zones
                                       render as removable chips with
                                       color swatches
  components/map/MapSidebar.tsx        +<ZonesPanel /> between nav and Add
                                       Room button
  components/map/RoomEditor.tsx        +RoomZonesField; form state includes
                                       zoneIds; passed to updateRoom on
                                       save

Schema migration (executed live)
  triggers.target_id                   text → bigint (NULLed 1 non-numeric
                                       value, dropped NOT NULL, converted)

Tests
  features/admin-zone-management.feature  new — 7 Gherkin scenarios:
                                       create zone, add room chip, removed
                                       room red chip, non-existent room
                                       rejected, multi-select in room
                                       editor, chip removal, zone delete
                                       cleanup

Bug fixes (v0.39.11)
  admin/src/hooks/useZones.ts          Two bugs fixed in this release:
                                       (1) Doubled-URL: apiGet/apiPost
                                       /apiPut/apiDelete prepended
                                       API_BASE, but call sites also
                                       prepended it. Result was URLs like
                                       http://100.67.206.65:5173http://100.67
                                       .206.65:5173/api/zones when window
                                       .location.origin was
                                       http://100.67.206.65:5173. Fixed
                                       by routing all calls through
                                       shared utils/apiFetch helpers which
                                       take a full URL and prepend
                                       API_BASE internally.
                                       (2) 401 Unauthorized: the hook had
                                       its own private apiGet/apiPost/etc.
                                       that didn't include the
                                       Authorization Bearer token. The
                                       backend AuthMiddleware rejected
                                       every call. Fixed by using the
                                       shared utils/apiFetch helpers,
                                       which auto-inject
                                       `Authorization: Bearer <token>`
                                       from localStorage.
  admin/src/utils/apiFetch.ts          Add "zones" to the unwrapKeys list.
                                       apiFetch auto-unwraps response
                                       envelopes of the form { "zones":
                                       [...] } (and the existing
                                       { skills | npcs | ... } keys).
                                       Fixes `filtered.map is not a
                                       function` in RoomZonesField when
                                       editing a room — the useZones hook
                                       previously returned the raw
                                       envelope object, so the field's
                                       `zones.filter(...).map(...)` call
                                       threw because the value was an
                                       object, not an array.

Verification
  go build ./...                       exits 0 in both server/ and herbst/
  server binary                        59MB, runs cleanly
  /api/zones CRUD                      201/200/200/200 across create/read/
                                       update/delete
  /api/zones/:id/rooms                 Returns room list with exists flag;
                                       ghost rooms (id 9999) marked
                                       exists=false
  PUT /api/rooms/:id zoneIds           Bi-directional sync verified: room
                                       90 in zones A and B → both zones
                                       have room_ids=[90]; remove from A
                                       → B still has 90
  Map editor admin UI                  Zones section visible in sidebar;
                                       chip UI works; ghost rooms render
                                       in red; multi-select in room
                                       editor works
  Browser smoke test                   Logged in as sma, switched world to
                                       Ooze Surfers via dropdown, opened
                                       /map?room=80&floor=1, all 47 rooms
                                       loaded on Floor 1, room list shows
                                       "Southwest Scrap Heap" (id 80) with
                                       navigable edges
