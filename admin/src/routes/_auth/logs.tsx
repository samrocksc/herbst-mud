import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect, useRef, useCallback } from 'react'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/_auth/logs')({
  component: LogsPage,
})

type LogEntry = Readonly<{
  id: number
  level: string
  message: string
  service?: string
  character_id?: number
  room_id?: number
  template_id?: string
  metadata?: Record<string, unknown>
  created_at: string
}>

type LogsResponse = Readonly<{
  logs: LogEntry[]
  total: number
  limit: number
  offset: number
}>

type ServicesResponse = Readonly<{
  services: string[]
}>

const LEVEL_COLORS: Record<string, string> = {
  DEBUG: 'var(--color-text-muted)',
  INFO: 'var(--color-info)',
  WARN: 'var(--color-warning)',
  ERROR: 'var(--color-error)',
}

const LEVEL_BADGE_COLORS: Record<string, string> = {
  DEBUG: 'var(--color-bg-subtle)',
  INFO: 'var(--color-info-bg)',
  WARN: 'var(--color-warning-bg)',
  ERROR: 'var(--color-error-bg)',
}

const ALL_LEVELS = ['DEBUG', 'INFO', 'WARN', 'ERROR']

function LogsPage() {
  const [live, setLive] = useState(false)
  const [levels, setLevels] = useState<Set<string>>(new Set(['WARN', 'ERROR', 'INFO']))
  const [service, setService] = useState('')
  const [search, setSearch] = useState('')
  const [services, setServices] = useState<string[]>([])
  const [entries, setEntries] = useState<LogEntry[]>([])
  const [paused, setPaused] = useState(false)
  const tableRef = useRef<HTMLDivElement>(null)
  const esRef = useRef<EventSource | null>(null)

  // Fetch initial log history
  const fetchHistory = useCallback(async () => {
    const levelParam = levels.size > 0 && levels.size < 4
      ? Array.from(levels).join(',')
      : ''
    const params = new URLSearchParams({ limit: '200' })
    if (levelParam) params.set('level', levelParam)
    if (service) params.set('service', service)

    try {
      const token = localStorage.getItem('auth_token')
      const resp = await fetch(`/api/logs?${params}`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      })
      if (!resp.ok) return
      const data: LogsResponse = await resp.json()
      setEntries(data.logs)
    } catch { /* ignore */ }
  }, [levels, service])

  // Fetch available services
  const fetchServices = useCallback(async () => {
    try {
      const token = localStorage.getItem('auth_token')
      const resp = await fetch('/api/logs/services', {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      })
      if (!resp.ok) return
      const data: ServicesResponse = await resp.json()
      setServices(data.services)
    } catch { /* ignore */ }
  }, [])

  // SSE connection for live tail
  useEffect(() => {
    if (!live) {
      esRef.current?.close()
      esRef.current = null
      return
    }

    fetchServices()

    const token = localStorage.getItem('auth_token')
    const url = new URL('/api/logs/stream', window.location.origin)
    const es = new EventSource(url.toString())
    esRef.current = es

    es.onmessage = (event) => {
      const entry: LogEntry = JSON.parse(event.data)
      setEntries((prev) => {
        // Apply client-side filters
        if (levels.size > 0 && !levels.has(entry.level)) return prev
        if (service && entry.service !== service) return prev
        return [...prev, entry].slice(-500) // keep last 500
      })
    }
    
    es.onerror = () => {
      es.close()
      esRef.current = null
      setLive(false)
    }

    return () => {
      es.close()
      esRef.current = null
    }
  }, [live, levels, service, fetchServices])

  // Load history on mount and when filters change (non-live)
  useEffect(() => {
    fetchHistory()
    fetchServices()
  }, [fetchHistory, fetchServices])

  // Auto-scroll when new entries arrive
  useEffect(() => {
    if (!paused && tableRef.current) {
      tableRef.current.scrollTop = tableRef.current.scrollHeight
    }
  }, [entries, paused])

  const toggleLevel = (level: string) => {
    setLevels((prev) => {
      const next = new Set(prev)
      if (next.has(level)) {
        next.delete(level)
      } else {
        next.add(level)
      }
      return next
    })
  }

  const filtered = entries.filter((e) => {
    if (levels.size > 0 && !levels.has(e.level)) return false
    if (service && e.service !== service) return false
    if (search && !e.message.toLowerCase().includes(search.toLowerCase())) return false
    return true
  })

  return (
    <div className="h-full flex flex-col">
      <PageHeader backTo="/dashboard" title="Logs" />

      {/* Filter bar */}
      <div className="flex flex-wrap items-center gap-3 px-4 py-2 border-b border-[var(--color-border)] bg-[var(--color-bg)] sticky top-0 z-10">
        {/* Live toggle */}
        <button
          onClick={() => setLive((v) => !v)}
          aria-label={live ? 'Disable live tail' : 'Enable live tail'}
          className="px-3 py-1 text-xs font-medium rounded-full border transition-colors"
          style={{
            background: live ? 'var(--color-accent)' : 'transparent',
            color: live ? 'var(--color-bg)' : 'var(--color-text)',
            borderColor: live ? 'var(--color-accent)' : 'var(--color-border)',
          }}
        >
          {live ? '● Live' : '○ Live'}
        </button>

        {/* Level chips */}
        {ALL_LEVELS.map((lvl) => (
          <button
            key={lvl}
            onClick={() => toggleLevel(lvl)}
            aria-label={`Toggle ${lvl} level`}
            className="px-2 py-0.5 text-xs font-medium rounded-full border transition-colors"
            style={{
              background: levels.has(lvl) ? LEVEL_BADGE_COLORS[lvl] : 'transparent',
              color: levels.has(lvl) ? LEVEL_COLORS[lvl] : 'var(--color-text-muted)',
              borderColor: levels.has(lvl) ? LEVEL_COLORS[lvl] : 'var(--color-border)',
              opacity: levels.has(lvl) ? 1 : 0.5,
            }}
          >
            {lvl}
          </button>
        ))}

        {/* Service filter */}
        {services.length > 0 && (
          <select
            value={service}
            onChange={(e) => setService(e.target.value)}
            aria-label="Filter by service"
            className="text-xs px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-[var(--color-text)]"
          >
            <option value="">All services</option>
            {services.map((svc) => (
              <option key={svc} value={svc}>{svc}</option>
            ))}
          </select>
        )}

        {/* Search */}
        <input
          type="text"
          placeholder="Search logs..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          aria-label="Search log messages"
          className="text-xs px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-[var(--color-text)] flex-1 min-w-[120px]"
        />

        {/* Pause auto-scroll */}
        <button
          onClick={() => setPaused((v) => !v)}
          aria-label={paused ? 'Resume auto-scroll' : 'Pause auto-scroll'}
          className="px-2 py-1 text-xs font-medium rounded border border-[var(--color-border)] text-[var(--color-text-muted)] hover:text-[var(--color-text)] transition-colors"
        >
          {paused ? '▶ Resume' : '⏸ Pause'}
        </button>
      </div>

      {/* Log table */}
      <div ref={tableRef} className="flex-1 overflow-auto font-mono text-xs" role="table" aria-label="Log entries">
        {filtered.length === 0 ? (
          <div className="p-4 text-center text-[var(--color-text-muted)]">No log entries</div>
        ) : (
          filtered.map((entry) => (
            <LogRow key={entry.id} entry={entry} />
          ))
        )}
      </div>
    </div>
  )
}

function LogRow({ entry }: Readonly<{ entry: LogEntry }>) {
  const ts = new Date(entry.created_at).toLocaleTimeString()
  return (
    <div
      role="row"
      className="flex items-start gap-2 px-4 py-1 border-b border-[var(--color-border-subtle)] hover:bg-[var(--color-bg-hover)]"
      style={{ color: LEVEL_COLORS[entry.level] ?? 'var(--color-text)' }}
    >
      <span className="shrink-0 w-[70px] text-[var(--color-text-muted)]">{ts}</span>
      <span
        className="shrink-0 w-[50px] font-semibold"
        style={{ color: LEVEL_COLORS[entry.level] }}
      >
        {entry.level}
      </span>
      {entry.service && (
        <span className="shrink-0 w-[80px] text-[var(--color-text-muted)] truncate">
          [{entry.service}]
        </span>
      )}
      <span className="flex-1 break-all">{entry.message}</span>
    </div>
  )
}
