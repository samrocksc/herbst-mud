#!/usr/bin/env python3
"""
Scaffold a test character for /player-crawler crawls.

Standing imperative (Sam, 2026-06-10): unless the crawl's primary goal
IS character creation, scaffold via this script and skip the 5-minute
React form fill. Pre-fills Ooze Surfers (world 2) with the Ooze race
and "It" gender — matches what the server validates.

Usage:
    python3 scaffold-character.py --name Crawler-1
    python3 scaffold-character.py --name Crawler-1 --room 2
    python3 scaffold-character.py --list
    python3 scaffold-character.py --cleanup Crawler-1

Never touches the `users` table.
"""
from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
import urllib.error
import urllib.request
from pathlib import Path

# ---- configuration (hardcoded for dev environment) ----
API = "http://localhost:8080"
WEB = "http://100.67.206.65:5174"
SMA_EMAIL = "sma@herbstmud.local"
SMA_PASSWORD = "sma"
SMA_USER_ID = 2  # verified at skill creation

# World 2 (Ooze Surfers) — starting room and assets
OOZE_WORLD = "2"
OOZE_ROOM = 1  # "A dank ass closet"
OOZE_RACE = "Ooze"
OOZE_GENDER = "It"

# DB connection — read from .env at repo root
# Path: .../references/scaffold-character.py → .../.agents/skills/player-crawler/references/
# parents[0] = .../references
# parents[1] = .../player-crawler
# parents[2] = .../skills
# parents[3] = .../.agents
# parents[4] = repo root (/home/sam/GitHub/herbst-mud)
ENV_FILE = Path(__file__).resolve().parents[4] / ".env"


def get_db_pw() -> str:
    if not ENV_FILE.exists():
        sys.exit(f"ERROR: {ENV_FILE} not found")
    for line in ENV_FILE.read_text().splitlines():
        if line.startswith("DB_PASSWORD="):
            return line.split("=", 1)[1].strip()
    sys.exit("ERROR: DB_PASSWORD not in .env")


def psql(sql: str) -> str:
    env = {**os.environ, "PGPASSWORD": get_db_pw()}
    r = subprocess.run(
        ["psql", "-h", "localhost", "-U", "herbst", "-d", "herbst_mud", "-t", "-A", "-c", sql],
        env=env, capture_output=True, text=True, check=True,
    )
    return r.stdout.strip()


def get_token() -> str:
    req = urllib.request.Request(
        f"{API}/users/auth",
        data=json.dumps({"email": SMA_EMAIL, "password": SMA_PASSWORD}).encode(),
        headers={"Content-Type": "application/json"},
    )
    with urllib.request.urlopen(req) as r:
        return json.loads(r.read())["token"]


def list_chars() -> list[dict]:
    out = psql(
        f"SELECT id, name, race, gender, level, current_room_id, current_world "
        f"FROM characters WHERE user_characters={SMA_USER_ID} AND is_npc=false ORDER BY id;"
    )
    rows = []
    for line in out.splitlines():
        if not line.strip():
            continue
        parts = [p.strip() for p in line.split("|")]
        rows.append({
            "id": int(parts[0]),
            "name": parts[1],
            "race": parts[2],
            "gender": parts[3],
            "level": int(parts[4]),
            "current_room_id": int(parts[5]) if parts[5] else None,
            "current_world": parts[6],
        })
    return rows


def sanitize_for_api(name: str) -> str:
    """The server validates name is letters only. Strip hyphens/digits for API mode."""
    return "".join(c for c in name if c.isalpha())


def create_via_api(name: str, room: int) -> dict:
    token = get_token()
    safe_name = sanitize_for_api(name) or "Crawler"
    body = {
        "name": safe_name,
        "race": OOZE_RACE,
        "gender": OOZE_GENDER,
        "world": OOZE_WORLD,  # server field is `world`, not `current_world`
    }
    req = urllib.request.Request(
        f"{API}/user-characters/{SMA_USER_ID}",  # no /api prefix
        data=json.dumps(body).encode(),
        headers={
            "Content-Type": "application/json",
            "Authorization": f"Bearer {token}",
        },
    )
    try:
        with urllib.request.urlopen(req) as r:
            return json.loads(r.read())
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        sys.exit(f"API ERROR {e.code}: {body}\nRequest was: {json.dumps(body, indent=2)}")


def create_via_db(name: str, room: int) -> int:
    """Direct insert. Fast, no API contract risk. Returns new character id.

    Note: bypasses the server's letter-only name validation, so hyphenated
    names like Crawler-1 work here. The web-client form will reject them,
    but this scaffolder writes directly to the DB.
    """
    # Use a unique placeholder to grab the new id from psql output
    out = psql(f"""
        WITH new_char AS (
            INSERT INTO characters
                (name, race, gender, level, hitpoints, max_hitpoints,
                 stamina, max_stamina, mana, max_mana,
                 current_world, starting_room_id, respawn_room_id,
                 current_room_id, user_characters)
            VALUES
                ('{name}', '{OOZE_RACE}', '{OOZE_GENDER}', 1, 100, 100,
                 50, 50, 25, 25,
                 '{OOZE_WORLD}', {room}, {room},
                 {room}, {SMA_USER_ID})
            RETURNING id
        )
        SELECT id FROM new_char;
    """)
    last = out.splitlines()[-1].strip() if out.strip() else ""
    try:
        return int(last)
    except ValueError:
        sys.exit(f"DB INSERT failed, psql output was:\n{out}")


def delete_chars(ids: list[int]) -> int:
    if not ids:
        return 0
    id_list = ",".join(str(i) for i in ids)
    # Guard: never delete NPCs, never touch the user row
    before = psql(f"SELECT count(*) FROM characters WHERE id IN ({id_list}) AND user_characters={SMA_USER_ID} AND is_npc=false;")
    if int(before) != len(ids):
        sys.exit(f"REFUSING: {len(ids)} ids requested, {before} match (user_characters={SMA_USER_ID}, is_npc=false). Aborting.")
    psql(f"DELETE FROM characters WHERE id IN ({id_list}) AND user_characters={SMA_USER_ID} AND is_npc=false;")
    return len(ids)


def main() -> None:
    p = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("--name", help="Character name to create (e.g. Crawler-1)")
    p.add_argument("--room", type=int, default=OOZE_ROOM, help=f"Starting room id (default {OOZE_ROOM})")
    p.add_argument("--method", choices=["api", "db"], default="db", help="Creation method (default: db — faster, bypasses server name validation)")
    p.add_argument("--list", action="store_true", help="List sma's characters and exit")
    p.add_argument("--cleanup", metavar="NAME_PATTERN", help="Delete characters whose name matches (substring). NEVER touches users.")
    args = p.parse_args()

    if args.list:
        chars = list_chars()
        if not chars:
            print("sma has 0 characters")
        else:
            print(f"sma has {len(chars)} character(s):")
            for c in chars:
                print(f"  id={c['id']:<4} name={c['name']:<20} race={c['race']:<8} gender={c['gender']:<5} "
                      f"level={c['level']} room={c['current_room_id']} world={c['current_world']}")
        return

    if args.cleanup:
        all_chars = list_chars()
        targets = [c for c in all_chars if args.cleanup.lower() in c["name"].lower()]
        if not targets:
            print(f"No characters match '{args.cleanup}'")
            return
        print(f"Will delete {len(targets)} character(s):")
        for c in targets:
            print(f"  id={c['id']} name={c['name']}")
        confirm = input("Type 'yes' to confirm: ")
        if confirm.strip().lower() != "yes":
            print("Aborted.")
            return
        n = delete_chars([c["id"] for c in targets])
        print(f"Deleted {n} character(s). User {SMA_EMAIL} untouched.")

        # Verify user still present
        user_count = psql(f"SELECT count(*) FROM users WHERE email='{SMA_EMAIL}';")
        if int(user_count) != 1:
            sys.exit(f"FATAL: user count is {user_count}, expected 1. Investigate immediately.")
        print(f"User {SMA_EMAIL} confirmed present (count={user_count}).")
        return

    if not args.name:
        p.error("--name is required (or use --list / --cleanup)")

    existing = [c for c in list_chars() if c["name"] == args.name]
    if existing:
        print(f"Character '{args.name}' already exists with id={existing[0]['id']}.")
        print("Use a different --name or run --cleanup first.")
        sys.exit(1)

    if args.method == "api":
        result = create_via_api(args.name, args.room)
        if result.get("name") != args.name:
            print(f"NOTE: requested name '{args.name}' was sanitized to '{result['name']}' (server requires letters only).")
        print(f"Created via API: id={result.get('id')} name={result.get('name')} world={result.get('currentWorld')}")
    else:
        new_id = create_via_db(args.name, args.room)
        print(f"Created via direct DB insert: id={new_id}, name={args.name}, room={args.room}, world={OOZE_WORLD}")


if __name__ == "__main__":
    main()
