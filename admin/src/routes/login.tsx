import { createFileRoute, Link } from '@tanstack/react-router';
import { useState } from 'react';
import { Button } from '../components/Button';
import { FormField } from '../components/fields/FormField';
import { FormError } from '../components/fields/FormError';

export const Route = createFileRoute('/login')({
  component: Login,
});

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await fetch(`${window.location.origin}/users/auth`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: username, password }),
      });

      if (!response.ok) {
        const data = await response.json();
        setError(data.error || 'Login failed');
        return;
      }

      const data = await response.json();
      localStorage.setItem('token', data.token);
      localStorage.setItem('userId', data.id);
      localStorage.setItem('email', data.email);
      localStorage.setItem('isAdmin', data.is_admin);

      window.location.href = '/dashboard';
    } catch {
      setError('Connection error. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-surface flex flex-col items-center justify-center text-text">
      <h1 className="text-3xl mb-2">Herbst MUD Admin</h1>
      <Link to="/" className="text-primary no-underline mb-6 hover:text-primary">
        ← Back to Home
      </Link>

      {error && <FormError message={error} />}

      <form onSubmit={handleSubmit} className="w-[300px] max-w-[90vw]">
        <div className="mb-4">
          <FormField label="Username" value={username} onChange={setUsername} />
        </div>
        <div className="mb-4">
          <FormField label="Password" value={password} onChange={setPassword} type="password" />
        </div>
        <Button type="submit" disabled={loading} variant="primary" fullWidth>
          {loading ? 'Logging in...' : 'Login'}
        </Button>
      </form>
    </div>
  );
}