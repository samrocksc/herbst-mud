import { useState } from 'react';
import {
  useDialogNodes,
  useCreateDialogNode,
  useUpdateDialogNode,
  useDeleteDialogNode,
  type DialogNode,
  type DialogResponse,
} from '../hooks/useDialogNodes';
import { Button } from './Button';
import { FormField, TextareaField } from './FormFields';
import { DeleteConfirmation } from './DeleteConfirmation';

type Props = { npcTemplateId: string }

const EMPTY_RESPONSE: DialogResponse = {
  label: '', next_node_id: '', condition: '', quest_offer_id: '', decline_node_id: '', effects: [],
};

export function DialogNodesPanel({ npcTemplateId }: Props) {
  const { data: nodes, isLoading } = useDialogNodes(npcTemplateId);
  const createNode = useCreateDialogNode();
  const updateNode = useUpdateDialogNode();
  const deleteNode = useDeleteDialogNode();
  const [editing, setEditing] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);
  const [deleting, setDeleting] = useState<string | null>(null);

  const [newText, setNewText] = useState('');
  const [newId, setNewId] = useState('');
  const [newEntry, setNewEntry] = useState(false);
  const [editForm, setEditForm] = useState<Partial<DialogNode> | null>(null);

  const entryNode = (nodes ?? []).find((n) => n.is_entry);

  const handleAdd = async () => {
    const id = newId || `node_${Date.now()}`;
    await createNode.mutateAsync({
      id,
      npc_text: newText || '...',
      npc_template_id: npcTemplateId,
      is_entry: newEntry && !entryNode,
      responses: [],
      on_enter_effects: [],
    });
    setAdding(false);
    setNewText('');
    setNewId('');
    setNewEntry(false);
  };

  const startEdit = (node: DialogNode) => {
    setEditing(node.id);
    setEditForm({ ...node });
  };

  const handleSave = async () => {
    if (!editing || !editForm) return;
    await updateNode.mutateAsync({ id: editing, input: editForm });
    setEditing(null);
    setEditForm(null);
  };

  const handleDelete = async () => {
    if (!deleting) return;
    await deleteNode.mutateAsync(deleting);
    setDeleting(null);
  };

  if (isLoading) return <div className="text-text-muted text-sm">Loading dialog nodes...</div>;

  const sorted = [...(nodes ?? [])].sort((a, b) => {
    if (a.is_entry) return -1;
    if (b.is_entry) return 1;
    return a.id.localeCompare(b.id);
  });

  return (
    <div className="mt-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="m-0 text-text text-lg font-semibold">Dialog Tree</h2>
        <Button variant="primary" size="sm" onClick={() => setAdding(true)} disabled={adding}>
          + Add Node
        </Button>
      </div>

      {adding && (
        <div className="form-card mb-4 space-y-3">
          <FormField label="Node ID" value={newId} onChange={setNewId} placeholder="e.g. node_001_greeting" />
          <TextareaField label="NPC Text" value={newText} onChange={setNewText} rows={2} placeholder="What the NPC says..." />
          <div className="flex items-center gap-2">
            <input type="checkbox" checked={newEntry} onChange={(e) => setNewEntry(e.target.checked)} id="new-entry" disabled={!!entryNode} />
            <label htmlFor="new-entry" className="text-sm text-text">Entry Node</label>
          </div>
          <div className="flex gap-2">
            <Button variant="primary" size="sm" onClick={handleAdd} disabled={createNode.isPending}>
              {createNode.isPending ? 'Creating...' : 'Create'}
            </Button>
            <Button variant="secondary" size="sm" onClick={() => setAdding(false)}>Cancel</Button>
          </div>
        </div>
      )}

      {sorted.length === 0 && !adding && (
        <div className="text-text-muted text-sm py-4">No dialog nodes yet. Add one to start building the conversation tree.</div>
      )}

      {sorted.map((node) => (
        <div key={node.id} className="form-card mb-3">
          {editing === node.id && editForm ? (
            <div className="space-y-3">
              <div className="flex items-center gap-2 text-sm">
                <span className="font-mono text-primary">{node.id}</span>
                {node.is_entry && <span className="bg-success/20 text-success px-2 py-0.5 rounded text-xs">Entry</span>}
              </div>
              <TextareaField label="NPC Text" value={editForm.npc_text ?? ''} onChange={(v) => setEditForm({ ...editForm, npc_text: v })} rows={3} />
              <div className="flex items-center gap-2">
                <input type="checkbox" checked={editForm.is_entry ?? false} onChange={(e) => setEditForm({ ...editForm, is_entry: e.target.checked })} id={`entry-${node.id}`} disabled={!!entryNode && entryNode.id !== node.id} />
                <label htmlFor={`entry-${node.id}`} className="text-sm text-text">Entry Node</label>
              </div>
              <FormField label="Entry Condition" value={editForm.entry_condition ?? ''} onChange={(v) => setEditForm({ ...editForm, entry_condition: v })} placeholder="e.g. character.tags.has(wizard_complete)" />
              <ResponsesEditor responses={editForm.responses ?? []} onChange={(r) => setEditForm({ ...editForm, responses: r })} />
              <div className="flex gap-2">
                <Button variant="primary" size="sm" onClick={handleSave} disabled={updateNode.isPending}>
                  {updateNode.isPending ? 'Saving...' : 'Save'}
                </Button>
                <Button variant="secondary" size="sm" onClick={() => { setEditing(null); setEditForm(null); }}>Cancel</Button>
              </div>
            </div>
          ) : (
            <div>
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2 text-sm">
                  <span className="font-mono text-primary font-bold">{node.id}</span>
                  {node.is_entry && <span className="bg-success/20 text-success px-2 py-0.5 rounded text-xs">Entry</span>}
                </div>
                <div className="flex gap-2">
                  <Button variant="ghost" size="sm" onClick={() => startEdit(node)}>Edit</Button>
                  <Button variant="danger" size="sm" onClick={() => setDeleting(node.id)}>Delete</Button>
                </div>
              </div>
              <p className="text-text text-sm mb-2 mt-0">"{node.npc_text}"</p>
              {node.responses && node.responses.length > 0 && (
                <div className="text-sm">
                  <span className="text-text-muted">Responses:</span>
                  <ol className="ml-4 mt-1 space-y-0.5">
                    {node.responses.map((r, i) => (
                      <li key={i} className="text-text">
                        <span className="text-primary">{i + 1}.</span> {r.label || '[Leave]'}
                        {r.next_node_id && <span className="text-text-muted"> → {r.next_node_id}</span>}
                        {r.quest_offer_id && <span className="text-success ml-1">[Quest: {r.quest_offer_id}]</span>}
                      </li>
                    ))}
                  </ol>
                </div>
              )}
              {(!node.responses || node.responses.length === 0) && (
                <div className="text-text-muted text-xs">No responses (end of conversation)</div>
              )}
            </div>
          )}
        </div>
      ))}

      <DeleteConfirmation
        open={!!deleting}
        title="Delete Dialog Node"
        message="Are you sure? This will permanently delete this dialog node and break any references to it."
        onConfirm={handleDelete}
        onCancel={() => setDeleting(null)}
        isLoading={deleteNode.isPending}
      />
    </div>
  );
}

function ResponsesEditor({ responses, onChange }: { responses: DialogResponse[]; onChange: (r: DialogResponse[]) => void }) {
  const add = () => onChange([...responses, { ...EMPTY_RESPONSE }]);
  const remove = (i: number) => onChange(responses.filter((_, idx) => idx !== i));
  const update = (i: number, patch: Partial<DialogResponse>) => onChange(responses.map((r, idx) => idx === i ? { ...r, ...patch } : r));

  return (
    <div>
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-semibold text-text">Responses</span>
        <Button variant="ghost" size="sm" onClick={add}>+ Response</Button>
      </div>
      {responses.map((r, i) => (
        <div key={i} className="grid grid-cols-4 gap-2 mb-2 items-end">
          <FormField label="Label" value={r.label} onChange={(v) => update(i, { label: v })} placeholder="What troubles you?" />
          <FormField label="Next Node" value={r.next_node_id} onChange={(v) => update(i, { next_node_id: v })} placeholder="node_002" />
          <FormField label="Quest Offer" value={r.quest_offer_id ?? ''} onChange={(v) => update(i, { quest_offer_id: v })} placeholder="quest ID" />
          <Button variant="danger" size="sm" onClick={() => remove(i)}>×</Button>
        </div>
      ))}
      {responses.length === 0 && <div className="text-text-muted text-xs">No responses. Player will see [Leave] only.</div>}
    </div>
  );
}