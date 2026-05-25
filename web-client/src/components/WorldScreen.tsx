import { useEffect, useState } from "react";
import { me, listWorlds, type World } from "../lib/api";
import { useTheme } from "../lib/theme";
import { Button, Card, CardBody } from "../ui";

export type WorldScreenProps = {
  onLogout: () => void
  onSelectWorld: (name: string) => void
};

export default function WorldScreen( { onLogout, onSelectWorld  }: Readonly<WorldScreenProps>) {
  const [user, setUser] = useState<{ email: string; is_admin: boolean } | null>(null);
  const [worlds, setWorlds] = useState<readonly World[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const { theme, toggle } = useTheme();

  useEffect(() => {
    setLoading(true);
    Promise.all([me(), listWorlds()])
      .then(([u, w]) => {
        setUser(u);
        setWorlds(w.worlds);
        setError("");
      })
      .catch((e) => {
        setError(e.message || "Failed to load");
        if (e.message?.includes("Session") || e.message?.includes("token")) {
          onLogout();
        }
      })
      .finally(() => setLoading(false));
  }, [onLogout]);

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      <header className="shrink-0 flex items-center justify-between px-3 py-2 border-b border-border bg-surface">
        <span className="font-bold text-lg text-accent">HERBST MUD</span>
        <div className="flex items-center gap-2">
          <span className="text-xs text-muted hidden sm:inline">{user?.email}</span>
          <Button variant="ghost" size="sm" onClick={toggle}>
            {theme === "dark" ? "🌙" : "☀️"}
          </Button>
          <Button variant="ghost" size="sm" onClick={onLogout}>
            Logout
          </Button>
        </div>
      </header>

      <main className="flex-1 px-4 py-8">
        {loading && (
          <p className="text-center font-mono text-sm text-muted">
            Loading worlds...
          </p>
        )}

        {error && !loading && (
          <p className="text-center font-mono text-sm text-danger">{error}</p>
        )}

        {!loading && !error && (
          <div className="max-w-3xl mx-auto space-y-4">
            <h2 className="text-lg font-mono font-bold mb-4">
              Available Worlds
            </h2>

            {worlds.length === 0 && (
              <p className="text-sm font-mono text-muted">No worlds found.</p>
            )}

            <ul className="space-y-3">
              {worlds.map((w) => (
                <li key={w.name}>
                  <Card hover onClick={() => onSelectWorld(w.name)}>
                    <CardBody>
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-bold text-sm">{w.name}</h3>
                          <p className="text-[11px] text-muted mt-1">{w.file}</p>
                        </div>
                        <span className="text-xs text-accent">→</span>
                      </div>
                    </CardBody>
                  </Card>
                </li>
              ))}
            </ul>
          </div>
        )}
      </main>
    </div>
  );
}