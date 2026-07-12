🔵 fix(player): offline PC visibility + NPC conversation overlay + stale conversation state

Three regressions fixed across server and web-client, plus skill
updates with lessons learned.

## Fix 1: Offline PCs visible in room character list (server)

Regression: non-connected player characters appeared in the room
character list because buildRoomScreen returned ALL characters via
ListByRoom with no WebSocket connection check.

Fix: Added connectedCharacterIDs() helper and isCharacterConnected()
filter. buildRoomScreen and examine now filter offline PCs. NPCs
always pass.

## Fix 2: NPC conversation overlay never opened (server)

Regression: tryTalk only returned greeting text. handleDialogChoice
existed but was dead code — handleCommand never routed the "dialog"
prefix.

Fix: tryTalk now queries DialogNode.ListByTemplate and calls
sendConversationScreen when nodes exist. handleCommand routes
"dialog <template> <node> <choice>" to handleDialogChoice.

## Fix 3: Stale conversation state reopens overlay on movement (web-client)

Regression: clicking "Leave" in the conversation overlay only set
conversationOpen=false but never cleared the conversation state
object. When the player subsequently moved to another room, the
useEffect in GameScreen saw conversation was still truthy and
reopened the overlay instead of rendering the room screen.

Root cause: onClose handler was () => setConversationOpen(false)
— it didn't clear the underlying conversation data. The useEffect
at GameScreen:266 checks if (conversation && conversation.npc_name)
first, so the stale conversation won over the new roomScreen.

Fix: Added clearConversation() to useMUDSocket hook that nulls both
conversationRef.current and the React state. Threaded through
MUDConnectionProvider context. GameScreen's onClose now calls
clearConversation() before setConversationOpen(false).

Verified via browser: talk Theodore → Leave → travel east → room
screen renders correctly, overlay stays closed. Returning and
talking again opens the overlay fresh.

## CORS config fix

.env CORS_ORIGINS was missing Tailnet IP entries for ports 5173/5174.
Added http://100.67.206.65:5173 and :5174 to the allowed origins.

## Skill updates

Both player-crawler and admin-crawler skills updated with:
- 3 new failure patterns (CORS_ORIGINS, dead code WebSocket handlers,
  tryTalk text-only fallback)
- Server restart workflow
- Commit and release workflow
- Dialog node schema reference
- Updated verification checklists

Files:
  server/routes/ws_routes.go           offline PC filter + conversation overlay routing
  web-client/src/hooks/useMUDSocket.ts  +clearConversation()
  web-client/src/context/MUDConnectionProvider.tsx  thread clearConversation
  web-client/src/components/GameScreen.tsx  onClose clears conversation state
  features/player-room-offline-pc-visibility.feature  4 Gherkin scenarios
  features/player-conversation-overlay-initial-talk.feature  4 Gherkin scenarios
  .agents/skills/player-crawler/SKILL.md  troubleshooting + commit workflow
  .agents/skills/admin-crawler/SKILL.md   patterns 13-15 + server restart + commit workflow
  docs/plans/COMMIT_MSG_v0.39.12.md   this commit message