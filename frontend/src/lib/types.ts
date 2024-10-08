export interface Region {
  region_id: number;
  name: string;
  description: string;
  constellations: number[]; // Assuming constellations is an array of IDs
}

export interface CharacterStats {
  character_id: number;
  kill_count: number;
  total_isk: number;
  total_value: number;
  name: string;
}

export interface Character {
  id: number;
  name: string;
  security_status: number;
  title: string;
  race_id: number;
}

export interface Kill {
  ID: number;
  KillmailID: number;
  KillmailTime: string;
  SolarSystemID: number;
  CharacterID: number;
  Victim: Victim;
  Attackers: string; // This is a JSON string now
  ZkillData: Zkill;
}

export interface Victim {
  AllianceID: number;
  CharacterID: number;
  CorporationID: number;
  DamageTaken: number;
  ShipTypeID: number;
  Position: Position;
  Items: Item[];
}

export interface Attacker {
  AllianceID?: number;
  CharacterID: number;
  CorporationID: number;
  DamageDone: number;
  FinalBlow: boolean;
  SecurityStatus: number;
  ShipTypeID: number;
  WeaponTypeID: number;
}

export interface Zkill {
  ID: number;
  KillmailID: number;
  CharacterID: number;
  LocationID: number;
  Hash: string;
  FittedValue: number;
  DroppedValue: number;
  DestroyedValue: number;
  TotalValue: number;
  Points: number;
  NPC: boolean;
  Solo: boolean;
  Awox: boolean;
  Labels: string[] | null;
}

export interface Position {
  X: number;
  Y: number;
  Z: number;
}

export interface Item {
  ItemTypeID: number;
  QuantityDestroyed?: number;
  QuantityDropped?: number;
  Flag: number;
  Singleton: number;
}

export interface RegionKillsResponse {
  data: Kill[];
}

export interface ChartConfig {
  [key: string]: {
    label: string;
    color: string;
  };
}

export interface System {
  SystemID: number;
  ConstellationID: number;
  RegionID: number;
  Name: string;
  SecurityClass: string;
  SecurityStatus: number;
  StarID: number;
  Planets: string; // Change to string as it's stored as json.RawMessage in Go
  Stargates: string; // Change to string as it's stored as json.RawMessage in Go
  Stations: string; // Change to string as it's stored as json.RawMessage in Go
  Position: string; // Change to string as it's stored as json.RawMessage in Go
}

// New interfaces to replace 'any' types
export interface Planet {
  PlanetID: number;
  TypeID: number;
  Name: string;
}

export interface Stargate {
  StargateID: number;
  DestinationStargateID: number;
  DestinationSystemID: number;
  TypeID: number;
  Name: string;
}

export interface Station {
  StationID: number;
  TypeID: number;
  Name: string;
}

export interface ESIItem {
  TypeID: number;
  GroupID: number;
  Name: string;
  Description: string;
  Mass: number;
  Volume: number;
  Capacity: number;
  PortionSize: number;
  PackagedVolume: number;
  Published: boolean;
  Radius: number;
}

export interface ZKillboardItem extends Item {
  Items?: ZKillboardItem[];
}