import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from "react";

export type ThemeName = "dark" | "light";

type ThemeContextValue = {
  readonly theme: ThemeName;
  readonly toggle: () => void;
  readonly setTheme: (t: ThemeName) => void;
};

const STORAGE_KEY = "herbst-mud-theme";

const ThemeContext = createContext<ThemeContextValue | null>(null);

export function ThemeProvider({ children }: { readonly children: ReactNode }) {
  const [theme, setThemeState] = useState<ThemeName>(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored === "light") return "light";
    if (stored === "dark") return "dark";
    return "dark"; // default
  });

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem(STORAGE_KEY, theme);
  }, [theme]);

  const setTheme = useCallback((t: ThemeName) => {
    setThemeState(t);
  }, []);

  const toggle = useCallback(() => {
    setThemeState((prev) => (prev === "dark" ? "light" : "dark"));
  }, []);

  return (
    <ThemeContext.Provider value={{ theme, toggle, setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme(): ThemeContextValue {
  const ctx = useContext(ThemeContext);
  if (!ctx) {
    // React requires hooks to be called unconditionally, so use the callback
    // form of setState to avoid throwing in a hook body. Return a never-type
    // assertion so the compiler knows this branch is unreachable in practice.
    return (() => {
      throw new Error("useTheme must be used inside ThemeProvider");
    })();
  }
  return ctx;
}