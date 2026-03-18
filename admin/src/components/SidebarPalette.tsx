// Native HTML5 drag and drop sidebar palette
// This is simpler and more reliable than dnd-kit for our use case

interface DraggableItemProps {
  type: 'room' | 'npc'
  label: string
  icon: string
}

function DraggableItem({ type, label, icon }: DraggableItemProps) {
  const handleDragStart = (e: React.DragEvent) => {
    e.dataTransfer.setData('application/x-herbstmud-node-type', type)
    e.dataTransfer.effectAllowed = 'copy'
  }

  return (
    <div
      draggable
      onDragStart={handleDragStart}
      className="sidebar-item"
      style={{
        padding: '12px',
        background: '#2d2d44',
        border: '1px solid #444',
        borderRadius: '6px',
        cursor: 'grab',
        marginBottom: '8px',
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
      }}
    >
      <span style={{ fontSize: '16px' }}>{icon}</span>
      <span style={{ color: '#fff', fontSize: '13px' }}>{label}</span>
    </div>
  )
}

export function SidebarPalette() {
  return (
    <aside className="sidebar-palette" style={{
      width: '160px',
      padding: '12px',
      background: '#1a1a2e',
      borderRadius: '8px',
      border: '1px solid #333',
    }}>
      <h4 style={{ margin: '0 0 12px 0', color: '#fff', fontSize: '14px' }}>
        Templates
      </h4>
      <DraggableItem type="room" label="New Room" icon="[+]" />
      <DraggableItem type="npc" label="New NPC" icon="[@]" />
      <div style={{ marginTop: '16px', padding: '8px', background: '#222', borderRadius: '4px', fontSize: '11px', color: '#888' }}>
        <p style={{ margin: '0 0 4px 0' }}>Tip:</p>
        <p style={{ margin: 0 }}>Drag items onto canvas</p>
      </div>
    </aside>
  )
}