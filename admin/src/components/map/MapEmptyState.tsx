/* eslint-disable functional/no-mixed-types */
import { Button } from "../Button";

type MapEmptyStateProps = Readonly<{
  currentZLevel: number
  onCreateRoom: () => void
}>

export function MapEmptyState({ currentZLevel, onCreateRoom }: MapEmptyStateProps) {
  return (
    <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
      <div className="bg-surface border border-border rounded-lg p-6 max-w-[400px] text-center shadow-lg pointer-events-auto">
        <div className="text-4xl mb-2">🗺️</div>
        <h3 className="text-text text-lg font-semibold mb-2">This floor is empty</h3>
        <p className="text-text-muted text-sm mb-4">
          Floor {currentZLevel} has no rooms yet. Create the first room to start building this level.
        </p>
        <Button variant="primary" onClick={onCreateRoom}>+ Create First Room</Button>
      </div>
    </div>
  );
}
