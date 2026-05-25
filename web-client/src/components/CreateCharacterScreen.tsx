import { useEffect, useState } from "react";
import { createCharacter, listRaces, listGenders, type Race, type Gender, type Character } from "../lib/api";
import { useTheme } from "../lib/theme";
import { Button, Input, Select } from "../ui";

export type CreateCharacterScreenProps = {
  worldName: string
  onCreated: (char: Character) => void
  onBack: () => void
  onLogout: () => void
};

export default function CreateCharacterScreen( { worldName, onCreated, onBack, onLogout  }: Readonly<CreateCharacterScreenProps>) {
  const [name, setName] = useState("");
  const [race, setRace] = useState("");
  const [gender, setGender] = useState("");
  const [races, setRaces] = useState<readonly Race[]>([]);
  const [genders, setGenders] = useState<readonly Gender[]>([]);
  const [loadingMeta, setLoadingMeta] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [nameError, setNameError] = useState("");
  const { theme, toggle } = useTheme();

  useEffect(() => {
    setLoadingMeta(true);
    Promise.all([listRaces(), listGenders()])
      .then(([r, g]) => {
        setRaces(r);
        setGenders(g);
        if (r.length > 0) setRace(r[0].name);
        if (g.length > 0) setGender(g[0].name);
      })
      .catch((e) => setError(e.message || "Failed to load options"))
      .finally(() => setLoadingMeta(false));
  }, []);

  function validateName(v: string): boolean {
    if (v.length < 1 || v.length > 23) {
      setNameError("Name must be 1–23 characters");
      return false;
    }
    if (!/^[a-zA-Z]+$/.test(v)) {
      setNameError("Letters only (a-z, A-Z)");
      return false;
    }
    setNameError("");
    return true;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!validateName(name)) return;
    setSaving(true);
    setError("");
    try {
      const char = await createCharacter({ name, race, gender, world: worldName });
      onCreated(char);
    } catch (e: unknown) {
      setError((e instanceof Error ? e.message : String(e)) || "Failed to create character");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      <header className="shrink-0 flex items-center justify-between px-3 py-2 border-b border-border bg-surface">
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={onBack}>← Back</Button>
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
        <div className="max-w-lg mx-auto space-y-6">
          <h2 className="text-lg font-mono font-bold">
            Create Character — {worldName}
          </h2>

          {loadingMeta && (
            <p className="text-center font-mono text-sm text-muted">
              Loading options...
            </p>
          )}

          {!loadingMeta && (
            <form onSubmit={handleSubmit} className="space-y-5">
              <Input
                label="Name"
                type="text"
                value={name}
                onChange={(e) => {
                  setName(e.target.value);
                  if (nameError) validateName(e.target.value);
                }}
                error={nameError || undefined}
                maxLength={23}
                autoFocus
                autoComplete="off"
                spellCheck={false}
              />

              <Select
                label="Race"
                options={races.map((r) => ({
                  value: r.name,
                  label: r.display_name || r.name,
                }))}
                value={race}
                onChange={(e) => setRace(e.target.value)}
              />

              <Select
                label="Gender"
                options={genders.map((g) => ({
                  value: g.name,
                  label: g.display_name || g.name,
                }))}
                value={gender}
                onChange={(e) => setGender(e.target.value)}
              />

              {error && (
                <p className="text-sm font-mono text-danger">{error}</p>
              )}

              <Button type="submit" variant="primary" size="lg" fullWidth disabled={saving}>
                {saving ? "Creating..." : "Create Character"}
              </Button>
            </form>
          )}
        </div>
      </main>
    </div>
  );
}