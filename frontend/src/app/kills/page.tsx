'use client';

import { useState } from 'react';
import { InfiniteKillList } from '../components/InfiniteKillList';
import { ShipTypeFilter } from '../components/ShipTypeFilter';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export default function KillsPage() {
  const [selectedShipTypes, setSelectedShipTypes] = useState<number[]>([]);
  const [characterId, setCharacterId] = useState<string>('');
  const [startDate, setStartDate] = useState<string>('');
  const [endDate, setEndDate] = useState<string>('');

  const characterIdNum = characterId ? parseInt(characterId, 10) : undefined;
  const selectedShipType = selectedShipTypes.length > 0 ? selectedShipTypes[0] : undefined;

  return (
    <div className="container mx-auto py-8 px-4 max-w-7xl">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Kill Feed</h1>
        <p className="text-muted-foreground">
          Browse all kills with advanced filtering and infinite scroll
        </p>
      </div>

      {/* Filters Card */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Filters</CardTitle>
          <CardDescription>
            Filter kills by character, date range, and ship type
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-3">
            <div>
              <Label htmlFor="character-id">Character ID</Label>
              <Input
                id="character-id"
                type="number"
                placeholder="e.g., 123456789"
                value={characterId}
                onChange={(e) => setCharacterId(e.target.value)}
              />
            </div>
            <div>
              <Label htmlFor="start-date">Start Date</Label>
              <Input
                id="start-date"
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
              />
            </div>
            <div>
              <Label htmlFor="end-date">End Date</Label>
              <Input
                id="end-date"
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
              />
            </div>
          </div>

          <ShipTypeFilter
            selectedShipTypes={selectedShipTypes}
            onShipTypesChange={setSelectedShipTypes}
          />
        </CardContent>
      </Card>

      {/* Kill List with Infinite Scroll */}
      <InfiniteKillList
        characterId={characterIdNum}
        startDate={startDate || undefined}
        endDate={endDate || undefined}
        shipTypeId={selectedShipType}
      />
    </div>
  );
}
