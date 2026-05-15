import { useState, useEffect, useCallback, createContext, useContext } from 'react';

const STORAGE_KEY = 'herbst_current_world';

interface WorldStore {
  currentWorld: string
  setWorld: (world: string) => void
}

const WorldStoreContext = createContext<WorldStore | undefined>(undefined);

export function WorldStoreProvider({ children }: { children: React.ReactNode }) {
  const [currentWorld, setCurrentWorld] = useState<string>(() => {
    try {
      return localStorage.getItem(STORAGE_KEY) || 'default';
    } catch {
      return 'default';
    }
  });

  useEffect(() => {
    const handleStorage = (e: StorageEvent) => {
      if (e.key === STORAGE_KEY && e.newValue !== null) {
        setCurrentWorld(e.newValue);
      }
    };
    window.addEventListener('storage', handleStorage);
    return () => window.removeEventListener('storage', handleStorage);
  }, []);

  const setWorld = useCallback((world: string) => {
    try {
      localStorage.setItem(STORAGE_KEY, world);
      setCurrentWorld(world);
      window.dispatchEvent(new StorageEvent('storage', {
        key: STORAGE_KEY,
        newValue: world
      }));
    } catch {
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
  if (!context) {
    throw new Error('useWorldStore must be used within WorldStoreProvider');
  }
  return context;
}

export function getCurrentWorld(): string {
  try {
    return localStorage.getItem(STORAGE_KEY) || 'default';
  } catch {
    return 'default';
  }
}
