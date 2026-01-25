import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: Home,
})

function Home() {
  return (
    <div className="home">
      <h1>Herbst MUD Admin</h1>
      <p>Welcome to the administrative backend for Herbst MUD.</p>
      <a href="/login">Login to Admin Panel</a>
    </div>
  )
}