import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useUsers, useResetPassword, useDeleteCharacter, type User } from "../../hooks/useUsers";
import { apiGet } from "../../utils/apiFetch";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";

export const Route = createFileRoute("/_auth/players")({
  component: PlayersManagement,
});

type Character = Readonly<{
  id: number
  name: string
  isNPC: boolean
  race: string
  class: string
  level: number
  hitpoints: number
  max_hitpoints: number
  currentRoomId: number
}>

type UserCharacter = Readonly<{
  id: number
  name: string
  hitpoints: number
  max_hitpoints: number
}>

function PlayersManagement() {
  const { data: users, isLoading, error } = useUsers();
  const resetPassword = useResetPassword();
  const deleteCharacter = useDeleteCharacter();
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [showDetail, setShowDetail] = useState(false);
  const [showCharacters, setShowCharacters] = useState(false);
  const [resetSuccess, setResetSuccess] = useState<string | null>(null);
  const [resetError, setResetError] = useState<string | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const charactersQuery = useQuery<Character[]>({
    queryKey: ["characters"],
    queryFn: () => apiGet<Character[]>(`${window.location.origin}/characters`),
    enabled: showCharacters,
  });

  const userCharactersQuery = useQuery<UserCharacter[]>({
    queryKey: ["user-characters", selectedUser?.id],
    queryFn: () => {
      if (!selectedUser) return Promise.resolve([]);
      return apiGet<UserCharacter[]>(`${window.location.origin}/user-characters/${selectedUser.id}`);
    },
    enabled: !!selectedUser && showDetail,
  });

  const handleReset = async (user: User) => {
    setResetSuccess(null); setResetError(null);
    try { await resetPassword.mutateAsync(user.id); setResetSuccess(`Password reset for ${user.email}`); }
    catch { setResetError(`Failed to reset password for ${user.email}`); }
  };

  const handleDelete = async (id: number) => {
    setDeleteError(null);
    try { await deleteCharacter.mutateAsync(id); }
    catch { setDeleteError("Failed to delete character"); }
  };

  const formatDate = (dateStr: string) =>
    new Date(dateStr).toLocaleDateString("en-US", { year: "numeric", month: "short", day: "numeric", hour: "2-digit", minute: "2-digit" });

  const userColumns: Column<User>[] = [
    { header: "ID", accessor: "id" },
    { header: "Email", accessor: "email" },
    { header: "Role", accessor: "is_admin", render: (val: unknown) =>
      val ? <span className="badge badge-admin">Admin</span> : <span className="badge badge-player">Player</span> },
    { header: "Created", accessor: "created_at", render: (val: unknown) => formatDate(String(val ?? "")) },
    { header: "Actions", accessor: "_actions", render: (_: unknown, row: User) => (
      <Button variant="secondary" size="sm" onClick={(e) => { e.stopPropagation(); handleReset(row); }} disabled={resetPassword.isPending}>Reset Password</Button>
    )},
  ];

  const charColumns: Column<Character>[] = [
    { header: "ID", accessor: "id" },
    { header: "Name", accessor: "name", render: (_: unknown, row: Character) => (
      <Link to="/characters/$characterId" params={{ characterId: String(row.id) }} className="text-primary no-underline hover:underline font-bold">{row.name}</Link>
    )},
    { header: "Race", accessor: "race" },
    { header: "Class", accessor: "class" },
    { header: "Level", accessor: "level" },
    { header: "HP", accessor: "hitpoints", render: (_: unknown, row: Character) => `${row.hitpoints}/${row.max_hitpoints}` },
    { header: "Room", accessor: "currentRoomId" },
    { header: "Actions", accessor: "_actions", render: (_: unknown, row: Character) => (
      <Button variant="danger" size="sm" onClick={() => handleDelete(row.id)} disabled={deleteCharacter.isPending}>Delete</Button>
    )},
  ];

  if (isLoading) return <div className="loading">Loading players...</div>;
  if (error) return <div className="error">Failed to load players: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader title="Players Management" backTo="/dashboard" />
      {resetSuccess && <div className="success-message">{resetSuccess}</div>}
      {resetError && <div className="error-message">{resetError}</div>}
      {deleteError && <div className="error-message">{deleteError}</div>}
      <DataTable columns={userColumns} data={users ?? []} getKey={(row: User) => row.id}
        onRowClick={(row: User) => { setSelectedUser(row); setShowDetail(true); }} emptyMessage="No players found." />

      <div className="mt-6">
        <div className="flex items-center justify-between mb-3">
          <h2 className="m-0 text-text text-lg font-semibold">Characters</h2>
          <Button variant="secondary" size="sm" onClick={() => setShowCharacters(!showCharacters)}>
            {showCharacters ? "Hide Characters" : "Show Characters"}
          </Button>
        </div>
        {showCharacters && (
          charactersQuery.isLoading ? <div className="text-text-muted text-sm">Loading characters...</div> :
          charactersQuery.isError ? <div className="text-danger text-sm">Failed to load characters</div> :
          <DataTable columns={charColumns} data={charactersQuery.data ?? []} getKey={(row: Character) => row.id} emptyMessage="No characters found." />
        )}
      </div>

      {showDetail && selectedUser && (
        <div className="modal-overlay" onClick={() => setShowDetail(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Player Details</h3>
              <Button variant="ghost" size="sm" onClick={() => setShowDetail(false)}>×</Button>
            </div>
            <div className="modal-body">
              <div className="detail-row"><label>ID:</label><span>{selectedUser.id}</span></div>
              <div className="detail-row"><label>Email:</label><span>{selectedUser.email}</span></div>
              <div className="detail-row"><label>Admin:</label><span>{selectedUser.is_admin ? "Yes" : "No"}</span></div>
              <div className="detail-row"><label>Created:</label><span>{formatDate(selectedUser.created_at)}</span></div>

              <div className="mt-4">
                <h4 className="text-text text-base font-semibold mb-2">Characters</h4>
                {userCharactersQuery.isLoading ? <div className="text-text-muted text-sm">Loading...</div> :
                 userCharactersQuery.isError ? <div className="text-danger text-sm">Failed to load characters</div> :
                 (userCharactersQuery.data?.length ?? 0) === 0 ? <div className="text-text-muted text-sm">No characters</div> :
                 <div className="space-y-2">
                   {userCharactersQuery.data!.map((char) => (
                     <div key={char.id} className="flex items-center justify-between bg-surface-muted rounded px-3 py-2">
                       <div>
                         <span className="font-semibold text-text">{char.name}</span>
                         <span className="text-text-muted text-sm ml-2">HP {char.hitpoints}/{char.max_hitpoints}</span>
                       </div>
                       <Button variant="danger" size="sm" onClick={() => handleDelete(char.id)} disabled={deleteCharacter.isPending}>Delete</Button>
                     </div>
                   ))}
                 </div>
                }
              </div>
            </div>
            <div className="modal-footer">
              <Button variant="secondary" size="sm" onClick={() => handleReset(selectedUser)}>Reset Password</Button>
              <Button variant="secondary" onClick={() => setShowDetail(false)}>Close</Button>
            </div>
          </div>
        </div>
      )}
    </PageContainer>
  );
}
