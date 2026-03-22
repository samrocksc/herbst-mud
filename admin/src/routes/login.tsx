import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'

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
      const response = await fetch('http://localhost:8080/users/auth', {
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
      console.error('Login error:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-[#1a1612] flex flex-col items-center justify-center text-[#e8dcc4]">
      <h1 className="text-3xl mb-2">Herbst MUD Admin</h1>
      <Link to="/" className="text-[#4a7c4e] no-underline mb-6 hover:text-[#5a9c5e]">
        ← Back to Home
      </Link>

      {error && <div className="text-[#8b4444] mb-4">{error}</div>}

      <form onSubmit={handleSubmit} className="w-[300px] max-w-[90vw]">
        <div className="mb-4">
          <label htmlFor="username" className="block mb-2">Username:</label>
          <input
            type="text" id="username" name="username" value={username}
            onChange={(e) => setUsername(e.target.value)} required
            className="w-full p-2 bg-[#2d2416] border border-[#5a4a35] rounded text-[#e8dcc4]"
          />
        </div>
        <div className="mb-4">
          <label htmlFor="password" className="block mb-2">Password:</label>
          <input
            type="password" id="password" name="password" value={password}
            onChange={(e) => setPassword(e.target.value)} required
            className="w-full p-2 bg-[#2d2416] border border-[#5a4a35] rounded text-[#e8dcc4]"
          />
        </div>
        <button
          type="submit" disabled={loading}
          className="w-full py-3 bg-[#4a7c4e] border-none rounded text-[#e8dcc4] cursor-pointer disabled:opacity-70 hover:bg-[#5a9c5e]"
        >
          {loading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  )
}