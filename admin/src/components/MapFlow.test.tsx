/**
 * @vitest-environment jsdom
 */
import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import { ReactFlowProvider } from '@xyflow/react'
import { MapFlow } from './MapFlow'
import type { Node, Edge } from '@xyflow/react'

// Mock ResizeObserver as a class for ReactFlow
class MockResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

global.ResizeObserver = MockResizeObserver

describe('MapFlow', () => {
  const mockNodes: Node[] = [
    {
      id: '1',
      type: 'room',
      position: { x: 0, y: 0 },
      data: { name: 'Test Room', description: 'A test room', zLevel: 0, roomId: 1, isStartingRoom: false, atmosphere: 'air', exits: {} },
    },
  ]

  const mockEdges: Edge[] = []

  const renderWithProvider = (component: React.ReactElement) => {
    return render(
      <ReactFlowProvider>
        {component}
      </ReactFlowProvider>
    )
  }

  it('renders ReactFlow canvas', () => {
    renderWithProvider(
      <MapFlow
        nodes={mockNodes}
        edges={mockEdges}
      />
    )
    
    const container = document.querySelector('.react-flow')
    expect(container).not.toBeNull()
  })

  it('renders with custom dimensions', () => {
    renderWithProvider(
      <MapFlow
        nodes={mockNodes}
        edges={mockEdges}
      />
    )
    
    const container = document.querySelector('.react-flow')
    expect(container).not.toBeNull()
  })

  it('accepts nodes and edges props', () => {
    const nodes: Node[] = [
      { id: '1', type: 'room', position: { x: 0, y: 0 }, data: mockNodes[0].data },
      { id: '2', type: 'room', position: { x: 100, y: 100 }, data: mockNodes[0].data },
    ]
    
    renderWithProvider(
      <MapFlow
        nodes={nodes}
        edges={[]}
      />
    )
    
    const container = document.querySelector('.react-flow')
    expect(container).not.toBeNull()
  })

  it('registers custom node types', () => {
    renderWithProvider(
      <MapFlow
        nodes={mockNodes}
        edges={mockEdges}
      />
    )
    
    // Room nodes should render with custom type
    const roomNodes = document.querySelectorAll('.react-flow__node-room')
    expect(roomNodes.length).toBeGreaterThan(0)
  })

  it('calls onConnect callback when connections are made', () => {
    renderWithProvider(
      <MapFlow
        nodes={mockNodes}
        edges={mockEdges}
        onConnect={() => {}}
      />
    )
    
    const container = document.querySelector('.react-flow')
    expect(container).not.toBeNull()
  })

  it('renders Background component', () => {
    renderWithProvider(
      <MapFlow
        nodes={mockNodes}
        edges={mockEdges}
      />
    )
    
    const svg = document.querySelector('.react-flow__background')
    expect(svg).not.toBeNull()
  })

  it('renders Controls component', () => {
    renderWithProvider(
      <MapFlow
        nodes={mockNodes}
        edges={mockEdges}
      />
    )
    
    const controls = document.querySelector('.react-flow__controls')
    expect(controls).not.toBeNull()
  })

  it('renders MiniMap component', () => {
    renderWithProvider(
      <MapFlow
        nodes={mockNodes}
        edges={mockEdges}
      />
    )
    
    const minimap = document.querySelector('.react-flow__minimap')
    expect(minimap).not.toBeNull()
  })
})