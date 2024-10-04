export interface Region {
  region_id: number
  name: string
  description: string
  constellations: number[]
}

export interface CharacterStats {
  character_id: number
  kill_count: number
  total_isk: number
  total_value: number
  name: string
}

export interface Character {
  id: number
  name: string
  security_status: number
  title: string
  race_id: number
}

export interface Kill {
  id: number
  killmail_id: number
  killmail_time: string
  solar_system_id: number
  victim: Victim
  attackers: string // Base64 encoded JSON string
  zkill_data: Zkill
}

export interface Victim {
  alliance_id: number
  character_id: number
  corporation_id: number
  damage_taken: number
  ship_type_id: number
  position: Position
  items: Item[]
}

export interface Attacker {
  alliance_id?: number
  character_id: number
  corporation_id: number
  damage_done: number
  final_blow: boolean
  security_status: number
  ship_type_id: number
  weapon_type_id: number
}

export interface Zkill {
  id: number
  killmail_id: number
  character_id: number
  location_id: number
  hash: string
  fitted_value: number
  dropped_value: number
  destroyed_value: number
  total_value: number
  points: number
  npc: boolean
  solo: boolean
  awox: boolean
  labels: string[]
}

export interface Position {
  x: number
  y: number
  z: number
}

export interface Item {
  item_type_id: number
  quantity_destroyed?: number
  quantity_dropped?: number
  flag: number
  singleton: number
}

export interface RegionKillsResponse {
  data: Kill[]
}

export interface ChartConfig {
  [key: string]: {
    label: string
    color: string
  }
}