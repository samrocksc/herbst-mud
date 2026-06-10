# COMPLETION_REPORT_COMBAT_QA_B.md

## Sub-ticket B: Player client — WebSocket auto-reconnect

**Finding #5:** WebSocket drops to empty page on tab idle

---

### What Was Done

Added auto-reconnect with exponential backoff to `web-client/src/hooks/useMUDSocket.ts`.

**Reconnect logic:**
- Triggered on `onclose` event (server idle timeout)
- NOT triggered on `onerror` (fatal — server may be down, let user decide)
- Delays: 1s → 2s → 4s → 8s → 16s → 30s (cap)
- `shouldReconnect` ref gates reconnect; set to `false` on explicit `disconnect()` call
- `urlRef` preserves the WebSocket URL across reconnects
- `connectRef` breaks circular dependency (onclose timer can't reference `connect` directly in useCallback deps)

**On explicit logout:**
- `disconnect()` sets `shouldReconnect = false`
- Clears any pending reconnect timer
- Closes socket — no background reconnect

---

### Files Modified

- `web-client/src/hooks/useMUDSocket.ts`
  - Added refs: `urlRef`, `shouldReconnect`, `reconnectAttempt`, `reconnectTimer`, `connectRef`
  - Updated `connect()` to reset reconnect state and store URL
  - Updated `onclose` handler to schedule reconnect with backoff
  - Updated `disconnect()` to disable reconnect

---

### Verification

```
cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit  # PASS (no output)
cd /home/sam/GitHub/herbst-mud && make build-all              # PASS (SSH + Web binaries)
```

---

### Not Changed

- WebSocket message protocol
- State management library
- Login/logout flow
- Any server-side code