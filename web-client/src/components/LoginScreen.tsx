import { useState } from "react";
import { login } from "../lib/api";
import { useTheme } from "../lib/theme";
import { Button, Card, CardBody, Input } from "../ui";

export type LoginScreenProps = {
  onLogin: (token: string) => void
  error?: string
};

export default function LoginScreen( { onLogin, error  }: Readonly<LoginScreenProps>) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState(error || "");
  const { toggle, theme } = useTheme();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setErr("");
    try {
      const { token } = await login(email, password);
      localStorage.setItem("herbst_token", token);
      onLogin(token);
    } catch (e: unknown) {
      setErr((e instanceof Error ? e.message : String(e)) || "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center px-4 bg-background">
      <Button
        variant="ghost"
        size="sm"
        onClick={toggle}
        className="absolute top-4 right-4"
      >
        {theme === "dark" ? "🌙" : "☀️"}
      </Button>

      <Card className="w-full max-w-sm">
        <CardBody className="space-y-5">
          <h1 className="text-center text-2xl font-mono font-bold tracking-wide text-accent">
            HERBST MUD
          </h1>

          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="Username / Email"
              type="text"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              placeholder="username"
              autoComplete="username"
            />

            <Input
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              autoComplete="current-password"
            />

            {err && (
              <p className="text-sm font-mono text-center text-danger">{err}</p>
            )}

            <Button type="submit" variant="primary" size="lg" fullWidth disabled={loading}>
              {loading ? "Connecting..." : "ENTER THE WORLD"}
            </Button>
          </form>
        </CardBody>
      </Card>
    </div>
  );
}