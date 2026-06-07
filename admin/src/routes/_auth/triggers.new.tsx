import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useCreateTrigger,
  type TriggerInput,
} from "../../hooks/useTriggers";
import { PageHeader } from "../..//components/PageHeader";
import { Button } from "../../components/Button";
import {
  FormField,
  TextareaField,
  SelectField,
  CheckboxField,
  NumberField,
} from "../../components/FormFields";
import { showToast } from "../../components/Toast";
import { PageContainer } from "../../components/PageContainer";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import { useWorlds } from "../../hooks/useWorlds";
import { ResourceIdField } from "../../components/ResourceIdField";
import { useNPCTemplates } from "../../hooks/useNPCTemplates";
import { useDialogNodes, useCreateDialogNode, type DialogNodeInput } from "../../hooks/useDialogNodes";
import { SearchableSelect } from "../../components/SearchableSelect";
import { DialogTreeEditor } from "../../components/DialogTreeEditor";

export const Route = createFileRoute("/_auth/triggers/new")({
  component: CreateTriggerPage,
});

const TRIGGER_TYPE_OPTS = [
  { value: "use", label: "Use" },
  { value: "touch", label: "Touch" },
  { value: "press", label: "Press" },
  { value: "enter", label: "Enter Room" },
  { value: "examine", label: "Examine" },
  { value: "talk", label: "Talk" },
];

const TARGET_TYPE_OPTS = [
  { value: "recipe", label: "Recipe" },
  { value: "effect", label: "Effect" },
  { value: "dialog_node", label: "Dialog Node" },
];

export function CreateTriggerPage() {
  const navigate = useNavigate();
  const createTrigger = useCreateTrigger();
  const createDialogNode = useCreateDialogNode();
  const { currentWorld } = useWorldStore();
  const { data: worlds } = useWorlds();
  const { data: npcTemplates } = useNPCTemplates();

  const [formData, setFormData] = useState<TriggerInput>({
    name: "",
    world_id: currentWorld || "default",
    trigger_type: "use",
    target_type: "recipe",
    target_id: "",
    room_id: null,
    equipment_id: null,
    condition: "",
    enabled: true,
    examine_weight: 0,
  });

  const [talkMode, setTalkMode] = useState<"existing" | "new">("existing");
  const [selectedNpcTemplateId, setSelectedNpcTemplateId] = useState<string>("");
  const [selectedEntryNodeId, setSelectedEntryNodeId] = useState<string>("");
  const [createdNodes, setCreatedNodes] = useState<DialogNodeInput[]>([]);
  const [submitting, setSubmitting] = useState(false);

  const set = (patch: Partial<TriggerInput>) => setFormData((prev) => ({ ...prev, ...patch }));

  const isTalk = formData.trigger_type === "talk";

  const { data: dialogNodes } = useDialogNodes(selectedNpcTemplateId);

  const entryNodes = (dialogNodes ?? []).filter((n) => n.is_entry);

  const npcOptions = (npcTemplates ?? []).map((t) => ({ id: t.id, name: t.name }));

  const entryNodeOptions = entryNodes.map((n) => ({ id: n.id, name: `${n.id} — "${n.npc_text.slice(0, 30)}"` }));

  const handleApplyTalkMode = (): string | null => {
    if (talkMode === "existing") {
      if (!selectedEntryNodeId) {
        showToast("Select an entry dialog node", "error");
        return null;
      }
      return selectedEntryNodeId;
    }
    if (createdNodes.length === 0) {
      showToast("Create at least one entry node in the dialog tree", "error");
      return null;
    }
    const entryNode = createdNodes.find((n) => n.is_entry);
    if (!entryNode) {
      showToast("Mark one node as the entry node", "error");
      return null;
    }
    return entryNode.id || null;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      let targetId = formData.target_id;

      if (isTalk) {
        if (talkMode === "new") {
          for (const nodeInput of createdNodes) {
            await createDialogNode.mutateAsync({
              ...nodeInput,
              npc_template_id: selectedNpcTemplateId,
            } as DialogNodeInput & { npc_template_id: string });
          }
          const entryNode = createdNodes.find((n) => n.is_entry);
          targetId = entryNode?.id || "";
        } else {
          targetId = selectedEntryNodeId;
        }
      }

      const payload: TriggerInput = {
        ...formData,
        target_type: isTalk ? "dialog_node" : formData.target_type,
        target_id: isTalk ? targetId : formData.target_id,
      };

      await createTrigger.mutateAsync(payload);
      showToast("Trigger created", "success");
      navigate({ to: "/triggers" });
    } catch (err) {
      console.error("Trigger creation error:", err);
      const message = err instanceof Error ? err.message : "Failed to create trigger";
      showToast(message, "error");
    } finally {
      setSubmitting(false);
    }
  };

  const isPending = createTrigger.isPending || submitting;

  return (
    <PageContainer>
      <PageHeader title="Create Trigger" showBack backTo="/triggers" />
      <div className="card bg-surface p-6 border border-border rounded">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-text font-semibold mb-4">Basic Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField label="Name *" value={formData.name} onChange={(v) => set({ name: v })} required />
            <SelectField
              label="World *"
              value={formData.world_id}
              onChange={(v) => set({ world_id: v })}
              options={(worlds || []).map(w => ({ value: String(w.id), label: w.name }))}
            />
            <SelectField label="Trigger Type" value={formData.trigger_type} onChange={(v) => set({ trigger_type: v })} options={TRIGGER_TYPE_OPTS} />
            {!isTalk && (
              <SelectField label="Target Type" value={formData.target_type} onChange={(v) => set({ target_type: v })} options={TARGET_TYPE_OPTS} />
            )}
            {!isTalk && (
              <ResourceIdField
                label="Target ID *"
                value={formData.target_id}
                onChange={(v) => set({ target_id: v != null ? String(v) : "" })}
                resourceType="targets"
                apiBase=""
                worldId={currentWorld}
              />
            )}
          </div>

          {isTalk && (
            <div className="space-y-4 mt-4 p-4 bg-surface-muted rounded border border-border">
              <h3 className="text-text font-semibold mt-0 mb-2">Talk Trigger Configuration</h3>

              <div>
                <label className="text-text-muted text-xs block mb-1">NPC Template</label>
                <SearchableSelect
                  options={npcOptions}
                  value={selectedNpcTemplateId}
                  onChange={setSelectedNpcTemplateId}
                  placeholder="Select an NPC template..."
                />
              </div>

              {selectedNpcTemplateId && (
                <>
                  <div className="flex gap-4 items-center">
                    <label className="text-sm text-text">
                      <input
                        type="radio"
                        name="talkMode"
                        value="existing"
                        checked={talkMode === "existing"}
                        onChange={() => setTalkMode("existing")}
                        className="mr-1"
                      />
                      Link existing entry node
                    </label>
                    <label className="text-sm text-text">
                      <input
                        type="radio"
                        name="talkMode"
                        value="new"
                        checked={talkMode === "new"}
                        onChange={() => setTalkMode("new")}
                        className="mr-1"
                      />
                      Build new dialog tree
                    </label>
                  </div>

                  {talkMode === "existing" && (
                    <div>
                      <label className="text-text-muted text-xs block mb-1">Entry Dialog Node</label>
                      <SearchableSelect
                        options={entryNodeOptions}
                        value={selectedEntryNodeId}
                        onChange={setSelectedEntryNodeId}
                        placeholder="Select an entry node..."
                      />
                      {entryNodes.length === 0 && (
                        <p className="text-text-muted text-xs mt-1">No entry nodes found for this NPC template.</p>
                      )}
                    </div>
                  )}

                  {talkMode === "new" && (
                    <NewDialogTreeEditor
                      npcTemplateId={selectedNpcTemplateId}
                      createdNodes={createdNodes}
                      onNodesChange={setCreatedNodes}
                    />
                  )}
                </>
              )}
            </div>
          )}

          <h3 className="text-text font-semibold mt-6 mb-4">Target Object</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <ResourceIdField
              label="Room ID (optional)"
              value={formData.room_id ?? ""}
              onChange={(v) => set({ room_id: v === null || v === "" ? null : Number(v) > 0 ? Number(v) : null })}
              resourceType="rooms"
              apiBase={window.location.origin}
            />
            <ResourceIdField
              label="Equipment ID (optional)"
              value={formData.equipment_id ?? ""}
              onChange={(v) => set({ equipment_id: v === null || v === "" ? null : Number(v) > 0 ? Number(v) : null })}
              resourceType="equipment-templates"
              apiBase={window.location.origin}
            />
          </div>

          <h3 className="text-text font-semibold mt-6 mb-4">Conditions & Settings</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <NumberField
              label="Examine Weight"
              value={formData.examine_weight ?? 0}
              onChange={(v) => set({ examine_weight: v })}
              placeholder="0 = always fires"
              tooltip="Examine tier threshold - higher values require more examine skill to reveal. 0 = always visible."
            />
            <div>
              <TextareaField label="Condition (SPICE expression, optional)" value={formData.condition} onChange={(v) => set({ condition: v })} rows={3} placeholder="e.g., player_level >= 10" />
              <div className="mt-1">
                <a href="/docs/triggers" target="_blank" rel="noopener noreferrer" className="text-xs text-primary hover:underline">
                  Learn more about conditions
                </a>
              </div>
            </div>
            <div className="flex items-end">
              <CheckboxField label="Enabled" checked={formData.enabled} onChange={(v) => set({ enabled: v })} />
            </div>
          </div>

          <div className="flex gap-2 justify-end mt-6">
            <Button variant="secondary" onClick={() => navigate({ to: "/triggers" })}>Cancel</Button>
            <Button variant="primary" type="submit" disabled={isPending}>
              {isPending ? "Creating..." : "Create Trigger"}
            </Button>
          </div>
        </form>
      </div>
    </PageContainer>
  );
}

/** Local editor for creating new dialog nodes before submit */
function NewDialogTreeEditor({
  npcTemplateId,
  createdNodes,
  onNodesChange,
}: Readonly<{
  npcTemplateId: string;
  createdNodes: DialogNodeInput[];
  onNodesChange: (nodes: DialogNodeInput[]) => void;
}>) {
  const [newId, setNewId] = useState("");
  const [newText, setNewText] = useState("");
  const [newEntry, setNewEntry] = useState(false);

  const hasEntry = createdNodes.some((n) => n.is_entry);

  const handleAddNode = () => {
    const id = newId || `node_${Date.now()}`;
    const node: DialogNodeInput = {
      id,
      npc_text: newText || "...",
      is_entry: newEntry && !hasEntry,
      responses: [],
      on_enter_effects: [],
    };
    onNodesChange([...createdNodes, node]);
    setNewId("");
    setNewText("");
    setNewEntry(false);
  };

  const handleRemoveNode = (idx: number) => {
    onNodesChange(createdNodes.filter((_, i) => i !== idx));
  };

  return (
    <div>
      <h4 className="text-text font-semibold mb-2">New Dialog Nodes</h4>

      <div className="space-y-3 mb-4">
        {createdNodes.map((node, idx) => (
          <div key={node.id || idx} className="form-card p-3">
            <div className="flex items-center justify-between mb-1">
              <div className="flex items-center gap-2 text-sm">
                <span className="font-mono text-primary font-bold">{node.id}</span>
                {node.is_entry && <span className="bg-success/20 text-success px-2 py-0.5 rounded text-xs">Entry</span>}
              </div>
              <Button variant="danger" size="sm" onClick={() => handleRemoveNode(idx)}>Remove</Button>
            </div>
            <p className="text-text text-sm mt-0">&quot;{node.npc_text}&quot;</p>
          </div>
        ))}
      </div>

      <div className="form-card space-y-3">
        <FormField label="Node ID" value={newId} onChange={setNewId} placeholder="e.g. node_001_greeting" />
        <TextareaField label="NPC Text" value={newText} onChange={setNewText} rows={2} placeholder="What the NPC says..." />
        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={newEntry}
            onChange={(e) => setNewEntry(e.target.checked)}
            id="talk-new-entry"
            disabled={hasEntry}
          />
          <label htmlFor="talk-new-entry" className="text-sm text-text">Entry Node</label>
        </div>
        <Button variant="primary" size="sm" onClick={handleAddNode} disabled={!newText && !newId}>
          Add Node
        </Button>
      </div>
    </div>
  );
}