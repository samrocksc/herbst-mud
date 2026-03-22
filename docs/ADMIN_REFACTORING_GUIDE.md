# Admin Panel Refactoring Guide

## Current State Analysis

The admin panel has grown organically and now has several issues:

### File Size Violations
| File | Lines | Status |
|------|-------|--------|
| `routes/map.tsx` | ~500 | ❌ Needs refactoring |
| `routes/items.tsx` | ~560 | ❌ Needs refactoring |
| `routes/npcs.tsx` | ~400 | ❌ Needs refactoring |
| `routes/dashboard.tsx` | ~120 | ⚠️ Acceptable |
| `routes/login.tsx` | ~90 | ✅ OK |
| `routes/index.tsx` | ~30 | ✅ OK |

### Issues
1. **Monolithic Components**: Each route file mixes data fetching, state management, and rendering
2. **Duplicated Logic**: Similar CRUD patterns repeated in NPCs and Items
3. **Hardcoded URLs**: `http://localhost:8080` scattered throughout
4. **Inline Styles**: Tailwind classes are long and make files harder to read
5. **No Error Boundaries**: API errors aren't handled gracefully

---

## Target Architecture

```
admin/src/
├── api/
│   ├── client.ts          # Base fetch wrapper with auth
│   ├── rooms.ts           # Room API functions
│   ├── npcs.ts            # NPC API functions
│   ├── items.ts           # Item API functions
│   └── auth.ts            # Auth API functions
├── components/
│   ├── ui/
│   │   ├── Button.tsx     # Reusable button
│   │   ├── Input.tsx      # Reusable input
│   │   ├── Select.tsx     # Reusable select
│   │   ├── Modal.tsx      # ✅ Already exists
│   │   └── Panel.tsx      # Reusable panel
│   ├── map/
│   │   ├── MapCanvas.tsx  # Room visualization
│   │   ├── RoomNode.tsx   # Single room node
│   │   ├── RoomPanel.tsx  # Room details panel
│   │   ├── ZLevelNav.tsx  # Z-level navigation
│   │   └── CreateRoomModal.tsx
│   ├── npcs/
│   │   ├── NPCList.tsx
│   │   ├── NPCForm.tsx
│   │   └── NPCPanel.tsx
│   └── items/
│       ├── ItemList.tsx
│       ├── ItemForm.tsx
│       └── ItemPanel.tsx
├── hooks/
│   ├── useAuth.ts         # Auth state & login/logout
│   ├── useRooms.ts        # Room data fetching
│   ├── useNPCs.ts         # NPC data fetching
│   └── useItems.ts        # Item data fetching
├── routes/
│   ├── map.tsx            # Composes map components
│   ├── npcs.tsx           # Composes NPC components
│   ├── items.tsx          # Composes item components
│   ├── dashboard.tsx      # Stats dashboard
│   ├── login.tsx          # Login page
│   └── index.tsx          # Home/redirect
├── styles/
│   └── theme.ts           # Color constants
└── main.tsx               # Entry point
```

---

## Step-by-Step Refactoring

### Step 1: Create API Client Layer

**File: `src/api/client.ts`**
```typescript
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

interface FetchOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  body?: unknown
  headers?: Record<string, string>
}

export async function apiFetch<T>(
  endpoint: string,
  options: FetchOptions = {}
): Promise<T> {
  const token = localStorage.getItem('token')
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method: options.method || 'GET',
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(error.error || 'Request failed')
  }

  return response.json()
}
```

**File: `src/api/rooms.ts`**
```typescript
import { apiFetch } from './client'

export interface Room {
  id: number
  name: string
  description: string
  isStartingRoom?: boolean
  exits: Record<string, number>
  atmosphere?: string
}

export const roomsApi = {
  list: () => apiFetch<Room[]>('/rooms'),
  get: (id: number) => apiFetch<Room>(`/rooms/${id}`),
  create: (room: Partial<Room>) => apiFetch<Room>('/rooms', { method: 'POST', body: room }),
  update: (id: number, room: Partial<Room>) => apiFetch<Room>(`/rooms/${id}`, { method: 'PUT', body: room }),
  delete: (id: number) => apiFetch<void>(`/rooms/${id}`, { method: 'DELETE' }),
}
```

**File: `src/api/npcs.ts`**
```typescript
import { apiFetch } from './client'

export interface NPC {
  id: number
  name: string
  class: string
  race: string
  level: number
  currentRoomId: number
  isNPC: boolean
}

export const npcsApi = {
  list: () => apiFetch<{ npcs: NPC[] }>('/npcs'),
  get: (id: number) => apiFetch<NPC>(`/characters/${id}`),
  create: (npc: Partial<NPC>) => apiFetch<NPC>('/characters', { method: 'POST', body: { ...npc, isNPC: true } }),
  update: (id: number, npc: Partial<NPC>) => apiFetch<NPC>(`/characters/${id}`, { method: 'PATCH', body: npc }),
  delete: (id: number) => apiFetch<void>(`/characters/${id}`, { method: 'DELETE' }),
}
```

**File: `src/api/items.ts`**
```typescript
import { apiFetch } from './client'

export interface Item {
  id: number
  name: string
  description: string
  slot: string
  level: number
  weight: number
  isEquipped: boolean
  isImmovable: boolean
  isVisible: boolean
  itemType: string
}

export const itemsApi = {
  list: () => apiFetch<Item[]>('/equipment'),
  get: (id: number) => apiFetch<Item>(`/equipment/${id}`),
  create: (item: Partial<Item>) => apiFetch<Item>('/equipment', { method: 'POST', body: item }),
  update: (id: number, item: Partial<Item>) => apiFetch<Item>(`/equipment/${id}`, { method: 'PUT', body: item }),
  delete: (id: number) => apiFetch<void>(`/equipment/${id}`, { method: 'DELETE' }),
}
```

---

### Step 2: Create Custom Hooks

**File: `src/hooks/useRooms.ts`**
```typescript
import { useState, useEffect, useCallback } from 'react'
import { roomsApi, Room } from '../api/rooms'

export function useRooms() {
  const [rooms, setRooms] = useState<Room[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async () => {
    try {
      setLoading(true)
      const data = await roomsApi.list()
      setRooms(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load rooms')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const create = useCallback(async (room: Partial<Room>) => {
    const newRoom = await roomsApi.create(room)
    setRooms(prev => [...prev, newRoom])
    return newRoom
  }, [])

  const update = useCallback(async (id: number, room: Partial<Room>) => {
    const updated = await roomsApi.update(id, room)
    setRooms(prev => prev.map(r => r.id === id ? updated : r))
    return updated
  }, [])

  const remove = useCallback(async (id: number) => {
    await roomsApi.delete(id)
    setRooms(prev => prev.filter(r => r.id !== id))
  }, [])

  return { rooms, loading, error, create, update, remove, reload: load }
}
```

**File: `src/hooks/useNPCs.ts`**
```typescript
import { useState, useEffect, useCallback } from 'react'
import { npcsApi, NPC } from '../api/npcs'

export function useNPCs() {
  const [npcs, setNpcs] = useState<NPC[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async () => {
    try {
      setLoading(true)
      const data = await npcsApi.list()
      setNpcs(data.npcs || [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load NPCs')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const create = useCallback(async (npc: Partial<NPC>) => {
    const newNPC = await npcsApi.create(npc)
    setNpcs(prev => [...prev, newNPC])
    return newNPC
  }, [])

  const update = useCallback(async (id: number, npc: Partial<NPC>) => {
    const updated = await npcsApi.update(id, npc)
    setNpcs(prev => prev.map(n => n.id === id ? updated : n))
    return updated
  }, [])

  const remove = useCallback(async (id: number) => {
    await npcsApi.delete(id)
    setNpcs(prev => prev.filter(n => n.id !== id))
  }, [])

  return { npcs, loading, error, create, update, remove, reload: load }
}
```

**File: `src/hooks/useItems.ts`**
```typescript
import { useState, useEffect, useCallback } from 'react'
import { itemsApi, Item } from '../api/items'

export function useItems() {
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async () => {
    try {
      setLoading(true)
      const data = await itemsApi.list()
      setItems(Array.isArray(data) ? data : [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load items')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const create = useCallback(async (item: Partial<Item>) => {
    const newItem = await itemsApi.create(item)
    setItems(prev => [...prev, newItem])
    return newItem
  }, [])

  const update = useCallback(async (id: number, item: Partial<Item>) => {
    const updated = await itemsApi.update(id, item)
    setItems(prev => prev.map(i => i.id === id ? updated : i))
    return updated
  }, [])

  const remove = useCallback(async (id: number) => {
    await itemsApi.delete(id)
    setItems(prev => prev.filter(i => i.id !== id))
  }, [])

  return { items, loading, error, create, update, remove, reload: load }
}
```

---

### Step 3: Create Theme Constants

**File: `src/styles/theme.ts`**
```typescript
// Ninja Turtle-inspired brown/green/white theme
export const colors = {
  // Backgrounds
  bgDark: '#1a1612',
  bgMedium: '#2d2416',
  bgLight: '#3d3020',

  // Primary (Green)
  primary: '#4a7c4e',
  primaryLight: '#5a9c5e',
  primaryDark: '#3a5c3e',

  // Accent (Brown)
  accent: '#8b7355',
  accentLight: '#a89070',

  // Text
  text: '#e8dcc4',
  textMuted: '#a89070',

  // Borders
  border: '#5a4a35',
  borderLight: '#8b7355',

  // Status
  success: '#4a7c4e',
  danger: '#8b4444',
  warning: '#a87044',
} as const

export const layout = {
  sidebarWidth: '280px',
  panelWidth: '320px',
  headerHeight: '50px',
} as const
```

---

### Step 4: Create Reusable UI Components

**File: `src/components/ui/Button.tsx`**
```typescript
import { ButtonHTMLAttributes, ReactNode } from 'react'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  children: ReactNode
}

export function Button({ variant = 'primary', className = '', children, ...props }: ButtonProps) {
  const baseStyles = 'px-4 py-2 rounded border-none cursor-pointer font-medium transition-colors'

  const variants = {
    primary: 'bg-[#4a7c4e] text-[#e8dcc4] hover:bg-[#5a9c5e]',
    secondary: 'bg-[#3d3020] text-[#a89070] border border-[#5a4a35] hover:bg-[#4d4030]',
    danger: 'bg-[#8b4444] text-[#e8dcc4] hover:bg-[#a84444]',
    ghost: 'bg-transparent text-[#a89070] hover:text-[#e8dcc4]',
  }

  return (
    <button className={`${baseStyles} ${variants[variant]} ${className}`} {...props}>
      {children}
    </button>
  )
}
```

**File: `src/components/ui/Input.tsx`**
```typescript
import { InputHTMLAttributes } from 'react'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
}

export function Input({ label, className = '', ...props }: InputProps) {
  return (
    <div className="mb-3">
      {label && <label className="block text-[#a89070] text-xs mb-1">{label}</label>}
      <input
        className={`w-full p-2 bg-[#3d3020] border border-[#5a4a35] rounded text-[#e8dcc4] text-sm ${className}`}
        {...props}
      />
    </div>
  )
}
```

---

### Step 5: Refactor Route Files

**File: `src/routes/items.tsx` (After refactoring)**
```typescript
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useItems } from '../hooks/useItems'
import { ItemList } from '../components/items/ItemList'
import { ItemPanel } from '../components/items/ItemPanel'
import { ItemForm } from '../components/items/ItemForm'
import { Modal } from '../components/Modal'

export const Route = createFileRoute('/items')({
  component: ItemManager,
})

function ItemManager() {
  const { items, loading, error, create, update, remove } = useItems()
  const [selectedItem, setSelectedItem] = useState(null)
  const [editingItem, setEditingItem] = useState(null)
  const [showCreate, setShowCreate] = useState(false)

  if (loading) return <div className="p-8 text-[#e8dcc4]">Loading...</div>
  if (error) return <div className="p-8 text-[#8b4444]">Error: {error}</div>

  return (
    <div className="flex h-screen bg-[#1a1612]">
      <ItemList items={items} onSelect={setSelectedItem} onCreate={() => setShowCreate(true)} />

      <div className="flex-1 p-6">
        {editingItem ? (
          <ItemForm item={editingItem} onSave={update} onCancel={() => setEditingItem(null)} />
        ) : selectedItem ? (
          <ItemPanel item={selectedItem} onEdit={() => setEditingItem(selectedItem)} onDelete={remove} />
        ) : (
          <EmptyState />
        )}
      </div>

      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Create Item">
        <ItemForm onSave={create} onCancel={() => setShowCreate(false)} />
      </Modal>
    </div>
  )
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center h-full text-[#a89070]">
      <p>Select an item or create a new one</p>
    </div>
  )
}
```

---

## Migration Checklist

### Phase A: API Layer
- [ ] Create `src/api/client.ts`
- [ ] Create `src/api/rooms.ts`
- [ ] Create `src/api/npcs.ts`
- [ ] Create `src/api/items.ts`
- [ ] Create `src/api/auth.ts`
- [ ] Add `VITE_API_URL` to `.env`

### Phase B: Hooks
- [ ] Create `src/hooks/useRooms.ts`
- [ ] Create `src/hooks/useNPCs.ts`
- [ ] Create `src/hooks/useItems.ts`
- [ ] Create `src/hooks/useAuth.ts`

### Phase C: UI Components
- [ ] Create `src/components/ui/Button.tsx`
- [ ] Create `src/components/ui/Input.tsx`
- [ ] Create `src/components/ui/Select.tsx`
- [ ] Create `src/components/ui/Panel.tsx`
- [ ] Create `src/styles/theme.ts`

### Phase D: Map Refactoring
- [ ] Create `src/components/map/MapCanvas.tsx`
- [ ] Create `src/components/map/RoomNode.tsx`
- [ ] Create `src/components/map/RoomPanel.tsx`
- [ ] Create `src/components/map/ZLevelNav.tsx`
- [ ] Refactor `src/routes/map.tsx`

### Phase E: NPC Refactoring
- [ ] Create `src/components/npcs/NPCList.tsx`
- [ ] Create `src/components/npcs/NPCForm.tsx`
- [ ] Create `src/components/npcs/NPCPanel.tsx`
- [ ] Refactor `src/routes/npcs.tsx`

### Phase F: Item Refactoring
- [ ] Create `src/components/items/ItemList.tsx`
- [ ] Create `src/components/items/ItemForm.tsx`
- [ ] Create `src/components/items/ItemPanel.tsx`
- [ ] Refactor `src/routes/items.tsx`

---

## Environment Variables

Create `admin/.env`:
```env
VITE_API_URL=http://localhost:8080
```

Create `admin/.env.production`:
```env
VITE_API_URL=https://api.yourdomain.com
```

---

## Testing After Each Phase

1. Run `npm run build` - should complete without errors
2. Run `npm run dev` - app should load
3. Test each affected route:
   - Navigate to page
   - Test CRUD operations
   - Test form submissions
4. Verify no regressions in other routes