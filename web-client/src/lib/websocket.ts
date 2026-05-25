import { useState, useCallback } from "react";
import type { ServerMessage, ClientMessage, OutputLine } from "../types";

export type ConnectionState = "idle" | "connecting" | "connected" | "error" | "closed";

type UseWebSocketReturn = {
  state: ConnectionState;
  output: OutputLine[];
  lastScreen: ServerMessage | null;
  lastEvent: ServerMessage | null;
  connect: (url: string) => void;
  send: (msg: ClientMessage) => void;
  disconnect: () => void;
};

export const useWebSocket = (): UseWebSocketReturn => {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [state, setState] = useState<ConnectionState>("idle");
  const [output, setOutput] = useState<OutputLine[]>([]);
  const [lastScreen, setLastScreen] = useState<ServerMessage | null>(null);
  const [lastEvent, setLastEvent] = useState<ServerMessage | null>(null);

  const connect = useCallback((url: string) => {
    setState("connecting");
    const socket = new WebSocket(url);

    socket.onopen = () => {
      setState("connected");
    };

    socket.onmessage = (e) => {
      try {
        const msg: ServerMessage = JSON.parse(e.data);
        if (msg.type === "output") {
          setOutput((prev) => [...prev, ...msg.payload.lines]);
        } else if (msg.type === "screen") {
          setLastScreen(msg);
        } else if (msg.type === "event") {
          setLastEvent(msg);
        }
      } catch {
        // Non-JSON message — append as raw text line
        setOutput((prev) => [
          ...prev,
          { text: e.data, style: "default", timestamp: Date.now() },
        ]);
      }
    };

    socket.onerror = () => {
      setState("error");
    };

    socket.onclose = () => {
      setState("closed");
      setWs(null);
    };

    setWs(socket);
  }, []);

  const send = useCallback((msg: ClientMessage) => {
    ws?.send(JSON.stringify(msg));
  }, [ws]);

  const disconnect = useCallback(() => {
    ws?.close();
    setWs(null);
    setState("idle");
  }, [ws]);

  return { state, output, lastScreen, lastEvent, connect, send, disconnect };
};