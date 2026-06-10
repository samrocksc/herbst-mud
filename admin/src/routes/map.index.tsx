/* eslint-disable react-hooks/exhaustive-deps */

import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useCallback, useMemo } from "react";
import { useMapState } from "../hooks/useMapState";
import { MapSidebar } from "../components/map/MapSidebar";
import { MapToolbar } from "../components/map/MapToolbar";
import { MapCanvas } from "../components/map/MapCanvas";
import { RoomDetailPanel } from "../components/map/RoomDetailPanel";
import { RoomEditor } from "../components/map/RoomEditor";
import { DeleteConfirmation } from "../components/DeleteConfirmation";
import { useNavigate } from "@tanstack/react-router";

export const Route = createFileRoute("/map/")({
  component: MapBuilder,
});

function MapBuilder() {
  const state = useMapState();

  useEffect(() => {
    if (!localStorage.getItem("token")) state.navigate({ to: "/login" });
  }, [state.navigate]);

  const getNPCsInRoom = useCallback((roomId: number) => state.npcs.filter(n => n.currentRoomId === roomId), [state.npcs]);
  const getEquipmentInRoom = useCallback((_roomId: number) => state.roomEquipment, [state.roomEquipment]);

  if (state.roomsLoading) return <div className="p-8 text-text">Loading map...</div>;

  return (
    <div className="flex h-[100dvh] bg-surface">
      {state.sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/30 z-40 lg:hidden"
          onClick={() => state.setSidebarOpen(false)}
        />
      )}
      <div className={`shrink-0 lg:relative ${state.sidebarOpen ? 'fixed inset-y-0 left-0 z-40 w-[220px]' : 'hidden'} lg:block lg:inset-auto lg:z-auto lg:w-auto`}>
        <MapSidebar
          rooms={state.rooms}
          npcs={state.npcs}
          zLevels={state.zLevels}
          currentZLevel={state.currentZLevel}
          selectedRoom={state.selectedRoom}
          setCurrentZLevel={state.setCurrentZLevel}
          setSelectedRoom={state.setSelectedRoom}
        />
      </div>

      <div className="flex-1 overflow-hidden relative">
        <MapToolbar
          currentZLevel={state.currentZLevel}
          zLevels={Array.from(new Set(Array.from(state.zLevels.values()))).sort((a, b) => a - b)}
          zoom={state.zoom}
          onZoom={state.handleZoom}
          onResetView={state.handleResetView}
          onRelayout={state.handleRelayout}
          onCleanupOrphanExits={state.handleRequestCleanupOrphanExits}
          onGoToFloor={state.handleSetZLevel}
          onAddFloor={state.handleAddFloor}
          onDeleteFloor={state.requestDeleteFloor}
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
          currentZLevel={state.currentZLevel}
          onCreateRoom={() => state.navigate({ to: "/map/rooms/new", search: { floor: state.currentZLevel } })}
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
            onRequestDeleteRoom={state.requestDeleteRoom}
            onAddRoom={state.requestAddRoom}
            addRoomModal={state.addRoomModal}
            onConfirmAddRoom={state.confirmAddRoom}
            onCancelAddRoom={state.cancelAddRoom}
            isAddingRoom={state.isAddingRoom}
            deleteRoomModalOpen={state.deleteRoomModalOpen}
            onConfirmDeleteRoom={state.confirmDeleteRoom}
            onCancelDeleteRoom={state.cancelDeleteRoom}
            isDeletingRoom={state.isDeletingRoom}
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

      {state.toast && (
        <div className="fixed bottom-4 right-4 bg-danger text-white px-4 py-2 rounded shadow-lg z-50">
          {state.toast}
        </div>
      )}
      <DeleteConfirmation
        open={state.cleanupConfirmOpen}
        title="Cleanup Orphan Exits"
        message="This will scan all rooms for exits pointing to deleted rooms and remove them. This action cannot be undone. Continue?"
        onConfirm={state.handleConfirmCleanupOrphanExits}
        onCancel={state.handleCancelCleanupOrphanExits}
      />
      <DeleteConfirmation
        open={state.deleteFloorModalOpen}
        title="Delete Floor"
        message={`This will delete all rooms on floor ${state.currentZLevel}. This cannot be undone.`}
        onConfirm={state.confirmDeleteFloor}
        onCancel={state.cancelDeleteFloor}
      />
    </div>
  );
}
