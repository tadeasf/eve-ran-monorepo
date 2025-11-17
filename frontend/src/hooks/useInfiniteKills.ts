import { useInfiniteQuery } from '@tanstack/react-query';
import { apiFetchJson } from '@/lib/api';

export interface Kill {
  id: number;
  killmail_id: number;
  killmail_time: string;
  solar_system_id: number;
  character_id: number;
  victim_ship_type_id: number;
  zkill_data?: {
    total_value: number;
    points: number;
    fitted_value: number;
    dropped_value: number;
    destroyed_value: number;
  };
}

export interface PaginatedKillsResponse {
  data: Kill[];
  total: number;
  limit: number;
  offset: number;
  page: number;
  totalPages: number;
}

export interface KillsQueryParams {
  character_id?: number;
  start_date?: string;
  end_date?: string;
  ship_type_id?: number;
  limit?: number;
}

export function useInfiniteKills(params: KillsQueryParams = {}) {
  const limit = params.limit || 50;

  return useInfiniteQuery({
    queryKey: ['kills', 'infinite', params],
    queryFn: async ({ pageParam = 0 }) => {
      const queryParams = new URLSearchParams({
        limit: limit.toString(),
        offset: pageParam.toString(),
      });

      if (params.character_id) {
        queryParams.append('character_id', params.character_id.toString());
      }
      if (params.start_date) {
        queryParams.append('start_date', params.start_date);
      }
      if (params.end_date) {
        queryParams.append('end_date', params.end_date);
      }
      if (params.ship_type_id) {
        queryParams.append('ship_type_id', params.ship_type_id.toString());
      }

      return apiFetchJson<PaginatedKillsResponse>(`/kills?${queryParams.toString()}`);
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => {
      const nextOffset = lastPage.offset + lastPage.limit;
      return nextOffset < lastPage.total ? nextOffset : undefined;
    },
    staleTime: 60 * 1000, // 1 minute
  });
}
