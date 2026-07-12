import { useMemo } from "react";
import { ReactFlow, type Node, type Edge, MiniMap, Controls, Background } from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { RoomNode } from "./RoomNode";

interface ReactFlowCanvasProps {
  nodes: Node[];
  edges: Edge[];
}

export function ReactFlowCanvas({ nodes, edges }: ReactFlowCanvasProps) {
  const memoizedNodes = useMemo(() => nodes, [nodes]);
  const memoizedEdges = useMemo(() => edges, [edges]);

  return (
    <ReactFlow
      nodes={memoizedNodes}
      edges={memoizedEdges}
      nodeTypes={{
        room: ({ data, selected, position }) => (
          <RoomNode
            {...data}
            pos={position}
            isSelected={selected}
            zoom={1} // Simplified, ideally from a store or context
            isDragging={false} // React Flow handles dragging
            onDragStart={() => {}}
            onDragEnd={(id, x, y) => {}}
          />
        ),
      }}
      fitView
      proOptions={{ hideAttribution: true }}
      minZoom={0.5}
      maxZoom={2}
    >
      <Background />
      <MiniMap />
      <Controls />
    </ReactFlow>
  );
}
