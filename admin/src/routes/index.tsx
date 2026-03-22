import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useEffect, useState } from 'react'

export const Route = createFileRoute('/')({
  component: Home,
})

function Home() {
  const [dimensions, setDimensions] = useState({ width: window.innerWidth, height: window.innerHeight })
  const navigate = useNavigate()

  useEffect(() => {
    const handleResize = () => {
      setDimensions({ width: window.innerWidth, height: window.innerHeight })
    }
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  // Check if already logged in
  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      navigate({ to: '/map' })
    }
  }, [navigate])

  return (
    <div
      className="home"
      style={{
        width: dimensions.width,
        height: dimensions.height,
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        margin: 0,
        padding: 0
      }}
    >
      <h1>Herbst MUD Admin</h1>
      <p>Welcome to the administrative backend for Herbst MUD.</p>
      <a href="/login" style={{
        color: '#61dafb',
        textDecoration: 'none',
        padding: '0.5rem 1rem',
        border: '1px solid #61dafb',
        borderRadius: '4px'
      }}>
        Login to Admin Panel
      </a>
    </div>
  )
}