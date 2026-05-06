import { Button } from '../Button'

type MapToolbarProps = {
  currentZLevel: number
  zoom: number
  onZoom: (delta: number) => void
  onResetView?: () => void
  onRelayout?: () => void
  onCleanupOrphanExits?: () => void
}

export function MapToolbar({ currentZLevel, zoom, onZoom, onResetView, onRelayout, onCleanupOrphanExits }: MapToolbarProps) {
  return (
    <div className="absolute top-0 left-0 right-0 p-3 bg-surface-muted border-b border-border flex justify-between items-center z-10">
      <h1 className="m-0 text-text text-lg">Map Builder — Floor {currentZLevel}</h1>
      <div className="flex gap-2 items-center">
        <Button variant="secondary" size="sm" onClick={() => onZoom(-0.25)} aria-label="Zoom out">−</Button>
        <span className="text-text-muted text-xs w-12 text-center">{Math.round(zoom * 100)}%</span>
        <Button variant="primary" size="sm" onClick={() => onZoom(0.25)} aria-label="Zoom in">+</Button>
        {onResetView && (
          <Button variant="outline" size="sm" onClick={onResetView} title="Reset zoom and pan to default">↺ Reset</Button>
        )}
        {onRelayout && (
          <Button variant="ghost" size="sm" onClick={onRelayout} title="Auto-relax overlapping nodes">✨</Button>
        )}
        {onCleanupOrphanExits && (
          <Button variant="ghost" size="sm" onClick={onCleanupOrphanExits} title="Remove exits pointing to deleted rooms">🧹</Button>
        )}
      </div>
    </div>
  )
}