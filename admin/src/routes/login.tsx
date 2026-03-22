import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/login')({
  component: Login,
})

function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

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

      navigate({ to: '/dashboard' })
    } catch (err) {
      setError('Connection error. Please try again.')
      console.error('Login error:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-[#0a0a0f] flex flex-col items-center justify-center text-white">
      <h1 className="text-3xl mb-2">Admin Login</h1>
      <Link to="/" className="text-[#61dafb] no-underline mb-6 hover:underline">
        ← Back to Home
      </Link>

      {error && <div className="text-red-400 mb-4">{error}</div>}

      <form onSubmit={handleSubmit} className="w-[300px] max-w-[90vw]">
        <div className="mb-4">
          <label htmlFor="username" className="block mb-2">Username:</label>
          <input
            type="text"
            id="username"
            name="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            className="w-full p-2 bg-[#1a1a2e] border border-[#333] rounded text-white"
          />
        </div>
        <div className="mb-4">
          <label htmlFor="password" className="block mb-2">Password:</label>
          <input
            type="password"
            id="password"
            name="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="w-full p-2 bg-[#1a1a2e] border border-[#333] rounded text-white"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full py-3 bg-[#27ae60] border-none rounded text-white cursor-pointer disabled:opacity-70 hover:bg-[#2ecc71]"
        >
          {loading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  )
}