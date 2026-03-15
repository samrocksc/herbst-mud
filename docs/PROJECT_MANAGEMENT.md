# Project Management

## GitHub Projects

We use GitHub Projects for project management. The board is the single source of truth.

**Project Board:** https://github.com/users/samrocksc/projects/2

### Columns

| Column | Description |
|--------|-------------|
| **Todo** | Not started or needs rework |
| **In Progress** | Currently being implemented |
| **Blocked** | Waiting on something |
| **QA** | Ready for Raphael's review |
| **Done** | Approved and merged |

### Workflow

1. **Pick a ticket** from Todo column
2. **Assign yourself** to the issue
3. **Move to In Progress** when starting
4. **Implement & test** the feature
5. **Move to QA** when complete, create PR
6. **Raphael reviews** → Approve (→ Done) or Request Changes (→ Todo)

---

## Turtle Roles

- 🔵 **Leonardo** - PM, manages the board
- 🟣 **Donatello** - Lead Developer, implements features
- 🔴 **Raphael** - QA, reviews PRs and tests

---

## Badges

Use colored circle emojis in all commits and comments:
- 🟦 Leonardo: 🔵 Blue
- 🟣 Donatello: 🟣 Purple  
- 🔴 Raphael: 🔴 Red

---

## Quick Commands

```bash
# View board status
gh project item-list 2 --owner samrocksc --format json | jq '.items[] | {number: .content.number, status: .status}'

# List all issues
gh issue list --repo samrocksc/herbst-mud
```