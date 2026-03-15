# AGENTS.md - Project Agent Instructions

You are a programming assistant for building a MUD(Multi-User Dungeon). Your job is to build solid testable code.

**IMPORTANT** the goal of this isn't to create a single mud, but to create an engine that can run multiple muds and multiple story lines.

- the administrative backend is a vite/tanstack router project in the `admin/` folder.
- It should use GoLang's Crush
- utilize `docs/` markdown for knowing how to interface with the project(e.g. `docs/TESTING.md` for test writing).
- It should have a `admin` folder, which should contain a vite react project.
- Utilize `features/` for seeing what features have been implemented. Ensure that the features are written using proper Gherkin format

## Onboarding Checklist

When starting work on this project, always:

1. **Read `AGENT_KNOWLEDGE.md`** - Contains project architecture and notes
2. **Read `docs/*.md`** - Each of the files in the docs will allow us to understand the project, and how it worksl
3. **Check `features/`** - See what features are implemented/in-progress
4. **Review existing code structure** before making changes

> ⚠️ **IMPORTANT:** Always create or update `AGENT_KNOWLEDGE.md` when you learn something new about the project. This helps future agents onboard quickly.
