// vitest.setup.ts - Global mocks for tests
import { vi } from 'vitest'

// Mock @xyflow/react module - must be hoisted
vi.mock('@xyflow/react', () => ({
  ReactFlowProvider: ({ children }: { children: React.ReactNode }) => children,
  ReactFlow: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="reactflow">{children}</div>
  ),
  Controls: () => <div data-testid="controls" />,
  MiniMap: ({ nodeColor }: { nodeColor?: string }) => <div data-testid="minimap" />,
  Background: ({ color, gap }: { color?: string; gap?: number }) => <div data-testid="background" />,
  Handle: ({ type, position, id, style }: any) => <div data-testid={`handle-${id}`} data-position={position} data-type={type} />,
  Position: {
    Left: 'left',
    Right: 'right',
    Top: 'top',
    Bottom: 'bottom',
  },
  useNodesState: (initial: any[]) => [initial, vi.fn(), vi.fn()],
  useEdgesState: (initial: any[]) => [initial, vi.fn(), vi.fn()],
  addEdge: vi.fn((edge: any, edges: any[]) => [...edges, edge]),
  Node: 'node',
  Edge: 'edge',
}))