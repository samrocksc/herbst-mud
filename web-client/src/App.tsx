import { useState, useEffect, useCallback } from "react";
import { ThemeProvider } from "./lib/theme";
import LoginScreen from "./components/LoginScreen";
import WorldScreen from "./components/WorldScreen";
import CharacterScreen from "./components/CharacterScreen";
import CreateCharacterScreen from "./components/CreateCharacterScreen";
import GameScreen from "./components/GameScreen";
import { MUDConnectionProvider } from "./context/MUDConnectionProvider";
import { me, type Character } from "./lib/api";

function AppInner() {
  const [phase, setPhase] = useState<"checking" | "login" | "world" | "character" | "create" | "playing">("checking");
  const [selectedWorldId, setSelectedWorldId] = useState<string | null>(null);
  const [selectedWorldName, setSelectedWorldName] = useState<string | null>(null);
  const [selectedCharacter, setSelectedCharacter] = useState<Character | null>(null);

  useEffect(() => {
    console.log("[APP] AppInner render cycle");
    const token = localStorage.getItem("herbst_token");
    if (!token) {
      setPhase("login");
      return;
    }
    me()
      .then(() => setPhase("world"))
      .catch(() => {
        localStorage.removeItem("herbst_token");
        setPhase("login");
      });
  }, []);

  const doLogout = useCallback(() => {
    localStorage.removeItem("herbst_token");
    setSelectedWorldId(null);
    setSelectedWorldName(null);
    setSelectedCharacter(null);
    setPhase("login");
  }, []);

  if (phase === "checking") {
    return (
      <div
        className="min-h-screen flex items-center justify-center font-mono"
        style={{ backgroundColor: "var(--mud-bg)", color: "var(--mud-fg)" }}
      >
        <p style={{ color: "var(--mud-muted)" }}>Verifying session...</p>
      </div>
    );
  }

  if (phase === "login") {
    return <LoginScreen onLogin={() => setPhase("world")} />;
  }

  if (phase === "create" && selectedWorldId && selectedWorldName) {
    return (
      <CreateCharacterScreen
        worldId={selectedWorldId}
        worldName={selectedWorldName}
        onCreated={(char) => {
          setSelectedCharacter(char);
          setPhase("playing");
        }}
        onBack={() => setPhase("character")}
        onLogout={doLogout}
      />
    );
  }

  if (phase === "playing" && selectedWorldId && selectedWorldName && selectedCharacter) {
    const token = localStorage.getItem("herbst_token") || "";
    return (
      <GameScreen
        worldId={selectedWorldId}
        worldName={selectedWorldName}
        character={selectedCharacter}
        token={token}
        onDisconnect={() => {
          setSelectedCharacter(null);
          setPhase("character");
        }}
      />
    );
  }

  if (phase === "character" && selectedWorldId && selectedWorldName) {
    return (
      <CharacterScreen
        worldId={selectedWorldId}
        worldName={selectedWorldName}
        onSelectCharacter={(char) => {
          setSelectedCharacter(char);
          setPhase("playing");
        }}
        onCreateNew={() => setPhase("create")}
        onBack={() => {
          setSelectedWorldId(null);
          setSelectedWorldName(null);
          setPhase("world");
        }}
        onLogout={doLogout}
      />
    );
  }

  return <WorldScreen
    onLogout={doLogout}
    onSelectWorld={(worldId, worldName) => {
      setSelectedWorldId(worldId);
      setSelectedWorldName(worldName);
      setPhase("character");
    }}
  />;
}

export default function App() {
  return (
    <ThemeProvider>
      <MUDConnectionProvider>
        <AppInner />
      </MUDConnectionProvider>
    </ThemeProvider>
  );
}
