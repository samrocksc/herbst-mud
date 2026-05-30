# Herbst MUD | Developer Portal

> 🔵 **Master Knowledge Base & Project Index**
> This is the single source of truth for operating, developing, and extending the Herbst MUD engine.

---

## 🗺️ Project Map

### 🚀 Quick Start & Operations
*   **[Getting Started](./DEVELOPER_GUIDE.md)**: Local setup, prerequisites, and environment variables.
*   **[Installation & Upgrade](./OPERATIONS/INSTALL.md)**: Docker Compose deployment, migrations, and upgrades.
*   **[Operations Manual](./OPERATIONS/INDEX.md)**: Running the stack, logs, health checks, and backups.

### 🏗️ Architecture & Specifications
*   **[Technical Specs](./SPECS/)**: All RFCs and deep-dives.
    *   **[Character System Deep-dive](./SPECS/CHARACTER-SYSTEM-DEEPDIVE.md)**: Core logic for attributes and progression.
*   **[Feature Set](./SPECS/FEATURES.md)**: List of implemented and planned features.
*   **[World System Overview](./SPECS/WORLD-SYSTEM.md)**: Multi-world architecture, `world_id` filtering, and whitelist-based access control.
*   **[Effects System](./SPECS/effects-system.md)**: Ability effects, hooks, and active effects system.

### 🛠️ Development Standards
*   **[Developer Guide](./DEVELOPER-GUIDE/INDEX.md)**: Build, test, and codebase conventions.

### 📖 User & Admin Manuals
*   **[Admin Panel Guide](./ADMIN-GUIDE/)**: World management via the web UI (coming soon).

---

## 🛠️ Core Tech Stack

| Layer | Technology | Role |
| :--- | :--- | :--- |
| **Backend** | Go + Gin | REST API (`server/`) |
| **Game Engine** | Go + BubbleTea | SSH TUI Server (`herbst/`) |
| **Database** | PostgreSQL + Ent | Schema-driven ORM |
| **Admin UI** | React + Vite | Management SPA (`admin/`) |
| **Routing** | TanStack Router | File-based frontend routing |

## 📜 The Golden Rules

1.  **The 100-Line Rule**: Files should generally not exceed 100 lines. If they do, split them by sub-domain (e.g., `routes/character/crud.go`, `routes/character/combat.go`).
2.  **TDD First**: No production code without a failing test first.
3.  **Minimalist Design**: Prefer functional-lite over heavy OOP.
4.  **Genre Agnostic**: Do not hardcode campaign/genre specifics; make them configurable via game settings.

---
*Created by Leonardo 🐢 | la testa e il cuore del progetto.*
