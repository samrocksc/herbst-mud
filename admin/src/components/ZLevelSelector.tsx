import { useState, useMemo } from 'react'

interface ZLevelSelectorProps {
  currentLevel: number
  onChange: (level: number) => void
}

const Z_LEVELS = [
  { level: -2, label: 'Deep Underground', short: 'Z: -2' },
  { level: -1, label: 'Underground', short: 'Z: -1' },
  { level: 0, label: 'Ground', short: 'Z: 0' },
  { level: 1, label: 'Upper Floor', short: 'Z: 1' },
  { level: 2, label: 'Tower', short: 'Z: 2' },
]

export function ZLevelSelector({ currentLevel, onChange }: ZLevelSelectorProps) {
  const [adjacentLevels, setAdjacentLevels] = useState<number[]>([])

  const handleLevelClick = (level: number) => {
    onChange(level)
    // Show adjacent levels for preview
    const adjacent = [level - 1, level + 1].filter(l => l >= -2 && l <= 2)
    setAdjacentLevels(adjacent)
  }

  return (
    <div className="z-level-selector" style={{
      display: 'flex',
      gap: '8px',
      padding: '8px 16px',
      background: '#1a1a2e',
      borderRadius: '8px',
      marginBottom: '16px',
      flexWrap: 'wrap'
    }}>
      <span style={{ color: '#888', alignSelf: 'center', marginRight: '8px' }}>Floor:</span>
      {Z_LEVELS.map(({ level, label, short }) => (
        <button
          key={level}
          onClick={() => handleLevelClick(level)}
          className={currentLevel === level ? 'active' : ''}
          style={{
            padding: '8px 16px',
            borderRadius: '6px',
            border: 'none',
            cursor: 'pointer',
            fontWeight: currentLevel === level ? 'bold' : 'normal',
            background: currentLevel === level ? '#00AAAA' : '#333',
            color: currentLevel === level ? '#fff' : '#aaa',
            opacity: !adjacentLevels.includes(level) && level !== currentLevel ? 0.5 : 1,
            transition: 'all 0.2s ease',
          }}
          title={label}
        >
          {short}
          <span style={{ display: 'block', fontSize: '10px', color: '#666' }}>
            {label}
          </span>
        </button>
      ))}
      
      {adjacentLevels.length > 0 && (
        <div style={{ 
          fontSize: '11px', 
          color: '#666', 
          alignSelf: 'center',
          marginLeft: '8px'
        }}>
          Adjacent levels shown faintly
        </div>
      )}
    </div>
  )
}

export function useZLevelFilter(nodes: any[], currentLevel: number) {
  return useMemo(() => {
    return nodes.filter(node => {
      const zLevel = node.data?.zLevel ?? 0
      return zLevel === currentLevel || Math.abs(zLevel - currentLevel) === 1
    })
  }, [nodes, currentLevel])
}

export function getZLevelLabel(level: number): string {
  const found = Z_LEVELS.find(z => z.level === level)
  return found ? found.label : `Z: ${level}`
}