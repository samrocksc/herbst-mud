import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect, useRef, useCallback } from 'react'
import { useLogs, useLogServices } from '../../hooks/useLogs'
import { PageHeader } from '../../components/PageHeader'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/logs')({
  component: LogsPage,
})

const LEVEL_BADGE: Record<string, { bg: string; text: string; dot: string }> = {
  DEBUG: { bg: 'bg-slate-700/50', text: 'text-slate-300', dot: 'bg-slate-400' },
  INFO:  { bg: 'bg-sky-900/40', text: 'text-sky-300', dot: 'bg-sky-400' },
  WARN:  { bg: 'bg-amber-900/40', text: 'text-amber-300', dot: 'bg-amber-400' },
  ERROR: { bg: 'bg-red-900/40', text: 'text-red-300', dot: 'bg-red-400' },
}

const LEVELS = ['ALL', 'DEBUG', 'INFO', 'WARN', 'ERROR'] as const

type LogLine = {
  id?: number
  level: string
  message: string
  service?: string
  character_id?: number
  room_id?: number
  template_id?: string
  created_at: string
}

function LogsPage() {
  const [level, setLevel] = useState<string>('')
  const [service, setService] = useState<string>('')
  const [search, setSearch] = useState<string>('')
  const [live, setLive] = useState(false)
  const [liveLines, setLiveLines] = useState<LogLine[]>([])
  const bottomRef = useRef<HTMLDivElement>(null)

  const filters = { level: level || undefined, service: service || undefined, limit: 200 }
  const { data, isLoading } = useLogs(live ? undefined : filters)
  const { data: services } = useLogServices()

  const logs = live ? liveLines : (data?.logs ?? [])
  const filtered = search
    ? logs.filter((l) => l.message?.toLowerCase().includes(search.toLowerCase()))
    : logs

  useEffect(() => {
    if (live && bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [filtered.length, live])

  const toggleLive = useCallback(() => {
    setLive((prev) => !prev)
    setLiveLines([])
  }, [])

  useEffect(() => {
    if (!live) return
    const token = localStorage.getItem('auth_token') || ''
    const es = new EventSource(`${window.location.origin}/api/logs/stream?token=${encodeURIComponent(token)}`)
    es.onmessage = (e) => {
      try {
        const entry = JSON.parse(e.data) as LogLine
        setLiveLines((prev) => [entry, ...prev].slice(0, 500))
      } catch { /* skip malformed */ }
    }
    es.onerror = () => es.close()
    return () => es.close()
  }, [live])

  const activeLevel = level || 'ALL'
  const total = data?.total ?? 0

  return (
    <div className="management-page">
      <PageHeader title="Logs" backTo="/dashboard" actions={
        <div className="flex items-center gap-2">
          <span className="text-xs text-text-muted">{live ? `${liveLines.length} live` : `${total} total`}</span>
          <Button
            variant={live ? 'primary' : 'secondary'}
            size="sm"
            onClick={toggleLive}
          >
            {live ? '● Live' : '○ Paused'}
          </Button>
        </div>
      } />

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-3 mb-4 p-3 bg-surface-muted rounded-lg border border-border">
        <div className="flex gap-1">
          {LEVELS.map((l) => {
            const isActive = (l === 'ALL' && !level) || l === level
            const style = LEVEL_BADGE[l]
            return (
              <button
                key={l}
                onClick={() => setLevel(l === 'ALL' ? '' : l)}
                className={`px-2.5 py-1 rounded text-xs font-semibold transition-all ${
                  isActive
                    ? `${style ? style.bg : 'bg-primary/30'} ${style ? style.text : 'text-white'} ring-1 ring-current`
                    : 'bg-surface text-text-muted hover:text-text hover:bg-surface/80'
                }`}
              >
                {l}
              </button>
            )
          })}
        </div>

        <select
          value={service}
          onChange={(e) => setService(e.target.value)}
          className="bg-surface border border-border rounded-md px-2.5 py-1.5 text-sm text-text min-w-[140px]"
        >
          <option value="">All Services</option>
          {(services ?? []).map((s) => (
            <option key={s} value={s}>{s}</option>
          ))}
        </select>

        <div className="flex-1 min-w-[200px]">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Filter messages..."
            className="w-full bg-surface border border-border rounded-md px-3 py-1.5 text-sm text-text placeholder:text-text-muted/50"
          />
        </div>
      </div>

      {/* Log entries */}
      {isLoading && !live && (
        <div className="text-center py-12 text-text-muted">Loading logs...</div>
      )}

      {!isLoading && filtered.length === 0 && (
        <div className="text-center py-12">
          <div className="text-text-muted text-lg mb-1">No logs found</div>
          <div className="text-text-muted/60 text-sm">
            {live ? 'Waiting for new log entries...' : 'Try adjusting your filters.'}
          </div>
        </div>
      )}

      <div className="space-y-0">
        {filtered.slice(0, 300).map((log, i) => {
          const style = LEVEL_BADGE[log.level] ?? LEVEL_BADGE.DEBUG
          return (
            <div
              key={log.id ?? i}
              className={`flex items-start gap-3 px-3 py-2 border-b border-border/20 hover:bg-surface-muted/30 transition-colors`}
            >
              {/* Level dot + badge */}
              <div className="flex-shrink-0 w-16">
                <span className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-xs font-semibold ${style.bg} ${style.text}`}>
                  <span className={`w-1.5 h-1.5 rounded-full ${style.dot}`} />
                  {log.level}
                </span>
              </div>

              {/* Timestamp */}
              <div className="flex-shrink-0 w-28 text-xs text-text-muted/70 pt-0.5">
                {formatTime(log.created_at)}
              </div>

              {/* Service */}
              {log.service && (
                <div className="flex-shrink-0 w-28 text-xs text-text-muted/80 pt-0.5 truncate" title={log.service}>
                  {log.service}
                </div>
              )}

              {/* Message */}
              <div className="flex-1 text-sm text-text break-all min-w-0">
                {log.message}
              </div>

              {/* Context badges */}
              <div className="flex-shrink-0 flex items-center gap-2 text-xs text-text-muted/60 pt-0.5">
                {log.character_id != null && (
                  <span className="bg-surface px-1.5 py-0.5 rounded" title="Character ID">
                    c:{log.character_id}
                  </span>
                )}
                {log.room_id != null && (
                  <span className="bg-surface px-1.5 py-0.5 rounded" title="Room ID">
                    r:{log.room_id}
                  </span>
                )}
              </div>
            </div>
          )
        })}
      </div>

      {filtered.length > 300 && (
        <div className="text-center py-3 text-text-muted text-sm">
          Showing 300 of {filtered.length} entries. Narrow your filters to see more.
        </div>
      )}

      <div ref={bottomRef} />
    </div>
  )
}

function formatTime(t: string): string {
  if (!t) return '—'
  try {
    const d = new Date(t)
    const now = new Date()
    const diffMs = now.getTime() - d.getTime()
    const diffMin = Math.floor(diffMs / 60000)
    if (diffMin < 1) return 'just now'
    if (diffMin < 60) return `${diffMin}m ago`
    const isToday = d.toDateString() === now.toDateString()
    if (isToday) return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    const isThisYear = d.getFullYear() === now.getFullYear()
    if (isThisYear) return d.toLocaleDateString([], { month: 'short', day: 'numeric' }) + ' ' + d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    return d.toLocaleDateString([], { year: 'numeric', month: 'short', day: 'numeric' }) + ' ' + d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  } catch {
    return t
  }
}