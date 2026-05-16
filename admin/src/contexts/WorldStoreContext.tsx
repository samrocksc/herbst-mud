/* eslint-disable react-refresh/only-export-components, functional/prefer-immutable-types, functional/no-mixed-types */
import type { ReactNode } from "react";
import { useState, useEffect, useCallback, createContext, useContext } from "react";

const STORAGE_KEY = "herbst_current_world";

type WorldStoreContextType = Readonly<{
  currentWorld: string;
  setWorld: (world: string) => void;
}>

const WorldStoreContext = createContext<WorldStoreContextType | undefined>(undefined);

export function WorldStoreProvider({ children }: { children: ReactNode }) {
  const [currentWorld, setCurrentWorld] = useState<string>(() => {
    try {
      return localStorage.getItem(STORAGE_KEY) || "default";
    } catch {
      return "default";
    }
  });

  useEffect(() => {
    const handleStorage = (e: StorageEvent) => {
      if (e.key === STORAGE_KEY && e.newValue !== null) {
        setCurrentWorld(e.newValue);
      }
    };
    window.addEventListener("storage", handleStorage);
    return () => window.removeEventListener("storage", handleStorage);
  }, []);

  const setWorld = useCallback((world: string) => {
    try {
      localStorage.setItem(STORAGE_KEY, world);
      setCurrentWorld(world);
      window.dispatchEvent(new StorageEvent("storage", {
        key: STORAGE_KEY,
        newValue: world,
      }));
    } catch {
      // Ignore errors
    }
  }, []);

  return (
    <WorldStoreContext.Provider value={{ currentWorld, setWorld }}>
      {children}
    </WorldStoreContext.Provider>
  );
}

export function useWorldStore() {
  const context = useContext(WorldStoreContext);
  if (context === undefined) {
    return { currentWorld: "default", setWorld: (_w: string) => {} };
  }
  return context;
}
