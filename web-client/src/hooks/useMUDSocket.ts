import { useCallback, useEffect, useRef, useState } from "react";
import type { RoomScreenPayload } from "../lib/types";

export type ConnectionState = "idle" | "connecting" | "connected" | "error" | "closed";

export type WSLine = {
  id: number;
  text: string;
  kind: "system" | "output" | "error" | "input" | "ping";
  timestamp: number;
};

export type DebugEntry = {
  id: number;
  direction: "send" | "recv" | "state" | "error" | "local";
  label: string;
  payload?: string;
  timestamp: number;
};

let globalLineId = 0;
let globalDebugId = 0;

export function useMUDSocket() {
  const [state, setState] = useState<ConnectionState>("idle");
  const [lines, setLines] = useState<WSLine[]>([]);
  const [roomScreen, setRoomScreen] = useState<RoomScreenPayload | null>(null);
  const [debugLog, setDebugLog] = useState<DebugEntry[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  const logDebug = useCallback((direction: DebugEntry["direction"], label: string, payload?: string) => {
    const entry: DebugEntry = {
      id: ++globalDebugId,
      direction,
      label,
      payload,
      timestamp: Date.now(),
    };
    setDebugLog((prev) => [...prev.slice(-199), entry]); // keep last 200
    // Also mirror critical events to console for DevTools users
    console.log(`[DEBUG ${direction.toUpperCase()}] ${label}`, payload ?? "");
  }, []);

  const connect = useCallback((url: string) => {
    if (wsRef.current) {
      wsRef.current.close();
    }
    setState("connecting");
    setLines([]);
    setRoomScreen(null);
    logDebug("state", "connecting", url);

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      logDebug("state", "connected");
      setState("connected");
    };

    ws.onmessage = (e) => {
      logDebug("recv", "raw", String(e.data));
      try {
        const msg = JSON.parse(e.data) as { type: string; text?: string; data?: unknown; timestamp?: number };
        logDebug("recv", `type=${msg.type}`, JSON.stringify(msg).slice(0, 500));

        if (msg.type === "screen" && msg.data) {
          logDebug("state", "screen payload received");
          setRoomScreen(msg.data as RoomScreenPayload);
          return;
        }

        const kind: WSLine["kind"] =
          msg.type === "output" ? "output" :
          msg.type === "system" ? "system" :
          msg.type === "error" ? "error" :
          msg.type === "ping" ? "ping" : "output";

        setLines((prev) => [
          ...prev,
          { id: ++globalLineId, text: msg.text ?? "", kind, timestamp: msg.timestamp ?? Date.now() },
        ]);
      } catch (err) {
        const errMsg = err instanceof Error ? err.message : String(err);
        logDebug("error", "parse failed", `${errMsg} | raw: ${String(e.data).slice(0, 200)}`);
        setLines((prev) => [
          ...prev,
          { id: ++globalLineId, text: String(e.data), kind: "output", timestamp: Date.now() },
        ]);
      }
    };

    ws.onerror = (err) => {
      logDebug("error", "websocket error", String(err));
      setState("error");
    };

    ws.onclose = () => {
      logDebug("state", "closed");
      setState("closed");
      if (wsRef.current === ws) {
        wsRef.current = null;
      }
    };
  }, [logDebug]);

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

  const pushLocal = useCallback((text: string, kind: WSLine["kind"] = "output") => {
    logDebug("local", `kind=${kind}`, text);
    setLines((prev) => [
      ...prev,
      { id: ++globalLineId, text, kind, timestamp: Date.now() },
    ]);
  }, [logDebug]);

  const disconnect = useCallback(() => {
    logDebug("state", "disconnect called");
    wsRef.current?.close();
    wsRef.current = null;
    setState("idle");
  }, [logDebug]);

  // Heartbeat every 25s to keep connection alive
  useEffect(() => {
    if (state !== "connected") return;
    const id = setInterval(() => send("heartbeat", ""), 25000);
    return () => clearInterval(id);
  }, [state, send]);

  return { state, lines, roomScreen, debugLog, connect, send, disconnect, pushLocal };
}