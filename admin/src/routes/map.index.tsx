/* eslint-disable react-hooks/exhaustive-deps */

import { createFileRoute } from "@tanstack/react-router";
import { useEffect } from "react";
import { useMapState } from "../hooks/useMapState";
import { MapSidebar } from "../components/map/MapSidebar";
import { MapToolbar } from "../components/map/MapToolbar";
import { ReactFlowCanvas } from "../components/map/ReactFlowCanvas";
import { RoomDetailPanel } from "../components/map/RoomDetailPanel";
import { RoomEditor } from "../components/map/RoomEditor";
import { DeleteConfirmation } from "../components/DeleteConfirmation";

export const Route = createFileRoute("/map/")({
  component: MapBuilder,
});

function MapBuilder() {
  const state = useMapState();

  useEffect(() => {
    if (!localStorage.getItem("token")) state.navigate({ to: "/login" });
  }, [state.navigate]);

  if (state.roomsLoading) return <div className="p-8 text-text">Loading map...</div>;

  return (
    <div className="flex h-[100dvh] bg-surface">
      {state.sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/30 z-40 lg:hidden"
          onClick={() => state.setSidebarOpen(false)}
        />
      )}
      <div className={`shrink-0 lg:relative ${state.sidebarOpen ? "fixed inset-y-0 left-0 z-40 w-[220px]" : "hidden"} lg:block lg:inset-auto lg:z-auto lg:w-auto`}>
        <MapSidebar
          rooms={state.rooms}
          npcs={state.npcs}
          zLevels={state.zLevels}
          currentZLevel={state.currentZLevel}
          selectedRoom={state.selectedRoom}
          setCurrentZLevel={state.handleSetZLevel}
          setSelectedRoom={state.handleSelectRoom}
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
          onRealign={state.handleRealign}
          onCleanupOrphanExits={state.handleRequestCleanupOrphanExits}
          onGoToFloor={state.handleSetZLevel}
          onAddFloor={state.handleAddFloor}
          onDeleteFloor={state.requestDeleteFloor}
        />

        {/* React Flow Canvas - main map rendering */}
        <ReactFlowCanvas
          nodes={state.reactFlowNodes}
          edges={state.reactFlowEdges}
        />
      </div>

      {state.selectedRoom && !state.editingRoom && (
        <div className="w-[300px] h-full bg-surface-muted border-l border-border flex-col hidden lg:flex">
          <RoomDetailPanel
            selectedRoom={state.selectedRoom}
            rooms={state.rooms}
            zLevels={state.zLevels}
            onSelectRoom={state.handleSelectRoom}
            onEditRoom={state.handleEditRoom}
            onRequestDeleteRoom={state.requestDeleteRoom}
            deleteRoomModalOpen={state.deleteRoomModalOpen}
            onConfirmDeleteRoom={state.confirmDeleteRoom}
            onCancelDeleteRoom={state.cancelDeleteRoom}
            isDeletingRoom={state.isDeletingRoom}
            onRequestAddRoom={state.requestAddRoom}
            addRoomModal={state.addRoomModal}
            onConfirmAddRoom={state.confirmAddRoom}
            onCancelAddRoom={state.cancelAddRoom}
            isAddingRoom={state.isAddingRoom}
          />
        </div>
      )}

      {state.editingRoom && (
        <div className="w-[300px] h-full bg-surface-muted border-l border-border flex-col hidden lg:flex">
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
