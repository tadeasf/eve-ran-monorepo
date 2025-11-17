'use client';

import { useEffect, useRef } from 'react';
import { useInfiniteKills, type Kill } from '@/hooks/useInfiniteKills';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { KillListSkeleton } from '@/components/ui/skeleton';
import { formatDistanceToNow } from 'date-fns';
import { Loader2 } from 'lucide-react';

interface InfiniteKillListProps {
  characterId?: number;
  startDate?: string;
  endDate?: string;
  shipTypeId?: number;
}

export function InfiniteKillList({
  characterId,
  startDate,
  endDate,
  shipTypeId,
}: InfiniteKillListProps) {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading, isError } =
    useInfiniteKills({
      character_id: characterId,
      start_date: startDate,
      end_date: endDate,
      ship_type_id: shipTypeId,
    });

  const observerTarget = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
          fetchNextPage();
        }
      },
      { threshold: 0.1 }
    );

    const currentTarget = observerTarget.current;
    if (currentTarget) {
      observer.observe(currentTarget);
    }

    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget);
      }
    };
  }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

  const formatISK = (value: number) => {
    if (value >= 1e9) return `${(value / 1e9).toFixed(2)}B`;
    if (value >= 1e6) return `${(value / 1e6).toFixed(2)}M`;
    if (value >= 1e3) return `${(value / 1e3).toFixed(2)}K`;
    return value.toLocaleString();
  };

  if (isLoading) {
    return <KillListSkeleton />;
  }

  if (isError) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-destructive">
          Failed to load kills. Please try again later.
        </CardContent>
      </Card>
    );
  }

  const allKills = data?.pages.flatMap((page) => page.data) ?? [];

  if (allKills.length === 0) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-muted-foreground">
          No kills found matching your filters.
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-3">
      {allKills.map((kill: Kill) => (
        <Card key={kill.id} className="hover:shadow-md transition-shadow">
          <CardHeader className="pb-3">
            <div className="flex justify-between items-start gap-4">
              <div className="space-y-1 min-w-0">
                <CardTitle className="text-base font-semibold truncate">
                  Killmail #{kill.killmail_id}
                </CardTitle>
                <div className="flex flex-wrap gap-x-4 gap-y-1 text-sm text-muted-foreground">
                  <span>System: {kill.solar_system_id}</span>
                  <span>Ship: {kill.victim_ship_type_id}</span>
                  <span className="whitespace-nowrap">
                    {formatDistanceToNow(new Date(kill.killmail_time), {
                      addSuffix: true,
                    })}
                  </span>
                </div>
              </div>
              <div className="flex flex-col items-end gap-1 shrink-0">
                {kill.zkill_data && (
                  <>
                    <Badge variant="secondary" className="text-xs">
                      {formatISK(kill.zkill_data.total_value)} ISK
                    </Badge>
                    <span className="text-xs text-muted-foreground">
                      {kill.zkill_data.points} pts
                    </span>
                  </>
                )}
              </div>
            </div>
          </CardHeader>
        </Card>
      ))}

      {/* Intersection Observer Target */}
      <div ref={observerTarget} className="py-4">
        {isFetchingNextPage && (
          <div className="flex justify-center items-center gap-2 text-muted-foreground">
            <Loader2 className="h-4 w-4 animate-spin" />
            <span className="text-sm">Loading more kills...</span>
          </div>
        )}
        {!hasNextPage && allKills.length > 0 && (
          <div className="text-center text-sm text-muted-foreground">
            No more kills to load
          </div>
        )}
      </div>
    </div>
  );
}
