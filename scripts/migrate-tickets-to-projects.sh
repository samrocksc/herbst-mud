#!/usr/bin/env bash
# migrate-tickets-to-projects.sh
# One-shot: read all active tickets/*.md in herbst-mud, create a GitHub
# issue per ticket, add it to samrocksc/projects/2, and stamp the issue
# number back into the markdown file.
#
# DRY-RUN BY DEFAULT. Pass --apply to actually mutate.
#
# Usage:
#   ./migrate-tickets-to-projects.sh            # dry run, prints what it would do
#   ./migrate-tickets-to-projects.sh --apply    # create issues + project items
#   ./migrate-tickets-to-projects.sh --only TICKET-TIER-001.md  # single ticket
#
# Requires: gh CLI authenticated, OR GITHUB_TOKEN with project scope.
# Configured for: /home/sam/GitHub/herbst-mud (owner=samrocksc, project=2).

set -euo pipefail

REPO_ROOT="${REPO_ROOT:-/home/sam/GitHub/herbst-mud}"
TICKETS_DIR="${REPO_ROOT}/tickets"
OWNER="${OWNER:-samrocksc}"
REPO="${REPO:-herbst-mud}"
PROJECT_NUMBER="${PROJECT_NUMBER:-2}"

# Herbst-mud field map (per github-projects skill §8). Resolved at runtime
# via `gh project field-list` so we don't hardcode the project-specific IDs.

APPLY=0
ONLY_FILE=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply) APPLY=1; shift ;;
    --only)  ONLY_FILE="$2"; shift 2 ;;
    --repo)  REPO_ROOT="$2"; TICKETS_DIR="$2/tickets"; shift 2 ;;
    --help|-h)
      sed -n '2,20p' "$0"
      exit 0
      ;;
    *) echo "unknown arg: $1" >&2; exit 2 ;;
  esac
done

# --- Auth detection ---
if ! command -v gh &>/dev/null || ! gh auth status &>/dev/null 2>&1; then
  echo "error: gh not authenticated. Run \`gh auth login\` or set GITHUB_TOKEN." >&2
  exit 1
fi

echo "Mode: $([[ $APPLY -eq 1 ]] && echo APPLY || echo DRY-RUN)"
echo "Source: $TICKETS_DIR (active tickets only — skipping done/)"
echo "Target: $OWNER/$REPO  +  project #$PROJECT_NUMBER"
echo

# --- Discover active tickets ---
if [[ -n "$ONLY_FILE" ]]; then
  TICKETS=("$ONLY_FILE")
else
  mapfile -t TICKETS < <(find "$TICKETS_DIR" -maxdepth 1 -name 'TICKET-*.md' -type f | sort)
fi
[[ ${#TICKETS[@]} -eq 0 ]] && { echo "no tickets found"; exit 0; }

# --- Resolve project field IDs once ---
echo "Resolving project field IDs..."
FIELDS_JSON=$(gh project field-list "$PROJECT_NUMBER" --owner "$OWNER" --format json 2>/dev/null \
  || echo '{"fields":[]}')

field_id() { echo "$FIELDS_JSON" | jq -r --arg n "$1" '.fields[]? | select(.name==$n) | .id' | head -1; }
option_id() {
  local fname="$1" oname="$2"
  echo "$FIELDS_JSON" | jq -r --arg f "$fname" --arg o "$oname" \
    '.fields[]? | select(.name==$f) | .options[]? | select(.name==$o) | .id' | head -1
}

STATUS_FIELD_ID=$(field_id Status)
PRIORITY_FIELD_ID=$(field_id Priority)
DIFFICULTY_FIELD_ID=$(field_id Difficulty)
TYPE_FIELD_ID=$(field_id Type)
AREA_FIELD_ID=$(field_id Area)

if [[ -z "$STATUS_FIELD_ID" ]]; then
  echo "warning: Status field not found on project. The board may not have the herbst-mud field set yet." >&2
  echo "         create the fields in the GitHub UI first (see github-projects skill §8)." >&2
  echo "         proceeding without field updates — issues will still be created + added to project." >&2
fi

# --- For each ticket, extract metadata, create issue, add to project, stamp ---
for f in "${TICKETS[@]}"; do
  base=$(basename "$f" .md)

  # Title: first H1 (e.g. "# TICKET TIER-001: ...")
  TITLE=$(awk '/^# /{sub(/^# /,""); print; exit}' "$f")
  [[ -z "$TITLE" ]] && TITLE="$base"

  # Difficulty: frontmatter `difficulty: <value>`
  DIFFICULTY=$(awk '/^difficulty:/{print $2; exit}' "$f" | tr -d '[:space:]')
  [[ -z "$DIFFICULTY" ]] && DIFFICULTY="med"

  # Status heuristic for the project: everything in tickets/ (not done/) is
  # either Backlog or Ready. Treat the file's mtime as a "last touched" hint.
  if [[ $(find "$f" -mtime -3 2>/dev/null) ]]; then
    INITIAL_STATUS="Ready"
  else
    INITIAL_STATUS="Backlog"
  fi

  # Body: everything after the first H2 (preserves Problem / Findings / etc.)
  BODY=$(awk '/^## /{p=1} p' "$f")
  [[ -z "$BODY" ]] && BODY="$TITLE"

  # Source link footer
  FOOTER=$'\n\n---\n_Source: [tickets/'$base'.md]('$base'.md) in `'$REPO'`. Auto-migrated by `migrate-tickets-to-projects.sh`._'
  FULL_BODY="${BODY}${FOOTER}"

  # Labels from filename
  AREA=$(echo "$base" | sed -E 's/^TICKET-([A-Z]+)-.*/\L\1/' | tr -d '0-9-')
  [[ -z "$AREA" || "$AREA" == "$base" ]] && AREA="general"

  echo "─────────────────────────────────────────────────"
  echo "Ticket:  $base"
  echo "Title:   $TITLE"
  echo "Diff:    $DIFFICULTY"
  echo "Status:  $INITIAL_STATUS"
  echo "Area:    $AREA"
  echo "Source:  $f"

  if [[ $APPLY -eq 0 ]]; then
    echo "Action:  [dry-run] would create issue + add to project"
    continue
  fi

  # 1. Create the issue
  ISSUE_URL=$(gh issue create \
    --repo "$OWNER/$REPO" \
    --title "$TITLE" \
    --body "$FULL_BODY" \
    --label "$AREA" 2>/dev/null || echo "")

  if [[ -z "$ISSUE_URL" ]]; then
    echo "  [error] issue create failed for $base — skipping" >&2
    continue
  fi
  ISSUE_NUM=$(echo "$ISSUE_URL" | grep -oE '[0-9]+$')
  echo "Issue:   $ISSUE_URL  (#$ISSUE_NUM)"

  # 2. Add to project
  ITEM_OUT=$(gh project item-add "$PROJECT_NUMBER" --owner "$OWNER" --url "$ISSUE_URL" 2>&1)
  ITEM_ID=$(echo "$ITEM_OUT" | grep -oE '[0-9a-f]{8,}' | head -1)
  if [[ -z "$ITEM_ID" ]]; then
    # Fallback: parse from the table output
    ITEM_ID=$(gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json 2>/dev/null \
      | jq -r --arg n "$ISSUE_NUM" '.items[]? | select(.content.number == ($n|tonumber)) | .id' | head -1)
  fi
  echo "Item:    ${ITEM_ID:-<not resolved>}"

  # 3. Set custom fields (Status, Difficulty, Priority)
  if [[ -n "$ITEM_ID" && -n "$STATUS_FIELD_ID" ]]; then
    STATUS_OPT_ID=$(option_id Status "$INITIAL_STATUS")
    if [[ -n "$STATUS_OPT_ID" ]]; then
      gh project item-edit --id "$ITEM_ID" --field-id "$STATUS_FIELD_ID" \
        --single-select-option-id "$STATUS_OPT_ID" >/dev/null 2>&1 \
        || echo "  [warn] could not set Status=$INITIAL_STATUS"
      echo "  set:  Status = $INITIAL_STATUS"
    fi
    if [[ -n "$DIFFICULTY_FIELD_ID" ]]; then
      DIFF_OPT_ID=$(option_id Difficulty "$DIFFICULTY")
      if [[ -n "$DIFF_OPT_ID" ]]; then
        gh project item-edit --id "$ITEM_ID" --field-id "$DIFFICULTY_FIELD_ID" \
          --single-select-option-id "$DIFF_OPT_ID" >/dev/null 2>&1 \
          || echo "  [warn] could not set Difficulty=$DIFFICULTY"
        echo "  set:  Difficulty = $DIFFICULTY"
      fi
    fi
  fi

  # 4. Stamp the issue number into the markdown file as a banner
  BANNER="<!-- Migrated to GitHub Projects: $ISSUE_URL on $(date -I) -->\n"
  if head -1 "$f" | grep -q "Migrated to GitHub Projects"; then
    echo "  skip:  $base already has migration banner"
  else
    tmp=$(mktemp)
    printf '%s' "$BANNER" | cat - "$f" > "$tmp" && mv "$tmp" "$f"
    echo "  stamp: prepended migration banner to $base"
  fi

  # 5. Rate-limit courtesy
  sleep 1
done

echo
echo "Done."
[[ $APPLY -eq 0 ]] && echo "Re-run with --apply to actually create the issues."
