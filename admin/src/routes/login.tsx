import { createFileRoute, Link } from '@tanstack/react-router'
import { logError } from '../utils/log'
import { useState } from 'react'
import { Button } from '../components/Button'

export const Route = createFileRoute('/login')({
  component: Login,
})

function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await fetch(`${window.location.origin}/users/auth`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: username, password }),
      })

      if (!response.ok) {
        const data = await response.json()
        setError(data.error || 'Login failed')
        return
      }

      const data = await response.json()
      localStorage.setItem('token', data.token)
      localStorage.setItem('userId', data.id)
      localStorage.setItem('email', data.email)
      localStorage.setItem('isAdmin', data.is_admin)

      window.location.href = '/dashboard'
    } catch (err) {
      setError('Connection error. Please try again.')
      logError('Login error:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-surface flex flex-col items-center justify-center text-text">
      <h1 className="text-3xl mb-2">Herbst MUD Admin</h1>
      <Link to="/" className="text-primary no-underline mb-6 hover:text-primary">
        ← Back to Home
      </Link>

      {error && <div className="text-danger mb-4">{error}</div>}

      <form onSubmit={handleSubmit} className="w-[300px] max-w-[90vw]">
        <div className="mb-4">
          <label htmlFor="username" className="block mb-2">Username:</label>
          <input
            type="text" id="username" name="username" value={username}
            onChange={(e) => setUsername(e.target.value)} required
            className="w-full p-2 bg-surface border border-border rounded text-text"
          />
        </div>
        <div className="mb-4">
          <label htmlFor="password" className="block mb-2">Password:</label>
          <input
            type="password" id="password" name="password" value={password}
            onChange={(e) => setPassword(e.target.value)} required
            className="w-full p-2 bg-surface border border-border rounded text-text"
          />
        </div>
        <Button
          type="submit" disabled={loading}
          variant="primary"
          fullWidth
        >
          {loading ? 'Logging in...' : 'Login'}
        </Button>
      </form>
    </div>
  )
}