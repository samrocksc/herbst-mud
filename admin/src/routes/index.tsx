import { createFileRoute, Link } from '@tanstack/react-router'
import { useEffect } from 'react'
import { useNavigate } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: Home,
})

function Home() {
  const navigate = useNavigate()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      navigate({ to: '/dashboard' })
    }
  }, [navigate])

  return (
    <div className="min-h-screen bg-surface flex flex-col items-center justify-center text-text">
      <h1 className="text-4xl mb-4 text-primary">Herbst MUD Admin</h1>
      <p className="text-text-muted mb-8">Welcome to the administrative backend for Herbst MUD.</p>
      <Link
        to="/login"
        className="text-primary no-underline py-2 px-4 border-2 border-primary rounded hover:bg-surface-muted transition-colors"
      >
        Login to Admin Panel
      </Link>
    </div>
  )
}