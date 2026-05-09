import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect, useRef, useCallback } from 'react'
import { useLogs, useLogServices } from '../../hooks/useLogs'
import { PageHeader } from '../../components/PageHeader'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/logs')({
  component: LogsPage,
})

const LEVEL_COLORS: Record<string, string> = {
  DEBUG: 'bg-gray-500/20 text-gray-400',
  INFO: 'bg-blue-500/20 text-blue-400',
  WARN: 'bg-yellow-500/20 text-yellow-400',
  ERROR: 'bg-red-500/20 text-red-400',
}

const LEVELS = ['ALL', 'DEBUG', 'INFO', 'WARN', 'ERROR'] as const

type LogEntryLike = {
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
  const [liveLines, setLiveLines] = useState<LogEntryLike[]>([])
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

  const handleLiveToggle = useCallback(() => {
    setLive((prev) => !prev)
    setLiveLines([])
  }, [])

  useEffect(() => {
    if (!live) return
    const token = localStorage.getItem('auth_token') || ''
    const es = new EventSource(`${window.location.origin}/api/logs/stream?token=${encodeURIComponent(token)}`)
    es.onmessage = (e) => {
      try {
        const entry = JSON.parse(e.data) as LogEntryLike
        setLiveLines((prev) => [entry, ...prev].slice(0, 500))
      } catch { /* skip */ }
    }
    es.onerror = () => { es.close() }
    return () => { es.close() }
  }, [live])

  return (
    <div className="management-page">
      <PageHeader title="Logs" backTo="/dashboard" />
      <div className="filters-bar flex flex-wrap items-center gap-2 mb-4">
        <div className="flex gap-1">
          {LEVELS.map((l) => (
            <button
              key={l}
              onClick={() => setLevel(l === 'ALL' ? '' : l)}
              className={`px-2 py-1 rounded text-xs font-medium transition-colors ${
                (l === 'ALL' && !level) || l === level ? 'bg-primary text-white' : 'bg-surface-muted text-text-muted hover:bg-surface-muted/70'
              }`}
            >
              {l}
            </button>
          ))}
        </div>
        <select value={service} onChange={(e) => setService(e.target.value)} className="bg-surface border border-border rounded px-2 py-1 text-xs text-text">
          <option value="">All Services</option>
          {(services ?? []).map((s) => (<option key={s} value={s}>{s}</option>))}
        </select>
        <input type="text" value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Search logs..." className="bg-surface border border-border rounded px-2 py-1 text-xs text-text flex-1 min-w-[150px]" />
        <Button variant={live ? 'primary' : 'secondary'} size="sm" onClick={handleLiveToggle}>
          {live ? 'Live ON' : 'Live OFF'}
        </Button>
      </div>
      {isLoading && !live && <div className="loading">Loading logs...</div>}
      <div className="overflow-x-auto">
        <table className="w-full text-xs">
          <thead>
            <tr className="border-b border-border text-text-muted">
              <th className="text-left py-1 px-2 font-medium">Time</th>
              <th className="text-left py-1 px-2 font-medium w-16">Level</th>
              <th className="text-left py-1 px-2 font-medium w-24">Service</th>
              <th className="text-left py-1 px-2 font-medium">Message</th>
              <th className="text-left py-1 px-2 font-medium w-16">Char</th>
              <th className="text-left py-1 px-2 font-medium w-16">Room</th>
            </tr>
          </thead>
          <tbody>
            {filtered.slice(0, 500).map((log, i) => (
              <tr key={log.id ?? i} className="border-b border-border/30 hover:bg-surface-muted/50">
                <td className="py-1 px-2 text-text-muted whitespace-nowrap">{formatTime(log.created_at)}</td>
                <td className="py-1 px-2"><span className={`px-1.5 py-0.5 rounded text-xs font-medium ${LEVEL_COLORS[log.level] ?? 'bg-gray-500/20 text-gray-400'}`}>{log.level}</span></td>
                <td className="py-1 px-2 text-text-muted">{log.service || '—'}</td>
                <td className="py-1 px-2 text-text max-w-md truncate" title={log.message}>{log.message}</td>
                <td className="py-1 px-2 text-text-muted">{log.character_id ?? '—'}</td>
                <td className="py-1 px-2 text-text-muted">{log.room_id ?? '—'}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {filtered.length === 0 && <div className="text-center py-8 text-text-muted">No logs found.</div>}
      </div>
      <div ref={bottomRef} />
    </div>
  )
}

function formatTime(t: string): string {
  if (!t) return '—'
  try { return new Date(t).toLocaleTimeString() } catch { return t }
}