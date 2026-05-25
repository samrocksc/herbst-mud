import { useState, useEffect, useCallback } from "react";
import { ThemeProvider } from "./lib/theme";
import LoginScreen from "./components/LoginScreen";
import WorldScreen from "./components/WorldScreen";
import CharacterScreen from "./components/CharacterScreen";
import CreateCharacterScreen from "./components/CreateCharacterScreen";
import GameScreen from "./components/GameScreen";
import { me, type Character } from "./lib/api";

function AppInner() {
  const [phase, setPhase] = useState<"checking" | "login" | "world" | "character" | "create" | "playing">("checking");
  const [selectedWorld, setSelectedWorld] = useState<string | null>(null);
  const [selectedCharacter, setSelectedCharacter] = useState<Character | null>(null);

  useEffect(() => {
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
    setSelectedWorld(null);
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

  if (phase === "create" && selectedWorld) {
    return (
      <CreateCharacterScreen
        worldName={selectedWorld}
        onCreated={(char) => {
          // Auto-enter game with newly created character
          setSelectedCharacter(char);
          setPhase("playing");
        }}
        onBack={() => setPhase("character")}
        onLogout={doLogout}
      />
    );
  }

  if (phase === "playing" && selectedWorld && selectedCharacter) {
    const token = localStorage.getItem("herbst_token") || "";
    return (
      <GameScreen
        worldName={selectedWorld}
        character={selectedCharacter}
        token={token}
        onDisconnect={() => {
          setSelectedCharacter(null);
          setPhase("character");
        }}
      />
    );
  }

  if (phase === "character" && selectedWorld) {
    return (
      <CharacterScreen
        worldName={selectedWorld}
        onSelectCharacter={(char) => {
          setSelectedCharacter(char);
          setPhase("playing");
        }}
        onCreateNew={() => setPhase("create")}
        onBack={() => {
          setSelectedWorld(null);
          setPhase("world");
        }}
        onLogout={doLogout}
      />
    );
  }

  return <WorldScreen
    onLogout={doLogout}
    onSelectWorld={(name) => {
      setSelectedWorld(name);
      setPhase("character");
    }}
  />;
}

export default function App() {
  return (
    <ThemeProvider>
      <AppInner />
    </ThemeProvider>
  );
}