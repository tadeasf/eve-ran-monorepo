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
  ID: number
  KillmailID: number
  KillmailTime: string
  SolarSystemID: number
  CharacterID: number
  Victim: Victim
  Attackers: string
  ZkillData: Zkill
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
  ID: number
  KillmailID: number
  CharacterID: number
  LocationID: number
  Hash: string
  FittedValue: number
  DroppedValue: number
  DestroyedValue: number
  TotalValue: number
  Points: number
  NPC: boolean
  Solo: boolean
  Awox: boolean
  Labels: string[] | null
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