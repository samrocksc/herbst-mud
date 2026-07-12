import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "../../test/mocks/server";
import { RoomDetailPanel } from "./RoomDetailPanel";
import { ALL_DIRECTIONS } from "./DirectionUtils";
import type { Room } from "./types";

const baseRoom: Room = {
  id: 1,
  name: "Test Hub",
  description: "A test room",
  isStartingRoom: false,
  isRootRoom: true,
  exits: {},
  posZ: 0,
  atmosphere: "air",
  version: 1,
};

const noopRooms: Room[] = [baseRoom];
const noopZLevels = new Map<number, number>([[1, 0]]);

const noopSelect = vi.fn();
const noopEdit = vi.fn();
const noopDelete = vi.fn();
const noopAdd = vi.fn();
const noopConfirmAdd = vi.fn();
const noopCancelAdd = vi.fn();
const noopConfirmDelete = vi.fn();
const noopCancelDelete = vi.fn();

const renderPanel = (overrides: Partial<React.ComponentProps<typeof RoomDetailPanel>> = {}) =>
  render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <RoomDetailPanel
        selectedRoom={baseRoom}
        rooms={noopRooms}
        zLevels={noopZLevels}
        onSelectRoom={noopSelect}
        onEditRoom={noopEdit}
        onDeleteRoom={noopDelete}
        onRequestAddRoom={noopAdd}
        addRoomModal={undefined}
        onConfirmAddRoom={noopConfirmAdd}
        onCancelAddRoom={noopCancelAdd}
        deleteRoomModalOpen={false}
        onConfirmDeleteRoom={noopConfirmDelete}
        onCancelDeleteRoom={noopCancelDelete}
        {...overrides}
      />
    </QueryClientProvider>,
  );

describe("RoomDetailPanel — per-direction + Add Room button", () => {
  beforeEach(() => {
    noopAdd.mockClear();
    server.use(
      http.get("/api/rooms/:id/equipment", () => HttpResponse.json([])),
      http.get("/api/npc-instances", () => HttpResponse.json([])),
      http.get("/api/item-instances", () => HttpResponse.json([])),
    );
  });

  it("renders one + button per ALL_DIRECTIONS entry when the room has no exits", async () => {
    renderPanel();
    const addButtons = await screen.findAllByRole("button", { name: /^Add room to the / });
    expect(addButtons).toHaveLength(ALL_DIRECTIONS.length);
    for (const dir of ALL_DIRECTIONS) {
      expect(
        screen.getByRole("button", { name: `Add room to the ${dir}` }),
      ).toBeInTheDocument();
    }
  });

  it("clicking the + for a direction invokes onRequestAddRoom with that direction", async () => {
    renderPanel();
    const northButton = await screen.findByRole("button", { name: "Add room to the north" });
    fireEvent.click(northButton);
    expect(noopAdd).toHaveBeenCalledTimes(1);
    expect(noopAdd).toHaveBeenCalledWith("north");
  });

  it("clicking + for up direction invokes onRequestAddRoom with 'up'", async () => {
    renderPanel();
    const upButton = await screen.findByRole("button", { name: "Add room to the up" });
    fireEvent.click(upButton);
    expect(noopAdd).toHaveBeenCalledWith("up");
  });

  it("does not render a + button for a direction that already has an exit", async () => {
    const roomWithExit: Room = {
      ...baseRoom,
      exits: { north: 42 },
    };
    const roomMap = new Map<number, number>([[42, 0]]);
    renderPanel({ selectedRoom: roomWithExit, rooms: [...noopRooms, { ...baseRoom, id: 42, name: "Northern Room" }], zLevels: roomMap });

    expect(screen.queryByRole("button", { name: "Add room to the north" })).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Add room to the south" })).toBeInTheDocument();
  });

  it("renders the NewRoomModal when addRoomModal.open is true", async () => {
    renderPanel({
      addRoomModal: { open: true, fromRoom: baseRoom, dir: "north" },
    });
    // The modal title is rendered by the NewRoomModal inside its Modal wrapper.
    expect(await screen.findByText("Add Room to the north")).toBeInTheDocument();
  });

  it("does not render the NewRoomModal when addRoomModal is undefined", () => {
    renderPanel({ addRoomModal: undefined });
    expect(screen.queryByText(/^Add Room to the /)).not.toBeInTheDocument();
  });
});
