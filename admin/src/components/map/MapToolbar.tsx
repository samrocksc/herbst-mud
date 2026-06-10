/* eslint-disable functional/no-mixed-types */
import { Button } from "../Button";
import { Tooltip } from "../Tooltip";

type MapToolbarProps = Readonly<{
  currentZLevel: number
  zLevels: ReadonlyArray<number>
  zoom: number
  onZoom: (delta: number) => void
  onResetView?: () => void
  onRelayout?: () => void
  onCleanupOrphanExits?: () => void
  onGoToFloor: (z: number) => void
  onAddFloor: () => void
  onDeleteFloor: () => void
}>

export function MapToolbar({
  currentZLevel,
  zLevels,
  zoom,
  onZoom,
  onResetView,
  onRelayout,
  onCleanupOrphanExits,
  onGoToFloor,
  onAddFloor,
  onDeleteFloor,
}: MapToolbarProps) {
  const minZ = zLevels.length > 0 ? zLevels[0] : 0
  const maxZ = zLevels.length > 0 ? zLevels[zLevels.length - 1] : 0
  const hasCurrentFloor = zLevels.includes(currentZLevel)

  return (
    <div className="absolute top-0 left-0 right-0 p-3 bg-surface-muted border-b border-border flex flex-wrap justify-between items-center gap-2 z-10">
      <h1 className="m-0 text-text text-lg">Map Builder — Floor {currentZLevel}</h1>

      <div className="flex gap-2 items-center flex-wrap">
        {/* Z-level controls */}
        <div className="flex gap-1 items-center bg-surface rounded border border-border px-2 py-1">
          <span className="text-text-muted text-xs mr-1">Z:</span>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onGoToFloor(currentZLevel - 1)}
            disabled={currentZLevel <= minZ}
            title={`Go down to floor ${currentZLevel - 1}`}
            aria-label="Go down a floor"
          >
            ▼
          </Button>
          <select
            value={hasCurrentFloor ? currentZLevel : ""}
            onChange={(e) => onGoToFloor(Number(e.target.value))}
            className="bg-transparent text-text text-xs border-none focus:outline-none"
            title="Jump to floor"
            aria-label="Jump to floor"
          >
            {!hasCurrentFloor && <option value="">{currentZLevel} (empty)</option>}
            {zLevels.map((z) => (
              <option key={z} value={z}>
                {z === 0 ? "G" : z > 0 ? `+${z}` : `${z}`}
              </option>
            ))}
          </select>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onGoToFloor(currentZLevel + 1)}
            title={`Add a new floor at ${currentZLevel + 1}`}
            aria-label="Add floor above"
          >
            ▲
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={onAddFloor}
            title="Add empty floor (jump to it)"
            aria-label="Add empty floor"
          >
            +
          </Button>
          <Tooltip content="Delete all rooms on the current floor">
            <Button
              variant="ghost"
              size="sm"
              onClick={onDeleteFloor}
              disabled={zLevels.length <= 1}
              aria-label="Delete current floor"
            >
              🗑
            </Button>
          </Tooltip>
        </div>

        {/* Zoom controls */}
        <Button variant="secondary" size="sm" onClick={() => onZoom(-0.25)} aria-label="Zoom out">−</Button>
        <span className="text-text-muted text-xs w-12 text-center">{Math.round(zoom * 100)}%</span>
        <Button variant="primary" size="sm" onClick={() => onZoom(0.25)} aria-label="Zoom in">+</Button>

        {onResetView && (
          <Tooltip content="Reset zoom and pan to default">
            <Button
              variant="outline"
              size="sm"
              onClick={onResetView}
              aria-label="Reset zoom and pan to default"
            >
              ↺ Reset
            </Button>
          </Tooltip>
        )}
        {onRelayout && (
          <Tooltip content="Auto-relax overlapping nodes">
            <Button
              variant="ghost"
              size="sm"
              onClick={onRelayout}
              aria-label="Auto-relax overlapping nodes"
            >
              ✨
            </Button>
          </Tooltip>
        )}
        {onCleanupOrphanExits && (
          <Tooltip content="Remove exits pointing to deleted rooms">
            <Button
              variant="ghost"
              size="sm"
              onClick={onCleanupOrphanExits}
              aria-label="Remove exits pointing to deleted rooms"
            >
              🧹
            </Button>
          </Tooltip>
        )}
      </div>
    </div>
  );
}
