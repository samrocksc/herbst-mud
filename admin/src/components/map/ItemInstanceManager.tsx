 
import { Link } from "@tanstack/react-router";
import { Button } from "../Button";
import { ItemEditRow } from "./ItemEditRow";
import { ItemInstanceRow } from "./ItemInstanceRow";
import { useItemInstances } from "./useItemInstances";
import type { ItemInstanceView } from "./types";

type ItemInstanceManagerProps = Readonly<{ roomId: number }>

export function ItemInstanceManager({ roomId }: ItemInstanceManagerProps) {
  const {
    instancesQuery, updateMutation,
    editingId, setEditingId,
    confirmDeleteId, setConfirmDeleteId,
    editForm, setEditForm,
    handleUpdate, handleDelete,
  } = useItemInstances(roomId);

  const instances = instancesQuery.data ?? [];
  const startEdit = (inst: ItemInstanceView) => {
    setEditingId(inst.id); setConfirmDeleteId(null);
    setEditForm({ name: inst.name, description: inst.description, slot: inst.slot, level: inst.level, weight: inst.weight, color: inst.color });
  };

  if (instancesQuery.isLoading) return (
    <div className="mb-3">
      <strong className="text-success text-xs">Items:</strong>
      <div className="text-text-muted text-[10px] mt-1">Loading...</div>
    </div>
  );

  if (instancesQuery.error) return (
    <div className="mb-3">
      <strong className="text-success text-xs">Items:</strong>
      <div className="text-danger text-[10px] mt-1">Error loading items</div>
    </div>
  );

  return (
    <div className="mb-3">
      <div className="flex items-center justify-between mb-1">
        <strong className="text-success text-xs">Items:</strong>
        <Link to="/map/rooms/$roomId/items/spawn" params={{ roomId: String(roomId) }} className="no-underline">
          <Button variant="primary" size="sm" className="!px-1.5 !py-0 !text-[10px]">+ Add Instance</Button>
        </Link>
      </div>
      {instances.length === 0 ? (
        <div className="text-text-muted text-[10px]">No item instances in this room.</div>
      ) : (
        <div className="mt-1 flex flex-col gap-1">
          {instances.map((inst) => (
            <div key={inst.id} className="p-1 bg-surface-muted rounded text-xs text-text">
              {editingId === inst.id ? (
                <ItemEditRow inst={inst} editForm={editForm} setEditForm={setEditForm}
                  onSave={handleUpdate} onCancel={() => { setEditingId(null); setEditForm({}); }}
                  isPending={updateMutation.isPending} error={updateMutation.error as Error | null} />
              ) : (
                <ItemInstanceRow inst={inst} confirmDeleteId={confirmDeleteId}
                  onEdit={() => startEdit(inst)}
                  onDelete={() => confirmDeleteId === inst.id ? handleDelete(inst.id) : setConfirmDeleteId(inst.id)} />
              )}
              {confirmDeleteId === inst.id && editingId !== inst.id && (
                <div className="mt-1 text-[10px] text-text">Confirm delete?{" "}
                  <Button variant="danger" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={() => handleDelete(inst.id)}>Yes</Button>{" "}
                  <Button variant="ghost" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={() => setConfirmDeleteId(null)}>No</Button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
