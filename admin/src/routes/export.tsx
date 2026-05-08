import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState, useRef } from 'react'
import { Button } from '../components/Button'

export const Route = createFileRoute('/export')({
  component: ExportPage,
})

type ExportData = Readonly<{
  version: string
  exported_at: string
  rooms: readonly unknown[]
  npcs: readonly unknown[]
  skills: readonly unknown[]
  items: readonly unknown[]
}>

function ExportPage() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState('')
  const [messageType, setMessageType] = useState<'success' | 'error'>('success')
  const [exportPreview, setExportPreview] = useState<ExportData | null>(null)
  const [showImportConfirm, setShowImportConfirm] = useState(false)
  const [showWipeConfirm, setShowWipeConfirm] = useState(false)
  const [importData, setImportData] = useState<ExportData | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleExport = async () => {
    setLoading(true)
    setMessage('')
    
    try {
      const response = await fetch(`\${window.location.origin}/admin/export`)
      if (!response.ok) {
        throw new Error('Export failed: ' + response.statusText)
      }
      
      const data: ExportData = await response.json()
      setExportPreview(data)
      
      // Create and download JSON file
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `herbst-mud-export-${new Date().toISOString().split('T')[0]}.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
      
      setMessage(`Exported ${data.rooms.length} rooms, ${data.npcs.length} NPCs, ${data.skills.length} skills`)
      setMessageType('success')
    } catch (err) {
      setMessage('Export failed: ' + (err instanceof Error ? err.message : 'Unknown error'))
      setMessageType('error')
    } finally {
      setLoading(false)
    }
  }

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    setLoading(true)
    setMessage('')

    try {
      const text = await file.text()
      const data: ExportData = JSON.parse(text)
      
      // Validate the file
      const response = await fetch(`${window.location.origin}/admin/import/validate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: text
      })
      
      const validation = await response.json()
      
      if (!validation.is_valid) {
        setMessage('Validation failed: ' + validation.errors.join(', '))
        setMessageType('error')
        return
      }
      
      setImportData(data)
      setShowImportConfirm(true)
      setMessage(`File validated: ${validation.rooms} rooms, ${validation.npcs} NPCs, ${validation.skills} skills`)
      setMessageType('success')
    } catch (err) {
      setMessage('Failed to parse file: ' + (err instanceof Error ? err.message : 'Invalid JSON'))
      setMessageType('error')
    } finally {
      setLoading(false)
    }
  }

  const handleImport = async () => {
    if (!importData) return

    setLoading(true)
    setMessage('')

    try {
      const response = await fetch(`${window.location.origin}/admin/import`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(importData)
      })
      
      if (!response.ok) {
        throw new Error('Import failed: ' + response.statusText)
      }
      
      const result = await response.json()
      setMessage(`Import successful! imported ${result.imported.rooms} rooms, ${result.imported.npcs} NPCs`)
      setMessageType('success')
      setShowImportConfirm(false)
      setImportData(null)
    } catch (err) {
      setMessage('Import failed: ' + (err instanceof Error ? err.message : 'Unknown error'))
      setMessageType('error')
    } finally {
      setLoading(false)
    }
  }

  const handleWipe = async () => {
    setLoading(true)
    setMessage('')

    try {
      const response = await fetch(`${window.location.origin}/admin/wipe/full`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({}) // Full wipe uses POST body
      })
      
      if (!response.ok) {
        throw new Error('Wipe failed: ' + response.statusText)
      }
      
      const result = await response.json()
      setMessage(`Wiped ${result.npcs_wiped} NPCs, ${result.rooms_wiped} rooms. Reinitialized: ${result.reinitialized.join(', ')}`)
      setMessageType('success')
      setShowWipeConfirm(false)
    } catch (err) {
      setMessage('Wipe failed: ' + (err instanceof Error ? err.message : 'Unknown error'))
      setMessageType('error')
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
    localStorage.removeItem('email')
    localStorage.removeItem('isAdmin')
    navigate({ to: '/login' })
  }

  return (
    <div className="min-h-screen bg-surface text-text p-8">
      <div className="max-w-[1200px] mx-auto">
        <div className="flex justify-between items-center mb-8 border-b border-border pb-4">
          <div className="flex items-center gap-4">
            <Button
              onClick={() => navigate({ to: '/dashboard' })}
              variant="secondary"
              size="sm"
            >
              ← Back
            </Button>
            <h1 className="text-primary">Game Export / Import</h1>
          </div>
          <Button onClick={handleLogout} variant="danger">
            Logout
          </Button>
        </div>

        {message && (
          <div className={`rounded-lg p-4 mb-6 ${messageType === 'success' ? 'bg-green-900/30 border border-green-600' : 'bg-red-900/30 border border-red-600'}`}>
            {message}
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Export Section */}
          <div className="bg-surface-muted rounded-lg p-6 border border-border">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <span>💾</span>
              Export Game World
            </h2>
            
            <p className="text-text-muted mb-4">
              Export all game data including rooms, NPCs, skills, and items. 
              Player accounts and player characters are excluded.
            </p>

            <div className="space-y-2 mb-6 text-sm text-text-muted">
              <div>• Rooms, exits, and descriptions</div>
              <div>• NPC stats and locations (isNPC=true only)</div>
              <div>• Skills and abilities</div>
              <div>• Items in rooms and on NPCs</div>
            </div>

            <Button
              onClick={handleExport}
              disabled={loading}
              variant="primary"
              size="sm"
            >
              {loading ? 'Exporting...' : '📥 Export to JSON'}
            </Button>

            {exportPreview && (
              <div className="mt-4 p-4 bg-surface rounded border border-border">
                <h3 className="font-bold mb-2">Last Export Preview:</h3>
                <div className="text-sm text-text-muted space-y-1">
                  <div>Version: {exportPreview.version}</div>
                  <div>Exported: {new Date(exportPreview.exported_at).toLocaleString()}</div>
                  <div>Rooms: {exportPreview.rooms.length}</div>
                  <div>NPCs: {exportPreview.npcs.length}</div>
                  <div>Skills: {exportPreview.skills.length}</div>
                  <div>Items: {exportPreview.items.length}</div>
                </div>
              </div>
            )}
          </div>

          {/* Import Section */}
          <div className="bg-surface-muted rounded-lg p-6 border border-border">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
              <span>📂</span>
              Import Game World
            </h2>
            
            <p className="text-text-muted mb-4">
              Import a previously exported game world. Existing data will be updated, 
              new data will be created. Player data is never imported.
            </p>

            <div className="space-y-2 mb-6 text-sm text-text-muted">
              <div>• Validates JSON structure before import</div>
              <div>• Updates existing rooms/NPCs by name</div>
              <div>• Creates new entries for unknown data</div>
              <div>• Preserves player accounts and characters</div>
            </div>

            <input
              type="file"
              accept=".json"
              ref={fileInputRef}
              onChange={handleFileSelect}
              className="hidden"
            />

            {!showImportConfirm ? (
              <Button
                onClick={() => fileInputRef.current?.click()}
                disabled={loading}
                variant="primary"
                size="sm"
              >
                {loading ? 'Validating...' : '📤 Select JSON File'}
              </Button>
            ) : (
              <div className="space-y-3">
                <div className="p-4 bg-yellow-900/30 border border-yellow-600 rounded">
                  <p className="font-bold text-yellow-400">⚠️ Confirm Import</p>
                  <p className="text-sm text-text-muted mt-1">
                    This will modify your game world data. Make sure you have a backup!
                  </p>
                </div>
                <div className="flex gap-3">
                  <Button
                    onClick={handleImport}
                    disabled={loading}
                    variant="danger"
                    fullWidth
                  >
                    {loading ? 'Importing...' : '✓ Confirm Import'}
                  </Button>
                  <Button
                    onClick={() => {
                      setShowImportConfirm(false)
                      setImportData(null)
                      if (fileInputRef.current) fileInputRef.current.value = ''
                    }}
                    variant="secondary"
                    fullWidth
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            )}
          </div>

          {/* Wipe & Reload Section */}
          <div className="bg-red-900/20 rounded-lg p-6 border border-red-800 md:col-span-2">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2 text-red-400">
              <span>⚠️</span>
              Danger Zone: Wipe & Reload
            </h2>
            
            <p className="text-text-muted mb-4">
              <strong>WARNING:</strong> This will delete ALL game data (NPCs, rooms, items, skills, abilities)
              and reinitialize with fresh default data. Player accounts and characters are preserved.
            </p>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6 text-sm text-text-muted">
              <div className="flex items-center gap-2">
                <span className="text-red-500">✗</span> NPCs deleted
              </div>
              <div className="flex items-center gap-2">
                <span className="text-red-500">✗</span> Rooms cleared
              </div>
              <div className="flex items-center gap-2">
                <span className="text-red-500">✗</span> Items wiped
              </div>
              <div className="flex items-center gap-2">
                <span className="text-red-500">✗</span> Skills reset
              </div>
              <div className="flex items-center gap-2">
                <span className="text-green-500">✓</span> Users preserved
              </div>
              <div className="flex items-center gap-2">
                <span className="text-green-500">✓</span> Players kept
              </div>
              <div className="flex items-center gap-2">
                <span className="text-green-500">✓</span> Defaults reloaded
              </div>
            </div>

            {!showWipeConfirm ? (
              <Button
                onClick={() => setShowWipeConfirm(true)}
                disabled={loading}
                variant="danger"
                fullWidth
              >
                🗑️ Wipe & Reload Game World
              </Button>
            ) : (
              <div className="space-y-3">
                <div className="p-4 bg-red-900/50 border border-red-600 rounded">
                  <p className="font-bold text-red-400 text-lg">☠️ FINAL WARNING</p>
                  <p className="text-text-muted mt-2">
                    This action is <strong>IRREVERSIBLE</strong>. All game data will be lost.
                  </p>
                  <p className="text-sm text-text-muted mt-1">
                    Consider exporting your current game world before wiping!
                  </p>
                </div>
                <div className="flex gap-3">
                  <Button
                    onClick={handleWipe}
                    disabled={loading}
                    variant="danger"
                    fullWidth
                  >
                    {loading ? 'Wiping...' : '☠️ YES - WIPE EVERYTHING'}
                  </Button>
                  <Button
                    onClick={() => setShowWipeConfirm(false)}
                    variant="secondary"
                    fullWidth
                  >
                    Cancel - Keep My Data
                  </Button>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* JSON Format Documentation */}
        <div className="mt-8 bg-surface-muted rounded-lg p-6 border border-border">
          <h2 className="text-xl font-bold mb-4">Export Format Documentation</h2>
          <pre className="bg-surface p-4 rounded overflow-x-auto text-xs text-text-muted">
{`{
  "version": "1.0",
  "exported_at": "2026-04-05T14:30:00Z",
  "rooms": [
    {
      "id": 1,
      "name": "Town Square",
      "description": "The center of town...",
      "is_starting": true,
      "exits": [
        { "direction": "north", "target_room_id": 2 },
        { "direction": "south", "target_room_id": 3 }
      ]
    }
  ],
  "npcs": [
    {
      "id": 10,
      "name": "Aragorn",
      "current_room_id": 5,
      "race": "human",
      "class": "warrior",
      "level": 10,
      "hitpoints": 100,
      "max_hitpoints": 100,
      "stamina": 50,
      "max_stamina": 50,
      "mana": 30,
      "max_mana": 30,
      "strength": 15,
      "dexterity": 12,
      "constitution": 14,
      "intelligence": 10,
      "wisdom": 11,
      "npc_skill_id": "druid_heal",
      "is_immortal": false
    }
  ],
  "skills": [...],
  "items": []
}`}
          </pre>
        </div>
      </div>
    </div>
  )
}
