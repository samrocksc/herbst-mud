import { useCallback, useEffect, useRef, useState } from "react";
import type { RoomScreenPayload, VitalsPayload } from "../lib/types";

export type ConnectionState = "idle" | "connecting" | "connected" | "error" | "closed";

export type WSLine = {
  id: number;
  text: string;
  kind: "system" | "output" | "error" | "input" | "ping" | "notification";
  timestamp: number;
};

export type DebugEntry = {
  id: number;
  direction: "send" | "recv" | "state" | "error" | "local";
  label: string;
  payload?: string;
  timestamp: number;
};

export type ConversationState = {
  npc_name: string;
  npc_template_id: string;
  nodes: Record<string, import("../lib/types").DialogNode>;
  current_node_id: string;
};

let globalLineId = 0;
let globalDebugId = 0;

// Global cache for hook functions - keyed by component instance identity
interface HookFunctions {
  logDebug: (direction: DebugEntry["direction"], label: string, payload?: string) => void;
  connect: (url: string) => void;
  send: (type: "command" | "heartbeat", payload: string) => void;
  pushLocal: (text: string, kind?: WSLine["kind"]) => void;
  disconnect: () => void;
}

const hookFunctionsCache = new WeakMap<object, HookFunctions>();

// Ref storage - separate cache to hold refs that need to be shared with React
interface HookRefs {
  stateRef: { current: ConnectionState };
  linesRef: { current: WSLine[] };
  roomScreenRef: { current: RoomScreenPayload | null };
  conversationRef: { current: ConversationState | null };
  vitalsRef: { current: VitalsPayload | null };
  debugLogRef: { current: DebugEntry[] };
  wsRef: { current: WebSocket | null };
  urlRef: { current: string | null };
  reconnectTimerRef: { current: ReturnType<typeof setTimeout> | null };
  shouldReconnectRef: { current: boolean };
  reconnectAttemptRef: { current: number };
}
const hookRefsCache = new WeakMap<object, HookRefs>();

export function useMUDSocket() {
  const [state, setState] = useState<ConnectionState>("idle");
  const [lines, setLines] = useState<WSLine[]>([]);
  const [roomScreen, setRoomScreen] = useState<RoomScreenPayload | null>(null);
  const [conversation, setConversation] = useState<ConversationState | null>(null);

  const clearConversation = useCallback(() => {
    conversationRef.current = null;
    setConversation(null);
  }, []);
  const [vitals, setVitals] = useState<VitalsPayload | null>(null);
  const [debugLog, setDebugLog] = useState<DebugEntry[]>([]);

  // Create a unique identity for this hook instance - always first
  const instanceIdentityRef = useRef({});

  // Get or create refs for this component instance - always run first
  let refs = hookRefsCache.get(instanceIdentityRef.current);
  if (!refs) {
    refs = {
      stateRef: { current: "idle" },
      linesRef: { current: [] },
      roomScreenRef: { current: null },
      conversationRef: { current: null },
      vitalsRef: { current: null },
      debugLogRef: { current: [] },
      wsRef: { current: null },
      urlRef: { current: null },
      reconnectTimerRef: { current: null },
      shouldReconnectRef: { current: true },
      reconnectAttemptRef: { current: 0 },
    };
    hookRefsCache.set(instanceIdentityRef.current, refs);
  }

  // Always run hooks in consistent order, regardless of cache status
  // These refs are used by the functions below
  const stateRef = refs.stateRef;
  const linesRef = refs.linesRef;
  const roomScreenRef = refs.roomScreenRef;
  const conversationRef = refs.conversationRef;
  const vitalsRef = refs.vitalsRef;
  const debugLogRef = refs.debugLogRef;
  const wsRef = refs.wsRef;
  const urlRef = refs.urlRef;
  const reconnectTimerRef = refs.reconnectTimerRef;
  const shouldReconnectRef = refs.shouldReconnectRef;
  const reconnectAttemptRef = refs.reconnectAttemptRef;

  // logDebug - wrapped in useCallback with empty deps
  const logDebug = useCallback((direction: DebugEntry["direction"], label: string, payload?: string) => {
    const entry: DebugEntry = {
      id: ++globalDebugId,
      direction,
      label,
      payload,
      timestamp: Date.now(),
    };
    debugLogRef.current = [...debugLogRef.current.slice(-199), entry];
    console.log(`[DEBUG ${direction.toUpperCase()}] ${label}`, payload ?? "");
    setDebugLog(debugLogRef.current);
  }, []);

  // disconnect - wrapped in useCallback with empty deps
  const disconnect = useCallback(() => {
    logDebug("state", "disconnect called");
    shouldReconnectRef.current = false;
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current);
      reconnectTimerRef.current = null;
    }
    wsRef.current?.close();
    wsRef.current = null;
    stateRef.current = "idle";
    setState("idle");
  }, [logDebug]);

  // connect - wrapped in useCallback with empty deps
  const connect = useCallback((url: string) => {
    // Guard: if URL hasn't changed and we're already connecting/connected, skip
    if (urlRef.current === url && (stateRef.current === "connecting" || stateRef.current === "connected")) {
      logDebug("state", "connect skipped - already in flight or connected", url);
      return;
    }

    // Guard: close any existing socket BEFORE creating a new one
    if (wsRef.current) {
      const old = wsRef.current;
      wsRef.current = null;
      // Temporarily suppress reconnect so the close doesn't trigger a reconnect loop
      shouldReconnectRef.current = false;
      old.close();
    }

    stateRef.current = "connecting";
    setState("connecting");
    urlRef.current = url;
    shouldReconnectRef.current = true;
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current);
      reconnectTimerRef.current = null;
    }
    linesRef.current = [];
    roomScreenRef.current = null;
    conversationRef.current = null;
    logDebug("state", "connecting", url);

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      logDebug("state", "connected");
      stateRef.current = "connected";
      setState("connected");
      reconnectAttemptRef.current = 0;
    };

    ws.onmessage = (e) => {
      logDebug("recv", "raw", String(e.data));
      try {
        const msg = JSON.parse(e.data) as { type: string; text?: string; data?: unknown; timestamp?: number };
        logDebug("recv", `type=${msg.type}`, JSON.stringify(msg).slice(0, 500));

        if (msg.type === "screen" && msg.data) {
          const data = msg.data as { view_type: string };
          if (data.view_type === "conversation") {
            logDebug("state", "conversation screen payload received");
            conversationRef.current = msg.data as ConversationState;
            setConversation(msg.data as ConversationState);
          } else {
            logDebug("state", "room screen payload received");
            roomScreenRef.current = msg.data as RoomScreenPayload;
            setRoomScreen(msg.data as RoomScreenPayload);
          }
          return;
        }

        if (msg.type === "vitals" && msg.data) {
          logDebug("state", "vitals payload received");
          vitalsRef.current = msg.data as VitalsPayload;
          setVitals(msg.data as VitalsPayload);
          return;
        }

        const kind: WSLine["kind"] =
          msg.type === "output" ? "output" :
          msg.type === "system" ? "system" :
          msg.type === "error" ? "error" :
          msg.type === "ping" ? "ping" :
          msg.type === "notification" ? "notification" : "output";

        const newLine: WSLine = { id: ++globalLineId, text: msg.text ?? "", kind, timestamp: msg.timestamp ?? Date.now() };
        linesRef.current = [...linesRef.current, newLine];
        setLines(linesRef.current);
      } catch (err) {
        const errMsg = err instanceof Error ? err.message : String(err);
        logDebug("error", "parse failed", `${errMsg} | raw: ${String(e.data).slice(0, 200)}`);
        const newLine: WSLine = { id: ++globalLineId, text: String(e.data), kind: "output" as const, timestamp: Date.now() };
        linesRef.current = [...linesRef.current, newLine];
        setLines(linesRef.current);
      }
    };

    ws.onerror = (err) => {
      logDebug("error", "websocket error", String(err));
      stateRef.current = "error";
      setState("error");
    };

    ws.onclose = () => {
      logDebug("state", "closed");
      stateRef.current = "closed";
      setState("closed");
      if (wsRef.current === ws) {
        wsRef.current = null;
      }
      if (shouldReconnectRef.current && urlRef.current) {
        reconnectAttemptRef.current += 1;
        const delay = Math.min(1000 * Math.pow(2, reconnectAttemptRef.current - 1), 30000);
        logDebug("state", `reconnecting in ${delay}ms (attempt ${reconnectAttemptRef.current})`);
        const savedUrl = urlRef.current;
        reconnectTimerRef.current = setTimeout(() => {
          reconnectTimerRef.current = null;
          connect(savedUrl);
        }, delay);
      }
    };
  }, [logDebug]);

  // send - wrapped in useCallback with empty deps
  const send = useCallback((type: "command" | "heartbeat", payload: string) => {
    const ws = wsRef.current;
    if (ws?.readyState === WebSocket.OPEN) {
      const json = JSON.stringify({ type, payload });
      ws.send(json);
      logDebug("send", `${type}${payload ? ` | ${payload}` : ""}`, json);
    } else {
      logDebug("error", "send failed — not open", `readyState=${ws?.readyState} type=${type} payload=${payload}`);
    }
  }, [logDebug]);

  // pushLocal - wrapped in useCallback with empty deps
  const pushLocal = useCallback((text: string, kind: WSLine["kind"] = "output") => {
    logDebug("local", `kind=${kind}`, text);
    linesRef.current = [...linesRef.current, { id: ++globalLineId, text, kind, timestamp: Date.now() }];
    setLines(linesRef.current);
  }, [logDebug]);

  // Get cached functions or create them (only happens once per instance)
  let cached = hookFunctionsCache.get(instanceIdentityRef.current);
  if (!cached) {
    cached = { logDebug, connect, send, pushLocal, disconnect };
    hookFunctionsCache.set(instanceIdentityRef.current, cached);
  }

  // Refs for state updates - always run after the hooks above
  const setStateRef = useRef(setState);
  const setLinesRef = useRef(setLines);
  const setRoomScreenRef = useRef(setRoomScreen);
  const setConversationRef = useRef(setConversation);
  const setVitalsRef = useRef(setVitals);
  const setDebugLogRef = useRef(setDebugLog);

  useEffect(() => {
    setStateRef.current = setState;
    setLinesRef.current = setLines;
    setRoomScreenRef.current = setRoomScreen;
    setConversationRef.current = setConversation;
    setVitalsRef.current = setVitals;
    setDebugLogRef.current = setDebugLog;
  }, [setState, setLines, setRoomScreen, setConversation, setVitals, setDebugLog]);

  // Heartbeat every 25s to keep connection alive
  useEffect(() => {
    if (state !== "connected") return;
    const id = setInterval(() => cached.send("heartbeat", ""), 25000);
    return () => clearInterval(id);
  }, [state, cached.send]);

  // Sync internal state refs to React state for rendering
  useEffect(() => {
    setDebugLogRef.current(refs.debugLogRef.current);
    setLinesRef.current(refs.linesRef.current);
    if (refs.roomScreenRef.current) setRoomScreenRef.current(refs.roomScreenRef.current);
    if (refs.conversationRef.current) setConversationRef.current(refs.conversationRef.current);
    if (refs.vitalsRef.current) setVitalsRef.current(refs.vitalsRef.current);
    setStateRef.current(refs.stateRef.current);
  }, []);

  return {
    state,
    lines,
    roomScreen,
    conversation,
    clearConversation,
    vitals,
    debugLog,
    connect: cached.connect,
    send: cached.send,
    disconnect: cached.disconnect,
    pushLocal: cached.pushLocal,
  };
}
