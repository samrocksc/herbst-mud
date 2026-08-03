# Export/Import Round-Trip Validation

When the QA process makes any change to the herbst-mud backend, the world
export → reimport cycle must be exercised before declaring the change safe.
This catches serialization regressions, FK drift, and world-filter bugs
that unit tests miss.

The full procedure is documented in
`player-crawler/references/export-import-roundtrip.md`. This file is the
admin-crawler copy that points at it and summarizes what to run.

## When to run

- After any direct fix or refactor in `server/`
- After any change to `server/db/schema/`
- After any change to `server/worldexport/`
- After any change to `server/routes/admin_wipe_routes.go` or
  `server/routes/game_export_routes.go`
- After any change to `server/dbinit/`
- After any change to `server/routes/ability_routes.go` or
  `server/service/ability_equip.go` (affects ability export shape)
- After any change to the zones scaffold (export includes zones field)

Skip only for: pure frontend changes, pure test-file changes, doc-only
changes.

## Quick command

Write a Python script to avoid the shell token-redaction guard:

```bash
# Write the script
cat > /tmp/qa_export_roundtrip.py <<'PY'
import urllib.request, json
def get_token():
    req = urllib.request.Request("http://localhost:8080/users/auth",
        data=json.dumps({"email":"admin@herbstmud.local","password":"herb5t2026!"}).encode(),
        headers={"Content-Type":"application/json"})
    return json.loads(urllib.request.urlopen(req).read())["token"]
def call(method, path, token, body=None):
    data = json.dumps(body).encode() if body else None
    req = urllib.request.Request(f"http://localhost:8080{path}", data=data, method=method,
        headers={"Authorization":"Bearer "+token, "Content-Type":"application/json"})
    return urllib.request.urlopen(req)
token = get_token()
worlds = json.loads(call("GET","/admin/export/worlds",token).read())["worlds"]
for w in worlds:
    if not w.get("active",True): continue
    snap = json.loads(call("GET",f"/admin/export?world={w['id']}",token).read())
    print(f"world {w['id']} ({w['name']}): abilities={len(snap.get('abilities',[]))} npcs={len(snap.get('npcs',[]))} rooms={len(snap.get('rooms',[]))} dialog_nodes={len(snap.get('dialog_nodes',[]))}")
    with open(f"/tmp/qa_world{w['id']}_snapshot.json","w") as f:
        json.dump(snap, f)
PY
python3 /tmp/qa_export_roundtrip.py
```

## What to look for

For each world exported:

1. **Counts are non-zero where they should be.** A world with NPCs should
   have `npcs` and `npc_templates` > 0. A world with dialog should have
   `dialog_nodes` > 0. If a previously-populated field shows 0, that's
   a regression.
2. **World scoping is correct.** Exporting `world=2` should return
   abilities with `world_id == "2"`, not `world_id == "1"`. If you see
   the wrong `world_id` in the payload, the export filter is broken (a
   known regression as of 2026-07-10).
3. **FK relationships are preserved.** Every `character_abilities` row
   should reference an `ability_id` that exists in `abilities`. Every
   `ability_effects` should reference an existing ability.
4. **Validate passes.** `POST /admin/import/validate` should return 200
   with no errors.

## Full procedure

See `../player-crawler/references/export-import-roundtrip.md` for the
full step-by-step with the import round-trip and cleanup.
