export interface Region {
  region_id: number
  name: string
  description: string
  constellations: number[]
}

export interface CharacterStats {
  character_id: number
  name: string
  kill_count: number
  total_value: number
  formatted_total_value: string
}

export interface Character {
  id: number
  name: string
  security_status: number
  title: string
  race_id: number
}

export interface KillmailData {
  killmail_id: number
  character_id: number
  killmail_time: string
  solar_system_id: number
  locationID: number
  hash: string
  fitted_value: number
  dropped_value: number
  destroyed_value: number
  total_value: number
  points: number
  npc: boolean
  solo: boolean
  awox: boolean
  victim: {
    character_id: number
    corporation_id: number
    faction_id?: number
    damage_taken: number
    ship_type_id: number
    items: Array<{
      item_type_id: number
      singleton: number
      quantity_destroyed?: number
      quantity_dropped?: number
      flag: number
    }>
    position: {
      x: number
      y: number
      z: number
    }
  }
  attackers: Array<{
    alliance_id?: number
    character_id: number
    corporation_id: number
    damage_done: number
    final_blow: boolean
    security_status: number
    ship_type_id: number
    weapon_type_id: number
  }>
}

export interface RegionKillsResponse {
  data: KillmailData[]
}