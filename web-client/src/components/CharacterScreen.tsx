import { useEffect, useState } from "react";
import { listMyCharacters, type Character } from "../lib/api";
import { useTheme } from "../lib/theme";
import { Button, Card, CardBody } from "../ui";

export type CharacterScreenProps = {
  worldName: string
  onSelectCharacter: (char: Character) => void
  onCreateNew: () => void
  onBack: () => void
  onLogout: () => void
};

export default function CharacterScreen( { worldName, onSelectCharacter, onCreateNew, onBack, onLogout  }: Readonly<CharacterScreenProps>) {
  const [characters, setCharacters] = useState<Character[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const { theme, toggle } = useTheme();

  useEffect(() => {
    setLoading(true);
    listMyCharacters()
      .then((all) => {
        const filtered = all.filter((c) => c.currentWorld === worldName);
        setCharacters(filtered);
        setError("");
      })
      .catch((e) => {
        setError(e.message || "Failed to load characters");
      })
      .finally(() => setLoading(false));
  }, [worldName]);

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      <header className="shrink-0 flex items-center justify-between px-3 py-2 border-b border-border bg-surface">
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={onBack}>
            ← Worlds
          </Button>
          <span className="font-bold text-lg text-accent">HERBST MUD</span>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={toggle}>
            {theme === "dark" ? "🌙" : "☀️"}
          </Button>
          <Button variant="ghost" size="sm" onClick={onLogout}>
            Logout
          </Button>
        </div>
      </header>

      <main className="flex-1 px-4 py-8">
        <div className="max-w-3xl mx-auto">
          <h2 className="text-lg font-mono font-bold mb-4">
            {worldName} — Select Character
          </h2>

          {loading && (
            <p className="text-center font-mono text-sm text-muted">
              Loading characters...
            </p>
          )}

          {error && !loading && (
            <p className="text-center font-mono text-sm text-danger">{error}</p>
          )}

          {!loading && !error && characters.length === 0 && (
            <div className="text-center space-y-4">
              <p className="text-sm font-mono text-muted">
                No characters in this world yet.
              </p>
              <Button variant="primary" onClick={onCreateNew}>
                Create New Character
              </Button>
            </div>
          )}

          <ul className="space-y-3">
            {characters.map((char) => (
              <li key={char.id}>
                <Card hover onClick={() => onSelectCharacter(char)}>
                  <CardBody>
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-bold text-sm">{char.name}</h3>
                        <p className="text-[11px] text-muted mt-1">
                          {char.race} • {char.gender || "unspecified"} • Level {char.level}
                        </p>
                      </div>
                      <div className="text-[11px] text-right text-muted">
                        <div>HP {char.hitpoints}/{char.max_hitpoints}</div>
                        <div>STA {char.stamina}/{char.max_stamina}</div>
                      </div>
                    </div>
                    {char.class && (
                      <span className="inline-block mt-2 px-2 py-0.5 rounded text-[11px] bg-accent text-background">
                        {char.class}
                      </span>
                    )}
                  </CardBody>
                </Card>
              </li>
            ))}
          </ul>

          {characters.length > 0 && characters.length < 3 && (
            <div className="mt-6 text-center">
              <Button variant="secondary" onClick={onCreateNew}>
                Create New Character ({3 - characters.length} slot{3 - characters.length !== 1 ? "s" : ""} left)
              </Button>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}