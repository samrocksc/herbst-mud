/* eslint-disable functional/no-mixed-types, functional/prefer-immutable-types */
import { useState, useEffect, useRef } from "react";
import { Modal } from "./Modal";
import { Button } from "./Button";

type DialogResponse = Readonly<{
  text: string;
  next_node_id?: string;
  effect_ids?: number[];
  quest_offer_id?: string;
}>;

type DialogNode = Readonly<{
  id: string;
  npc_text: string;
  responses: DialogResponse[];
  is_entry: boolean;
  order?: number;
}>;

type NPCConversationModalProps = Readonly<{
  isOpen: boolean;
  onClose: () => void;
  npcName: string;
  npcTemplateId: string;
  onConversationEnd?: () => void;
}>;

export function NPCConversationModal({
  isOpen,
  onClose,
  npcName,
  npcTemplateId,
  onConversationEnd,
}: NPCConversationModalProps) {
  const [currentNodeId, setCurrentNodeId] = useState<string | null>(null);
  const [nodes, setNodes] = useState<Record<string, DialogNode>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const initialLoadRef = useRef(false);

  // Load dialog nodes when modal opens
  useEffect(() => {
    if (isOpen && !initialLoadRef.current) {
      initialLoadRef.current = true;
      loadDialogNodes();
    }
  }, [isOpen]);

  const loadDialogNodes = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`${window.location.origin}/api/npc-templates/${npcTemplateId}/dialog-nodes`);
      if (!response.ok) {
        throw new Error(`Failed to load dialog nodes: ${response.status} ${response.statusText}`);
      }
      const data = await response.json();
      const nodeMap: Record<string, DialogNode> = {};
      for (const node of data.nodes) {
        nodeMap[node.id] = node;
      }
      setNodes(nodeMap);

      // Find entry node
      const entryNode = data.nodes.find((n: DialogNode) => n.is_entry);
      if (entryNode) {
        setCurrentNodeId(entryNode.id);
      } else if (data.nodes.length > 0) {
        // Fallback to first node
        setCurrentNodeId(data.nodes[0].id);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load dialog nodes");
    } finally {
      setIsLoading(false);
    }
  };

  const handleResponse = (response: DialogResponse) => {
    if (response.effect_ids && response.effect_ids.length > 0) {
      // Trigger effects via event hooks
      // For now, just log - effects would be handled by backend event system
      console.log("Triggering effects:", response.effect_ids);
    }

    if (response.next_node_id) {
      setCurrentNodeId(response.next_node_id);
    } else {
      // Conversation ended
      onClose();
      if (onConversationEnd) {
        onConversationEnd();
      }
    }
  };

  const current_node = currentNodeId ? nodes[currentNodeId] : null;

  if (!isOpen) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Talk to ${npcName}`}>
      {isLoading ? (
        <div className="p-4 text-text-muted text-sm">Loading conversation...</div>
      ) : error ? (
        <div className="p-4 bg-danger/10 border border-danger rounded text-danger text-sm">
          Error: {error}
        </div>
      ) : !current_node ? (
        <div className="p-4 text-text-muted text-sm">
          This NPC has no conversation available.
        </div>
      ) : (
        <div className="space-y-4">
          {/* NPC Speaker Line */}
          <div className="bg-surface-muted p-4 rounded border border-border">
            <div className="font-semibold text-primary mb-2">{npcName}</div>
            <div className="text-text text-sm italic">"{current_node.npc_text}"</div>
          </div>

          {/* Player Responses */}
          {current_node.responses && current_node.responses.length > 0 ? (
            <div className="space-y-2">
              <div className="text-text-muted text-xs">Your responses:</div>
              {current_node.responses.map((response, idx) => (
                <Button
                  key={idx}
                  variant="secondary"
                  size="md"
                  fullWidth
                  onClick={() => handleResponse(response)}
                  className="text-left justify-start"
                >
                  {response.text}
                </Button>
              ))}
            </div>
          ) : (
            <div className="text-text-muted text-sm text-center py-4">
              Conversation complete.
            </div>
          )}
        </div>
      )}
    </Modal>
  );
}
