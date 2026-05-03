import { Button } from '../Button'

type MapToolbarProps = {
  currentZLevel: number
  zoom: number
  setZoom: (zoom: number | ((z: number) => number)) => void
}

export function MapToolbar({ currentZLevel, zoom, setZoom }: MapToolbarProps) {
  return (
    <div className="absolute top-0 left-0 right-0 p-3 bg-surface-muted border-b border-border flex justify-between items-center z-10">
      <h1 className="m-0 text-text text-lg">Map Builder — Floor {currentZLevel}</h1>
      <div className="flex gap-2 items-center">
        <Button
          variant="secondary"
          size="sm"
          onClick={() => setZoom((z) => Math.max(z - 0.25, 0.5))}
        >
          −
        </Button>
        <span className="text-text-muted text-xs w-12 text-center">
          {Math.round(zoom * 100)}%
        </span>
        <Button
          variant="primary"
          size="sm"
          onClick={() => setZoom((z) => Math.min(z + 0.25, 2))}
        >
          +
        </Button>
      </div>
    </div>
  )
}
