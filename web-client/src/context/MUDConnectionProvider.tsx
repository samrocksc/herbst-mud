import React, { createContext, useContext } from "react";
import { useMUDSocket } from "../hooks/useMUDSocket";
import type { ConversationState } from "../hooks/useMUDSocket";
import type { ShopSession } from "../lib/types";

interface MUDConnectionContextType {
  state: string;
  lines: any[];
  roomScreen: any;
  conversation: ConversationState | null;
  clearConversation: () => void;
  shopData: ShopSession | null;
  clearShopData: () => void;
  vitals: any;
  debugLog: any[];
  connect: (url: string) => void;
  send: (type: "command" | "heartbeat", payload: string) => void;
  disconnect: () => void;
  pushLocal: (text: string, kind?: any) => void;
}

const MUDContext = createContext<MUDConnectionContextType | null>(null);

export const MUDConnectionProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const socket = useMUDSocket();

  return (
    <MUDContext.Provider value={{
      state: socket.state,
      lines: socket.lines,
      roomScreen: socket.roomScreen,
      conversation: socket.conversation,
      clearConversation: socket.clearConversation,
      shopData: socket.shopData,
      clearShopData: socket.clearShopData,
      vitals: socket.vitals,
      debugLog: socket.debugLog,
      connect: socket.connect,
      send: socket.send,
      disconnect: socket.disconnect,
      pushLocal: socket.pushLocal,
    }}>
      {children}
    </MUDContext.Provider>
  );
};

export const useMUDConnection = () => {
  const context = useContext(MUDContext);
  if (!context) throw new Error("useMUDConnection must be used within MUDConnectionProvider");
  return context;
};
