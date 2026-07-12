fix(admin): clean up the path to immutability in the 3 map anchor files

TICKET-ARCH-001: the original ticket's "remove the off override" prescription
was a no-op (no such override existed) and the "remove disable from 3 files"
prescription covered <10% of the actual work. Revised approach: surgically
remove the blanket functional eslint disables from the 3 named anchor files,
convert data shapes to Readonly where it doesn't break the algorithm, and
add targeted per-line disables with explicit reasoning where the imperative
shape is intentional (hot numerical loops, React useRef, closure-local
dedup Sets).

Changes:
- install immer (no functional rule change yet — Phase 1 of the ticket)
- useMapState.ts:
  - convert reactFlowNodes/Edges builders from imperative for-push to filter/flatMap
  - refactor updateSearchParams to immutable object construction
  - simplify chained Array.from(new Set(...)) to [...new Set(values)]
  - remove dead imports (GRID, updateRoomAsync)
  - per-line disables only: React useRef re-entry guards, closure-local
    dedup Set, useCallback dep warnings (pre-existing)
- ExitLineRenderer.tsx:
  - ExitLinesProps.rooms: Room[] -> ReadonlyArray<Room>
  - ExitLinesProps.nodePositions: Map -> ReadonlyMap
  - Segment type -> Readonly<>
  - resolveOverlaps: positions param -> ReadonlyMap
  - replace `new Map() + for-set` with `new Map(input)`
  - per-line disables only: hot relaxation loop (Immer would be O(n^2)
    perf downgrade), segment builder (early-exit + dedup needs imperative)
- RoomNode.tsx (subagent-driven):
  - RoomNodeProps -> Readonly<> with ReadonlyArray<>
  - removed unused props from destructure (pos, zoom — not used in JSX)
  - per-line disable only: no-mixed-types (props legitimately mix data + handlers)

Verification:
- npx tsc --noEmit: clean
- npx eslint on the 3 files: 0 errors (2 pre-existing react-hooks warnings
  on useMapState.ts, 0 warnings on the other two)
- npx vitest run src/components/map/: 5 files, 62 tests, all pass
- dev server: HTTP 200 on / and /map

Out of scope: 47 other files in admin/src/ with blanket functional eslint
disables. Those should be done file-by-file in follow-up tickets per the
ticket's revised dispatch strategy.

LEONARDO default — pausing for Sam to review and invoke git commit. Per
the standing commit rule (2026-06-30), do not commit without Sam's
explicit invocation.
