import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useUsers, useResetPassword, type User } from '../../hooks/useUsers'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/_auth/players')({
  component: PlayersManagement,
})

function PlayersManagement() {
  const { data: users, isLoading, error } = useUsers()
  const resetPassword = useResetPassword()
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [showDetail, setShowDetail] = useState(false)
  const [resetSuccess, setResetSuccess] = useState<string | null>(null)
  const [resetError, setResetError] = useState<string | null>(null)

  const handleResetPassword = async (user: User) => {
    setResetSuccess(null)
    setResetError(null)
    try {
      await resetPassword.mutateAsync(user.id)
      setResetSuccess(`Password reset successfully for ${user.email}`)
    } catch {
      setResetError(`Failed to reset password for ${user.email}`)
    }
  }

  const handleRowClick = (user: User) => {
    setSelectedUser(user)
    setShowDetail(true)
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (isLoading) return <div className="loading">Loading players...</div>
  if (error) return <div className="error">Failed to load players: {error.message}</div>

  return (
    <div className="management-page">
      <PageHeader title="Players Management" backTo="/dashboard" />

      {resetSuccess && <div className="success-message">{resetSuccess}</div>}
      {resetError && <div className="error-message">{resetError}</div>}

      <div className="players-table-container">
        <table className="players-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Email</th>
              <th>Admin</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users?.map((user) => (
              <tr key={user.id} onClick={() => handleRowClick(user)} className="clickable-row">
                <td>{user.id}</td>
                <td>{user.email}</td>
                <td>
                  {user.is_admin ? (
                    <span className="badge badge-admin">Admin</span>
                  ) : (
                    <span className="badge badge-player">Player</span>
                  )}
                </td>
                <td>{formatDate(user.created_at)}</td>
                <td onClick={(e) => e.stopPropagation()}>
                  <button
                    className="btn-reset"
                    onClick={() => handleResetPassword(user)}
                    disabled={resetPassword.isPending}
                  >
                    Reset Password
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {showDetail && selectedUser && (
        <div className="modal-overlay" onClick={() => setShowDetail(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Player Details</h3>
              <button className="modal-close" onClick={() => setShowDetail(false)}>×</button>
            </div>
            <div className="modal-body">
              <div className="detail-row">
                <label>ID:</label>
                <span>{selectedUser.id}</span>
              </div>
              <div className="detail-row">
                <label>Email:</label>
                <span>{selectedUser.email}</span>
              </div>
              <div className="detail-row">
                <label>Admin:</label>
                <span>{selectedUser.is_admin ? 'Yes' : 'No'}</span>
              </div>
              <div className="detail-row">
                <label>Character:</label>
                <span>{selectedUser.character_name || 'No character'}</span>
              </div>
              <div className="detail-row">
                <label>Created:</label>
                <span>{formatDate(selectedUser.created_at)}</span>
              </div>
            </div>
            <div className="modal-footer">
              <button className="btn-reset" onClick={() => handleResetPassword(selectedUser)}>
                Reset Password
              </button>
              <button className="btn-cancel" onClick={() => setShowDetail(false)}>
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
