/** Faction category from the API. */
export type FactionCategory = Readonly<{
  id: number
  name: string
  display_name?: string
  description?: string
  max_memberships?: number
  auto_join?: boolean
}>

/** Faction record from the API. */
export type Faction = Readonly<{
  id: number
  name: string
  display_name?: string
  description?: string
  category_id?: number
  standing?: number
  is_universal?: boolean
  members?: number[]
  member_count?: number
  category?: { id: number; name: string; display_name?: string }
  required_tags?: string[]
  member_tags?: string[]
  created_at?: string
}>

/** Form state for creating/editing a faction. */
export type FactionForm = Readonly<{
  name: string
  display_name: string
  description: string
  category_id: number | ''
  standing: number
  is_universal: boolean
  member_tags: string[]
}>

export const EMPTY_FORM: FactionForm = {
  name: '',
  display_name: '',
  description: '',
  category_id: '',
  standing: 0,
  is_universal: false,
  member_tags: [],
}

export function factionToForm(f: Faction): FactionForm {
  return {
    name: f.name,
    display_name: f.display_name || f.name,
    description: f.description || '',
    category_id: f.category_id ?? '',
    standing: f.standing ?? 0,
    is_universal: f.is_universal ?? false,
    member_tags: f.member_tags ?? [],
  }
}