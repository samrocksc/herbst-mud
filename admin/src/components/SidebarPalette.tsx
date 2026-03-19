import { useCallback } from 'react'

export function SidebarPalette() {
  const onDragStart = useCallback((event: React.DragEvent, nodeType: string) => {
    event.dataTransfer.setData('application/reactflow', nodeType)
    event.dataTransfer.effectAllowed = 'move'
  }, [])

  return (
    <aside style={{
      width: '200px',
      padding: '16px',
      background: '#1a1a1a',
      borderRadius: '8px',
      border: '1px solid #333'
    }}>
      <h3 style={{ marginBottom: '16px', color: '#fff' }}>Drag to Add</h3>
      
      <div
        draggable
        onDragStart={(e) => onDragStart(e, 'room')}
        style={{
          padding: '12px 16px',
          marginBottom: '8px',
          background: '#2d5a27',
          borderRadius: '6px',
          color: '#fff',
          cursor: 'grab',
          display: 'flex',
          alignItems: 'center',
          gap: '8px',
          userSelect: 'none'
        }}
      >
        <span style={{ fontSize: '20px' }}>🟩</span>
        <span>New Room</span>
      </div>

      <div
        draggable
        onDragStart={(e) => onDragStart(e, 'npc')}
        style={{
          padding: '12px 16px',
          marginBottom: '8px',
          background: '#5a2727',
          borderRadius: '6px',
          color: '#fff',
          cursor: 'grab',
          display: 'flex',
          alignItems: 'center',
          gap: '8px',
          userSelect: 'none'
        }}
      >
        <span style={{ fontSize: '20px' }}>👤</span>
        <span>New NPC</span>
      </div>

      <div
        draggable
        onDragStart={(e) => onDragStart(e, 'item')}
        style={{
          padding: '12px 16px',
          background: '#27475a',
          borderRadius: '6px',
          color: '#fff',
          cursor: 'grab',
          display: 'flex',
          alignItems: 'center',
          gap: '8px',
          userSelect: 'none'
        }}
      >
        <span style={{ fontSize: '20px' }}>📦</span>
        <span>New Item</span>
      </div>

      <p style={{ marginTop: '16px', fontSize: '12px', color: '#888' }}>
        Drag items onto the map canvas to create them.
      </p>
    </aside>
  )
}