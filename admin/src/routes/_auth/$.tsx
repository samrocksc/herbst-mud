import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/$')({
  component: NotFound,
})

function NotFound() {
  return (
    <div 
      style={{ 
        display: 'flex', 
        flexDirection: 'column',
        alignItems: 'center', 
        justifyContent: 'center',
        height: '100vh',
        textAlign: 'center',
        padding: '20px'
      }}
    >
      <h1 style={{ fontSize: '4rem', margin: 0, color: '#e74c3c' }}>404</h1>
      <h2 style={{ color: '#ecf0f1' }}>Page Not Found</h2>
      <p style={{ color: '#95a5a6', marginBottom: '30px' }}>
        The page you're looking for doesn't exist or has been moved.
      </p>
      <a 
        href="/_auth/dashboard"
        style={{
          background: '#27ae60',
          color: 'white',
          padding: '12px 24px',
          borderRadius: '6px',
          textDecoration: 'none',
          fontWeight: 'bold'
        }}
      >
        ← Go to Dashboard
      </a>
    </div>
  )
}