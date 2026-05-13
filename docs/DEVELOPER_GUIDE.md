# Herbst MUD | Developer Portal

> 🔵 **Master Knowledge Base & Project Index**
> This is the single source of truth for operating, developing, and extending the Herbst MUD engine.

---

## 🗺️ Project Map

### 🚀 Quick Start & Operations
*   **[Getting Started](./DEVELOPER_GUIDE.md)**: Local setup, prerequisites, and environment variables.
*   **[Deployment](./ops/DEPLOYMENT.md)**: Digital Ocean specs and deployment pipeline.
*   **[Infrastructure](./ops/SSH_SERVER.md)**: SSH server implementation and port mapping.
*   **[Dependencies](./ops/DEPENDENCIES.md)**: Project-wide library requirements.

### 🏗️ Architecture & Specifications
*   **[Technical Specs](./specs/)**: All RFCs (Request for Comments) and deep-dives.
    *   `RFC-001` through `RFC-009` covering Web Arch, Effects, Quests, Dialog Trees, LLM NPCs, and the Messaging system.
    *   **[Character System Deep-dive](./specs/CHARACTER-SYSTEM-DEEPDIVE.md)**: Core logic for attributes and progression.
*   **[Feature Set](./specs/FEATURES.md)**: List of implemented and planned features.

### 🛠️ Development Standards
*   **[Go Standards](./docs/GO_BEST_PRACTICES.md)**: Minimalist Go, Gin patterns, and the 100-line file rule.
*   **[Frontend Standards](./docs/REACT.md)**: TanStack Router, React Query, and Tailwind CSS usage.
*   **[Testing Protocol](./guides/TESTING.md)**: TDD approach and Gherkin feature testing.
*   **[Contribution Workflow](./docs/GITHUB_ETIQUETTE.md)**: PR standards and commit etiquette.

### 📖 User & Admin Manuals
*   **[Admin Panel Guide](./guides/ADMINISTRATION.md)**: How to manage the world via the web UI.
*   **[Player's Guide](./guides/PLAYER_GUIDE.md)**: Game mechanics and user interface instructions.

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
