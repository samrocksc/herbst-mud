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
    <div className="min-h-screen bg-[#0a0a0f] flex flex-col items-center justify-center text-white">
      <h1 className="text-4xl mb-4">Herbst MUD Admin</h1>
      <p className="text-[#888] mb-8">Welcome to the administrative backend for Herbst MUD.</p>
      <Link
        to="/login"
        className="text-[#61dafb] no-underline py-2 px-4 border border-[#61dafb] rounded hover:bg-[rgba(97,218,251,0.1)]"
      >
        Login to Admin Panel
      </Link>
    </div>
  )
}