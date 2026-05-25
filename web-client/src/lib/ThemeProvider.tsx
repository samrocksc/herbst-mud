import { createContext, useState, useEffect, useCallback, type ReactNode } from "react";

export type ThemeName = "dark" | "light";

type ThemeContextValue = {
  readonly theme: ThemeName;
  readonly toggle: () => void;
  readonly setTheme: (t: ThemeName) => void;
};

const STORAGE_KEY = "herbst-mud-theme";

export const ThemeContext = createContext<ThemeContextValue | null>(null);

export function ThemeProvider({ children }: { readonly children: ReactNode }) {
  const [theme, setThemeState] = useState<ThemeName>(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored === "light") return "light";
    if (stored === "dark") return "dark";
    return "dark";
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