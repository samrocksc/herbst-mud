import { useEffect, useRef, useState } from "react";
import { deleteCharacter, listMyCharacters, type Character } from "../lib/api";
import { useTheme } from "../lib/theme";
import { Button, Card, CardBody } from "../ui";

export type CharacterScreenProps = {
  worldId: string;
  worldName: string;
  onSelectCharacter: (char: Character) => void;
  onCreateNew: () => void;
  onBack: () => void;
  onLogout: () => void;
};

export default function CharacterScreen({ worldId, worldName, onSelectCharacter, onCreateNew, onBack, onLogout }: Readonly<CharacterScreenProps>) {
  const [characters, setCharacters] = useState<Character[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<Character | null>(null);
  const [confirmName, setConfirmName] = useState("");
  const [deleting, setDeleting] = useState(false);
  const { theme, toggle } = useTheme();

  // Ref to store the handler functions - only set once
  const onSelectRef = useRef(onSelectCharacter);
  const onCreateRef = useRef(onCreateNew);
  const onBackRef = useRef(onBack);
  const onLogoutRef = useRef(onLogout);

  // Only update refs when first created, never again
  useEffect(() => {
    if (onSelectRef.current === onSelectCharacter && onCreateRef.current === onCreateNew) {
      return;
    }
    onSelectRef.current = onSelectCharacter;
    onCreateRef.current = onCreateNew;
    onBackRef.current = onBack;
    onLogoutRef.current = onLogout;
  }, [onSelectCharacter, onCreateNew, onBack, onLogout]);

  useEffect(() => {
    setLoading(true);
    listMyCharacters()
      .then((all) => {
        const filtered = all.filter((c) => {
          const currentWorld = String(c.currentWorld);
          const worldIdStr = String(worldId);
          if (currentWorld === worldIdStr) return true;
          if (worldName && currentWorld === worldName) return true;
          return false;
        });
        setCharacters(filtered);
        setError("");
      })
      .catch((e) => {
        setError(e.message || "Failed to load characters");
      })
      .finally(() => setLoading(false));
  }, [worldId, worldName]);

  const handleDeleteRequest = (char: Character, e: React.MouseEvent) => {
    e.stopPropagation();
    setDeleteTarget(char);
    setConfirmName("");
  };

  const handleDeleteConfirm = async () => {
    if (!deleteTarget || confirmName !== deleteTarget.name) return;
    setDeleting(true);
    try {
      await deleteCharacter(deleteTarget.id);
      setCharacters((prev) => prev.filter((c) => c.id !== deleteTarget.id));
      setDeleteTarget(null);
      setConfirmName("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Delete failed");
    } finally {
      setDeleting(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      <header className="shrink-0 flex items-center justify-between px-3 py-2 border-b border-border bg-surface">
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={onBackRef.current}>
            ← Worlds
          </Button>
          <span className="font-bold text-lg text-accent">HERBST MUD</span>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={toggle}>
            {theme === "dark" ? "🌙" : "☀️"}
          </Button>
          <Button variant="ghost" size="sm" onClick={onLogoutRef.current}>
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
              <Button variant="primary" onClick={onCreateRef.current}>
                Create New Character
              </Button>
            </div>
          )}

          <ul className="space-y-3">
            {characters.map((char) => (
              <li key={char.id}>
                <Card hover onClick={() => onSelectRef.current(char)}>
                  <CardBody>
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-bold text-sm">{char.name}</h3>
                        <p className="text-[11px] text-muted mt-1">
                          {char.race} • {char.gender || "unspecified"} • Level {char.level}
                        </p>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="text-[11px] text-right text-muted">
                          <div>HP {Math.min(char.hitpoints, char.max_hitpoints)}/{char.max_hitpoints}</div>
                          <div>STA {Math.min(char.stamina, char.max_stamina)}/{char.max_stamina}</div>
                        </div>
                        <Button
                          variant="danger"
                          size="sm"
                          onClick={(e) => handleDeleteRequest(char, e)}
                        >
                          Delete
                        </Button>
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
              <Button variant="secondary" onClick={onCreateRef.current}>
                Create New Character ({3 - characters.length} slot{3 - characters.length !== 1 ? "s" : ""} left)
              </Button>
            </div>
          )}
        </div>
      </main>

      {deleteTarget && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60">
          <div className="bg-surface border border-border rounded-lg p-6 w-80 space-y-4">
            <h3 className="font-bold text-sm text-danger">Delete Character</h3>
            <p className="text-xs text-muted">
              Type <span className="font-bold text-foreground">{deleteTarget.name}</span> to confirm deletion. This cannot be undone.
            </p>
            <input
              type="text"
              value={confirmName}
              onChange={(e) => setConfirmName(e.target.value)}
              placeholder={deleteTarget.name}
              className="w-full rounded border border-border bg-surface-alt px-3 py-2 font-mono text-sm text-foreground outline-none focus:border-accent"
            />
            <div className="flex justify-end gap-2">
              <Button variant="ghost" size="sm" onClick={() => setDeleteTarget(null)}>
                Cancel
              </Button>
              <Button
                variant="danger"
                size="sm"
                disabled={confirmName !== deleteTarget.name || deleting}
                onClick={handleDeleteConfirm}
              >
                {deleting ? "Deleting..." : "Delete"}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
