import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useCallback } from 'react'
import { useMapState } from '../hooks/useMapState'
import { MapSidebar } from '../components/map/MapSidebar'
import { MapToolbar } from '../components/map/MapToolbar'
import { MapCanvas } from '../components/map/MapCanvas'
import { RoomDetailPanel } from '../components/map/RoomDetailPanel'
import { RoomEditor } from '../components/map/RoomEditor'
import { CreateRoomModal } from '../components/map/CreateRoomModal'
import { MenuIcon } from '../components/icons/MenuIcon'
import { Button } from '../components/Button'

export const Route = createFileRoute('/map')({
  component: MapBuilder,
})

function MapBuilder() {
  const state = useMapState()

  useEffect(() => {
    if (!localStorage.getItem('token')) state.navigate({ to: '/login' })
  }, [state.navigate])

  if (state.roomsLoading) return <div className="p-8 text-text">Loading map...</div>

  const getNPCsInRoom = useCallback((roomId: number) => state.npcs.filter(n => n.currentRoomId === roomId), [state.npcs])
  const getEquipmentInRoom = useCallback((_roomId: number) => state.roomEquipment, [state.roomEquipment])

  return (
    <div className="flex h-screen bg-surface">
      <Button
        variant="ghost" size="sm"
        onClick={() => state.setSidebarOpen(true)}
        aria-label="Open map sidebar"
        className="fixed top-3 left-3 z-50 p-2 bg-surface border border-border text-text-muted hover:bg-surface-muted hover:text-text lg:hidden"
      >
        <MenuIcon stroke="currentColor" />
      </Button>

      <div className={['lg:block lg:relative lg:inset-auto lg:z-auto', state.sidebarOpen ? 'block' : 'hidden'].join(' ')}>
        <div className="fixed inset-y-0 left-0 z-40 lg:static">
          <MapSidebar
            rooms={state.rooms}
            npcs={state.npcs}
            zLevels={state.zLevels}
            currentZLevel={state.currentZLevel}
            selectedRoom={state.selectedRoom}
            setCurrentZLevel={state.setCurrentZLevel}
            setSelectedRoom={state.setSelectedRoom}
            setShowCreateModal={state.setShowCreateModal}
          />
        </div>
      </div>

      {state.sidebarOpen && (
        <div className="fixed inset-0 bg-black/30 z-30 lg:hidden" onClick={() => state.setSidebarOpen(false)} />
      )}

      <div className="flex-1 overflow-hidden relative">
        <MapToolbar
          currentZLevel={state.currentZLevel}
          zoom={state.zoom}
          onZoom={state.handleZoom}
          onResetView={state.handleResetView}
          onRelayout={state.handleRelayout}
          onCleanupOrphanExits={state.cleanupOrphanExits}
        />

        <MapCanvas
          rooms={state.rooms}
          nodePositions={state.nodePositions}
          selectedRoom={state.selectedRoom}
          zoom={state.zoom}
          panOffset={state.panOffset}
          isDragging={state.isDragging}
          onWheel={state.handleWheel}
          onSelectRoom={state.setSelectedRoom}
          onDragStart={state.handleDragStart}
          onDragEnd={state.handleRoomDragEnd}
          getNPCsInRoom={getNPCsInRoom}
          getEquipmentInRoom={getEquipmentInRoom}
          viewportRef={state.viewportRef}
        />
      </div>

      {state.selectedRoom && !state.editingRoom && (
        <div className="w-[300px] bg-surface-muted border-l border-border flex flex-col lg:block hidden">
          <RoomDetailPanel
            selectedRoom={state.selectedRoom}
            rooms={state.rooms}
            zLevels={state.zLevels}
            onSelectRoom={state.setSelectedRoom}
            onEditRoom={state.handleEditRoom}
            onDeleteRoom={state.deleteRoom}
          />
        </div>
      )}

      {state.editingRoom && (
        <div className="w-[300px] bg-surface-muted border-l border-border flex flex-col lg:block hidden">
          <RoomEditor
            room={state.editingRoom}
            onCancel={() => state.setEditingRoom(null)}
          />
        </div>
      )}

      <CreateRoomModal
        isOpen={state.showCreateModal}
        onClose={() => state.setShowCreateModal(false)}
        onCreate={(input) => state.createRoom({ ...input, isStartingRoom: false, exits: {}, posX: 0, posY: 0 })}
        creating={state.isCreating}
      />

      {state.toast && (
        <div className="fixed bottom-4 right-4 bg-danger text-white px-4 py-2 rounded shadow-lg z-50">
          {state.toast}
        </div>
      )}
    </div>
  )
}