const API = window.location.origin;

export type ResourceEndpoint = Readonly<{
  resourceType: string
  path: string
  apiBase: string
}>

export const RESOURCE_ENDPOINTS = {
  characters: {
    resourceType: 'characters',
    path: 'characters',
    apiBase: `${API}/api`,
  },
  rooms: {
    resourceType: 'rooms',
    path: 'rooms',
    apiBase: API,
  },
  races: {
    resourceType: 'races',
    path: 'races',
    apiBase: `${API}/api`,
  },
  abilities: {
    resourceType: 'abilities',
    path: 'abilities',
    apiBase: `${API}/api`,
  },
  factions: {
    resourceType: 'factions',
    path: 'factions',
    apiBase: `${API}/api`,
  },
  factionCategories: {
    resourceType: 'faction-categories',
    path: 'faction-categories',
    apiBase: `${API}/api`,
  },
  npcTemplates: {
    resourceType: 'npc-templates',
    path: 'npc-templates',
    apiBase: `${API}/api`,
  },
  quests: {
    resourceType: 'quests',
    path: 'quests',
    apiBase: `${API}/api`,
  },
  effectDefs: {
    resourceType: 'effect-defs',
    path: 'effect-defs',
    apiBase: `${API}/api`,
  },
} as const;