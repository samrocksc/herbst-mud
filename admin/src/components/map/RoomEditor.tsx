/* eslint-disable functional/prefer-immutable-types, functional/no-mixed-types */
import { useState } from "react";
import { ALL_DIRECTIONS } from "./DirectionUtils";
import { Button } from "../Button";
import { SearchableSelect } from "../SearchableSelect";
import { FormField } from "../fields/FormField";
import { TextareaField } from "../fields/TextareaField";
import { FormError } from "../fields/FormError";
import { showToast } from "../Toast";
import { useRooms } from "../../hooks/useRooms";
import type { Room } from "./types";

type RoomEditorProps = {
  room: Room
  onCancel: () => void
}

export function RoomEditor({ room, onCancel }: RoomEditorProps) {
  const { rooms, updateRoom, isUpdating, createBidirectionalExit, removeBidirectionalExit } = useRooms();
  const [form, setForm] = useState({
    name: room.name,
    description: room.description,
    exits: { ...room.exits },
    isStartingRoom: room.isStartingRoom,
    isRootRoom: room.isRootRoom,
  });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const handleExitChange = async (dir: string, val: string) => {
    if (val) {
      const targetId = parseInt(val);
      const oldTargetId = form.exits[dir];

      if (oldTargetId && oldTargetId !== targetId) {
        removeBidirectionalExit({ roomId: room.id, direction: dir }).catch(() => {});
      }

      setForm(f => ({ ...f, exits: { ...f.exits, [dir]: targetId } }));
      try {
        await createBidirectionalExit({ roomId: room.id, direction: dir, targetRoomId: targetId });
      } catch {
        setForm(f => ({ ...f, exits: { ...f.exits, [dir]: oldTargetId } }));
        showToast("Failed to create exit", "error");
      }
    } else {
      const oldTargetId = form.exits[dir];
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { [dir]: _exitDir, ...rest } = form.exits;
      setForm(f => ({ ...f, exits: rest }));
      if (oldTargetId) {
        removeBidirectionalExit({ roomId: room.id, direction: dir }).catch(() => {});
      }
    }
  };

  const handleSave = () => {
    setError("");
    setSaving(true);
    updateRoom({
      id: room.id,
      update: {
        name: form.name,
        description: form.description,
        isStartingRoom: form.isStartingRoom,
        isRootRoom: form.isRootRoom,
        version: room.version,
      },
    });
    onCancel();
  };

  return (
    <>
      <div className="p-3 border-b border-border flex justify-between items-center">
        <h3 className="m-0 text-text text-base font-semibold">Edit Room</h3>
        <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button>
      </div>

      <div className="p-3 flex-1 overflow-y-auto">
        {error && <FormError message={error} />}
        <div className="mb-3">
          <FormField label="Name" value={form.name} onChange={(v) => setForm({ ...form, name: v })} />
        </div>
        <div className="mb-3">
          <TextareaField label="Description" value={form.description} onChange={(v) => setForm({ ...form, description: v })} rows={4} />
        </div>
        <div className="mb-3 space-y-2">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="isStartingRoom"
              checked={form.isStartingRoom}
              onChange={(e) => setForm(f => ({ ...f, isStartingRoom: e.target.checked }))}
              className="w-4 h-4 rounded border-border bg-surface text-primary focus:ring-primary"
            />
            <label htmlFor="isStartingRoom" className="text-text text-sm">Starting Room</label>
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="isRootRoom"
              checked={form.isRootRoom}
              onChange={(e) => setForm(f => ({ ...f, isRootRoom: e.target.checked }))}
              className="w-4 h-4 rounded border-border bg-surface text-primary focus:ring-primary"
            />
            <label htmlFor="isRootRoom" className="text-text text-sm">
              Root Room
              <span className="text-text-muted text-xs ml-1">(new characters spawn here)</span>
            </label>
          </div>
        </div>
        <div className="mb-3">
          <label className="text-text-muted text-xs block mb-2">Exits</label>
          {ALL_DIRECTIONS.map((dir) => (
            <div key={dir} className="flex items-center gap-2 mb-2">
              <span className="w-[60px] text-text-muted text-xs">{dir}:</span>
              <SearchableSelect
                options={rooms.map(r => ({ id: String(r.id), name: `${r.name} (ID: ${r.id})` }))}
                value={form.exits[dir] ? String(form.exits[dir]) : ""}
                onChange={(val) => handleExitChange(dir, val)}
                placeholder="Pick destination..."
              />
              {form.exits[dir] && (
                <Button variant="ghost" size="sm" className="!px-1 !py-0.5 text-danger"
                  onClick={() => handleExitChange(dir, "")} aria-label={`Remove ${dir} exit`}>×</Button>
              )}
            </div>
          ))}
        </div>
      </div>

      <div className="p-3 border-t border-border flex gap-2">
        <Button variant="primary" size="md" fullWidth onClick={handleSave} disabled={isUpdating || saving}>
          {isUpdating || saving ? "Saving..." : "Save Changes"}
        </Button>
        <Button variant="secondary" size="md" fullWidth onClick={onCancel}>Cancel</Button>
      </div>
    </>
  );
}