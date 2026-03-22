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
    <div className="min-h-screen bg-[#1a1612] flex flex-col items-center justify-center text-[#e8dcc4]">
      <h1 className="text-4xl mb-4 text-[#4a7c4e]">Herbst MUD Admin</h1>
      <p className="text-[#a89070] mb-8">Welcome to the administrative backend for Herbst MUD.</p>
      <Link
        to="/login"
        className="text-[#4a7c4e] no-underline py-2 px-4 border-2 border-[#4a7c4e] rounded hover:bg-[#2d2416]"
      >
        Login to Admin Panel
      </Link>
    </div>
  )
}