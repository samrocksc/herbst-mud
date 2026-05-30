import type { Trigger } from "../../hooks/useTriggers";

export function TriggerDetailView({ trigger }: Readonly<{ trigger: Trigger }>) {
  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Trigger Details</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <DetailField label="ID" value={String(trigger.id)} />
        <DetailField label="Name" value={trigger.name} />
        <DetailField label="World" value={trigger.world_id} />
        <DetailField label="Trigger Type" value={trigger.trigger_type} />
        <DetailField label="Target Type" value={trigger.target_type} />
        <DetailField label="Target ID" value={String(trigger.target_id)} />
        <DetailField label="Room ID" value={trigger.room_id != null ? String(trigger.room_id) : "—"} />
        <DetailField label="Equipment ID" value={trigger.equipment_id != null ? String(trigger.equipment_id) : "—"} />
        <DetailField
          label="Enabled"
          value={trigger.enabled ? "Yes" : "No"}
        />
      </div>
      {trigger.condition && (
        <div className="mt-4">
          <h3 className="text-text-muted text-xs uppercase font-semibold mb-2">Condition</h3>
          <code className="block p-3 bg-surface rounded border border-border text-sm">
            {trigger.condition}
          </code>
        </div>
      )}
    </div>
  );
}

function DetailField({ label, value }: Readonly<{ label: string; value: string }>) {
  return (
    <div>
      <span className="text-text-muted text-xs block mb-0.5">{label}</span>
      <span className="text-text text-sm font-medium">{value}</span>
    </div>
  );
}
