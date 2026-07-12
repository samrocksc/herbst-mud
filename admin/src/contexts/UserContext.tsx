/* eslint-disable react-refresh/only-export-components, functional/prefer-immutable-types, functional/no-mixed-types */
import type { ReactNode } from "react";
import { createContext, useContext, useState, useEffect, useCallback } from "react";

type UserContextType = Readonly<{
  userId: number | null;
  email: string | null;
  isAdmin: boolean | null;
  setUserData: (userId: number | null, email: string | null, isAdmin: boolean | null) => void;
  clearUserData: () => void;
}>;

const UserContext = createContext<UserContextType | undefined>(undefined);

export function UserProvider({ children }: { children: ReactNode }) {
  const [userId, setUserId] = useState<number | null>(() => {
    try {
      const stored = localStorage.getItem("userId");
      return stored ? parseInt(stored, 10) : null;
    } catch {
      return null;
    }
  });
  const [email, setEmail] = useState<string | null>(() => {
    try {
      return localStorage.getItem("email");
    } catch {
      return null;
    }
  });
  const [isAdmin, setIsAdmin] = useState<boolean | null>(() => {
    try {
      const stored = localStorage.getItem("isAdmin");
      return stored ? stored === "true" : null;
    } catch {
      return null;
    }
  });

  useEffect(() => {
    const handleStorage = (e: StorageEvent) => {
      if (e.key === "userId") {
        setUserId(e.newValue ? parseInt(e.newValue, 10) : null);
      }
      if (e.key === "email") {
        setEmail(e.newValue);
      }
      if (e.key === "isAdmin") {
        setIsAdmin(e.newValue ? e.newValue === "true" : null);
      }
    };
    window.addEventListener("storage", handleStorage);
    return () => window.removeEventListener("storage", handleStorage);
  }, []);

  const setUserData = useCallback((newUserId: number | null, newEmail: string | null, newIsAdmin: boolean | null) => {
    try {
      if (newUserId !== null) {
        localStorage.setItem("userId", String(newUserId));
        setUserId(newUserId);
      }
      if (newEmail !== null) {
        localStorage.setItem("email", newEmail);
        setEmail(newEmail);
      }
      if (newIsAdmin !== null) {
        localStorage.setItem("isAdmin", String(newIsAdmin));
        setIsAdmin(newIsAdmin);
      }
      window.dispatchEvent(new StorageEvent("storage", {
        key: "userId",
        newValue: newUserId !== null ? String(newUserId) : null,
      }));
      window.dispatchEvent(new StorageEvent("storage", {
        key: "email",
        newValue: newEmail,
      }));
      window.dispatchEvent(new StorageEvent("storage", {
        key: "isAdmin",
        newValue: newIsAdmin !== null ? String(newIsAdmin) : null,
      }));
    } catch {
      // Ignore errors
    }
  }, []);

  const clearUserData = useCallback(() => {
    try {
      localStorage.removeItem("token");
      localStorage.removeItem("userId");
      localStorage.removeItem("email");
      localStorage.removeItem("isAdmin");
      localStorage.removeItem("herbst_current_world");
      setUserId(null);
      setEmail(null);
      setIsAdmin(null);
      window.dispatchEvent(new StorageEvent("storage", { key: "userId", newValue: null }));
      window.dispatchEvent(new StorageEvent("storage", { key: "email", newValue: null }));
      window.dispatchEvent(new StorageEvent("storage", { key: "isAdmin", newValue: null }));
      window.dispatchEvent(new StorageEvent("storage", { key: "herbst_current_world", newValue: null }));
    } catch {
      // Ignore errors
    }
  }, [setUserId, setEmail, setIsAdmin]);

  return (
    <UserContext.Provider value={{ userId, email, isAdmin, setUserData, clearUserData }}>
      {children}
    </UserContext.Provider>
  );
}

export function useUser() {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error("useUser must be used within a UserProvider");
  }
  return context;
}
