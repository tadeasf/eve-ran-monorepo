'use client';

import { useState } from 'react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { X } from 'lucide-react';

interface ShipType {
  type_id: number;
  name: string;
  group_name?: string;
}

interface ShipTypeFilterProps {
  selectedShipTypes: number[];
  onShipTypesChange: (shipTypes: number[]) => void;
}

// Common ship types for EVE Online
const COMMON_SHIP_TYPES: ShipType[] = [
  // Frigates
  { type_id: 587, name: 'Rifter', group_name: 'Frigate' },
  { type_id: 588, name: 'Punisher', group_name: 'Frigate' },
  { type_id: 589, name: 'Tormentor', group_name: 'Frigate' },
  { type_id: 590, name: 'Magnate', group_name: 'Frigate' },
  { type_id: 591, name: 'Atron', group_name: 'Frigate' },
  { type_id: 592, name: 'Tristan', group_name: 'Frigate' },
  { type_id: 593, name: 'Incursus', group_name: 'Frigate' },
  { type_id: 594, name: 'Imicus', group_name: 'Frigate' },
  { type_id: 595, name: 'Slasher', group_name: 'Frigate' },
  { type_id: 596, name: 'Breacher', group_name: 'Frigate' },
  { type_id: 597, name: 'Burst', group_name: 'Frigate' },
  { type_id: 598, name: 'Probe', group_name: 'Frigate' },
  
  // Destroyers
  { type_id: 16236, name: 'Coercer', group_name: 'Destroyer' },
  { type_id: 16238, name: 'Catalyst', group_name: 'Destroyer' },
  { type_id: 16240, name: 'Corax', group_name: 'Destroyer' },
  { type_id: 16242, name: 'Thrasher', group_name: 'Destroyer' },
  
  // Cruisers
  { type_id: 621, name: 'Rupture', group_name: 'Cruiser' },
  { type_id: 622, name: 'Stabber', group_name: 'Cruiser' },
  { type_id: 623, name: 'Maller', group_name: 'Cruiser' },
  { type_id: 624, name: 'Arbitrator', group_name: 'Cruiser' },
  { type_id: 625, name: 'Vexor', group_name: 'Cruiser' },
  { type_id: 626, name: 'Thorax', group_name: 'Cruiser' },
  { type_id: 627, name: 'Caracal', group_name: 'Cruiser' },
  { type_id: 628, name: 'Moa', group_name: 'Cruiser' },
  
  // Battlecruisers
  { type_id: 4302, name: 'Drake', group_name: 'Battlecruiser' },
  { type_id: 4308, name: 'Hurricane', group_name: 'Battlecruiser' },
  { type_id: 24688, name: 'Prophecy', group_name: 'Battlecruiser' },
  
  // Battleships
  { type_id: 638, name: 'Raven', group_name: 'Battleship' },
  { type_id: 639, name: 'Scorpion', group_name: 'Battleship' },
  { type_id: 640, name: 'Rokh', group_name: 'Battleship' },
  { type_id: 641, name: 'Hyperion', group_name: 'Battleship' },
  { type_id: 642, name: 'Megathron', group_name: 'Battleship' },
  { type_id: 643, name: 'Dominix', group_name: 'Battleship' },
  { type_id: 644, name: 'Armageddon', group_name: 'Battleship' },
  { type_id: 645, name: 'Apocalypse', group_name: 'Battleship' },
  { type_id: 24692, name: 'Abaddon', group_name: 'Battleship' },
  { type_id: 28659, name: 'Tempest', group_name: 'Battleship' },
  { type_id: 28710, name: 'Maelstrom', group_name: 'Battleship' },
  
  // Capital Ships
  { type_id: 671, name: 'Archon', group_name: 'Carrier' },
  { type_id: 19720, name: 'Aeon', group_name: 'Supercarrier' },
  { type_id: 19722, name: 'Wyvern', group_name: 'Supercarrier' },
  { type_id: 23773, name: 'Naglfar', group_name: 'Dreadnought' },
  { type_id: 23913, name: 'Moros', group_name: 'Dreadnought' },
  { type_id: 23917, name: 'Phoenix', group_name: 'Dreadnought' },
  { type_id: 24483, name: 'Revelation', group_name: 'Dreadnought' },
  { type_id: 3514, name: 'Revenant', group_name: 'Supercarrier' },
  { type_id: 11567, name: 'Ragnarok', group_name: 'Titan' },
  { type_id: 671, name: 'Avatar', group_name: 'Titan' },
  { type_id: 3764, name: 'Leviathan', group_name: 'Titan' },
  { type_id: 42124, name: 'Keepstar', group_name: 'Citadel' },
];

export function ShipTypeFilter({ selectedShipTypes, onShipTypesChange }: ShipTypeFilterProps) {
  const [selectedType, setSelectedType] = useState<string>('');

  const handleAddShipType = (typeIdStr: string) => {
    const typeId = parseInt(typeIdStr, 10);
    if (typeId && !selectedShipTypes.includes(typeId)) {
      onShipTypesChange([...selectedShipTypes, typeId]);
    }
    setSelectedType('');
  };

  const handleRemoveShipType = (typeId: number) => {
    onShipTypesChange(selectedShipTypes.filter((id) => id !== typeId));
  };

  const handleClearAll = () => {
    onShipTypesChange([]);
  };

  const getShipTypeName = (typeId: number) => {
    const shipType = COMMON_SHIP_TYPES.find((st) => st.type_id === typeId);
    return shipType ? shipType.name : `Ship Type ${typeId}`;
  };

  // Group ship types by category
  const groupedShipTypes = COMMON_SHIP_TYPES.reduce((acc, shipType) => {
    const group = shipType.group_name || 'Other';
    if (!acc[group]) {
      acc[group] = [];
    }
    acc[group].push(shipType);
    return acc;
  }, {} as Record<string, ShipType[]>);

  return (
    <div className="space-y-3">
      <div className="flex items-end gap-2">
        <div className="flex-1">
          <Label htmlFor="ship-type-select">Filter by Ship Type</Label>
          <Select value={selectedType} onValueChange={handleAddShipType}>
            <SelectTrigger id="ship-type-select">
              <SelectValue placeholder="Select ship types..." />
            </SelectTrigger>
            <SelectContent>
              {Object.entries(groupedShipTypes).map(([groupName, ships]) => (
                <div key={groupName}>
                  <div className="px-2 py-1.5 text-sm font-semibold text-muted-foreground">
                    {groupName}
                  </div>
                  {ships.map((ship) => (
                    <SelectItem
                      key={ship.type_id}
                      value={ship.type_id.toString()}
                      disabled={selectedShipTypes.includes(ship.type_id)}
                    >
                      {ship.name}
                    </SelectItem>
                  ))}
                </div>
              ))}
            </SelectContent>
          </Select>
        </div>
        {selectedShipTypes.length > 0 && (
          <Button variant="outline" size="sm" onClick={handleClearAll}>
            Clear All
          </Button>
        )}
      </div>

      {/* Selected Ship Types */}
      {selectedShipTypes.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {selectedShipTypes.map((typeId) => (
            <div
              key={typeId}
              className="inline-flex items-center gap-1 rounded-md bg-secondary px-2.5 py-0.5 text-sm"
            >
              <span>{getShipTypeName(typeId)}</span>
              <button
                onClick={() => handleRemoveShipType(typeId)}
                className="rounded-sm hover:bg-secondary-foreground/20"
              >
                <X className="h-3 w-3" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
