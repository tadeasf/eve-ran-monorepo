// EVE ESI API integration for character search and data fetching

import type { Kill } from '../../lib/types'

// Use Next.js API routes to avoid CORS issues
const API_BASE_URL = '/api'
const ESI_BASE_URL = process.env.NEXT_PUBLIC_EVE_ESI_BASE_URL || 'https://esi.evetech.net/latest'

export interface CharacterSearchResult {
    character?: number[]
}

export interface CharacterInfo {
    character_id: number
    name: string
    corporation_id: number
    alliance_id?: number
    birthday: string
    description?: string
    gender: string
    race_id: number
    bloodline_id: number
    ancestry_id: number
    security_status?: number
}

export interface CorporationInfo {
    corporation_id: number
    name: string
    ticker: string
    member_count: number
    alliance_id?: number
}

export interface AllianceInfo {
    alliance_id: number
    name: string
    ticker: string
    executor_corporation_id: number
}

// Backend character info structure (matches backend Character model)
export interface BackendCharacterInfo {
    id: number
    name: string
    security_status: number
    title: string
    race_id: number
}

/**
 * Search for characters by name using our backend API
 */
export async function searchCharactersByName(searchTerm: string): Promise<BackendCharacterInfo[]> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/characters/search?search=${encodeURIComponent(searchTerm)}`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`Backend API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error searching characters:', error)
        throw error
    }
}

/**
 * Get character information by ID
 */
export async function getCharacterInfo(characterId: number): Promise<CharacterInfo> {
    try {
        const response = await fetch(
            `${ESI_BASE_URL}/characters/${characterId}/`,
            {
                headers: {
                    'Accept': 'application/json',
                    'User-Agent': 'Tundragon Corporation Character Manager'
                }
            }
        )

        if (!response.ok) {
            throw new Error(`ESI API error: ${response.status} ${response.statusText}`)
        }

        const data = await response.json()
        return {
            character_id: characterId,
            ...data
        }
    } catch (error) {
        console.error('Error fetching character info:', error)
        throw error
    }
}

/**
 * Get corporation information by ID
 */
export async function getCorporationInfo(corporationId: number): Promise<CorporationInfo> {
    try {
        const response = await fetch(
            `${ESI_BASE_URL}/corporations/${corporationId}/`,
            {
                headers: {
                    'Accept': 'application/json',
                    'User-Agent': 'Tundragon Corporation Character Manager'
                }
            }
        )

        if (!response.ok) {
            throw new Error(`ESI API error: ${response.status} ${response.statusText}`)
        }

        const data = await response.json()
        return {
            corporation_id: corporationId,
            ...data
        }
    } catch (error) {
        console.error('Error fetching corporation info:', error)
        throw error
    }
}

/**
 * Get alliance information by ID
 */
export async function getAllianceInfo(allianceId: number): Promise<AllianceInfo> {
    try {
        const response = await fetch(
            `${ESI_BASE_URL}/alliances/${allianceId}/`,
            {
                headers: {
                    'Accept': 'application/json',
                    'User-Agent': 'Tundragon Corporation Character Manager'
                }
            }
        )

        if (!response.ok) {
            throw new Error(`ESI API error: ${response.status} ${response.statusText}`)
        }

        const data = await response.json()
        return {
            alliance_id: allianceId,
            ...data
        }
    } catch (error) {
        console.error('Error fetching alliance info:', error)
        throw error
    }
}

/**
 * Enhanced character search with full character information
 */
export async function searchCharactersWithInfo(searchTerm: string): Promise<BackendCharacterInfo[]> {
    try {
        // Our backend already returns full character info, so we can directly return it
        return await searchCharactersByName(searchTerm)
    } catch (error) {
        console.error('Error in enhanced character search:', error)
        throw error
    }
}

/**
 * Batch add characters to the system
 */
export async function batchAddCharacters(characterIds: number[]): Promise<{
    success_count: number;
    failure_count: number;
    total_processed: number;
    errors?: string[];
}> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/characters/batch`,
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(characterIds)
            }
        )

        if (!response.ok) {
            throw new Error(`Backend API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error batch adding characters:', error)
        throw error
    }
}

/**
 * Batch search characters by names (more efficient than individual searches)
 */
export async function batchSearchCharactersByNames(names: string[]): Promise<{
    [name: string]: BackendCharacterInfo[]
}> {
    const results: { [name: string]: BackendCharacterInfo[] } = {}
    
    // For now, we'll search each name individually
    // In the future, this could be optimized with a batch search endpoint
    for (const name of names) {
        try {
            const searchResults = await searchCharactersByName(name)
            results[name] = searchResults
        } catch (error) {
            console.error(`Error searching for character ${name}:`, error)
            results[name] = []
        }
    }
    
    return results
}

/**
 * Get all characters from the database
 */
export async function getAllCharacters(): Promise<BackendCharacterInfo[]> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/characters`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`Backend API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error fetching all characters:', error)
        throw error
    }
}

/**
 * Get character stats with filters (optimized for dashboard)
 */
export interface CharacterStatsResponse {
    character_id: number;
    name: string;
    kill_count: number;
    total_isk: number;
}

export async function getCharacterStats(params: {
    regionIDs?: number[];
    startDate?: string;
    endDate?: string;
}): Promise<CharacterStatsResponse[]> {
    try {
        const queryParams = new URLSearchParams();

        if (params.regionIDs && params.regionIDs.length > 0) {
            params.regionIDs.forEach(id => queryParams.append('regionID', id.toString()));
        }

        if (params.startDate) {
            queryParams.append('startDate', params.startDate);
        }

        if (params.endDate) {
            queryParams.append('endDate', params.endDate);
        }

        const response = await fetch(
            `${API_BASE_URL}/characters/stats?${queryParams.toString()}`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`Backend API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error fetching character stats:', error)
        throw error
    }
}

/**
 * Get kills with filters and pagination
 */
export interface KillsResponse {
    kills: Kill[];
    total_count: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export async function getKills(params: {
    characterId?: number;
    regionId?: number;
    startDate?: string;
    endDate?: string;
    page?: number;
    pageSize?: number;
}): Promise<KillsResponse> {
    try {
        const queryParams = new URLSearchParams();

        if (params.characterId) {
            queryParams.append('character_id', params.characterId.toString());
        }

        if (params.regionId) {
            queryParams.append('region_id', params.regionId.toString());
        }

        if (params.startDate) {
            queryParams.append('start_date', params.startDate);
        }

        if (params.endDate) {
            queryParams.append('end_date', params.endDate);
        }

        if (params.page) {
            queryParams.append('page', params.page.toString());
        }

        if (params.pageSize) {
            queryParams.append('page_size', params.pageSize.toString());
        }

        const response = await fetch(
            `${API_BASE_URL}/kills?${queryParams.toString()}`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`Backend API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error fetching kills:', error)
        throw error
    }
}

/**
 * Remove a character by ID
 */
export async function removeCharacter(characterId: number): Promise<void> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/characters/${characterId}`,
            {
                method: 'DELETE',
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`Backend API error: ${response.status} ${response.statusText}`)
        }
    } catch (error) {
        console.error('Error removing character:', error)
        throw error
    }
}

/**
 * Get dashboard statistics
 */
export async function getDashboardStats(): Promise<{
    total_characters: number;
    total_kills: number;
    total_regions: number;
    recent_activity: number;
}> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/admin/stats`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error fetching dashboard stats:', error)
        throw error
    }
}

/**
 * Competition API functions
 */

export interface CompetitionSettings {
    id?: number;
    metric: 'isk_destroyed' | 'kill_count';
    regions: number[];
    active: boolean;
    created_at?: string;
    updated_at?: string;
}

export interface CompetitionStanding {
    character_id: number;
    character_name: string;
    value: number;
    rank: number;
}

export interface CompetitionWinner {
    character_id: number;
    character_name: string;
    month: number;
    year: number;
    metric: string;
    value: number;
    rank: number;
}

/**
 * Get current competition settings
 */
export async function getCompetitionSettings(): Promise<CompetitionSettings> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/competition/settings`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`API error: ${response.status} ${response.statusText}`)
        }

        const data = await response.json()
        // Ensure regions is always an array, never null
        if (!data.regions) {
            data.regions = []
        }
        return data
    } catch (error) {
        console.error('Error fetching competition settings:', error)
        throw error
    }
}

/**
 * Update competition settings
 */
export async function updateCompetitionSettings(settings: CompetitionSettings): Promise<CompetitionSettings> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/competition/settings`,
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                body: JSON.stringify(settings)
            }
        )

        if (!response.ok) {
            throw new Error(`API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error updating competition settings:', error)
        throw error
    }
}

/**
 * Get current competition standings
 */
export async function getCurrentCompetitionStandings(): Promise<CompetitionStanding[]> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/competition/current`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`API error: ${response.status} ${response.statusText}`)
        }

        const data = await response.json()
        // Ensure we always return an array
        return Array.isArray(data) ? data : []
    } catch (error) {
        console.error('Error fetching current competition standings:', error)
        throw error
    }
}

/**
 * Get competition history for a specific month and year
 */
export async function getCompetitionHistory(month: number, year: number): Promise<CompetitionWinner[]> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/competition/history/${month}/${year}`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error fetching competition history:', error)
        throw error
    }
}

/**
 * Get all competition history
 */
export async function getAllCompetitionHistory(): Promise<CompetitionWinner[]> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/competition/history`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`API error: ${response.status} ${response.statusText}`)
        }

        return await response.json()
    } catch (error) {
        console.error('Error fetching all competition history:', error)
        throw error
    }
}

/**
 * Get recent competition winners (last month)
 */
export async function getRecentCompetitionWinners(): Promise<CompetitionWinner[]> {
    try {
        const response = await fetch(
            `${API_BASE_URL}/competition/recent-winners`,
            {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )

        if (!response.ok) {
            throw new Error(`API error: ${response.status} ${response.statusText}`)
        }

        const data = await response.json()
        // Ensure we always return an array
        return Array.isArray(data) ? data : []
    } catch (error) {
        console.error('Error fetching recent competition winners:', error)
        throw error
    }
}

/**
 * Get year-to-date competition winners (first place for each month in current year)
 */
export async function getYearToDateWinners(): Promise<CompetitionWinner[]> {
    const response = await fetch(`${API_BASE_URL}/competition/ytd-winners`, {
        credentials: 'include',
        headers: {
            Accept: 'application/json',
        },
    })

    if (!response.ok) {
        console.error('Failed to fetch year-to-date winners:', response.statusText)
        throw new Error('Failed to fetch year-to-date winners')
    }

    const data = await response.json()
    return data || []
}

export async function getAllRegions(): Promise<{ region_id: number; name: string; description: string; constellations: number[] }[]> {
    const response = await fetch(`${API_BASE_URL}/regions`, {
        credentials: 'include',
        headers: {
            Accept: 'application/json',
        },
    })

    if (!response.ok) {
        console.error('Failed to fetch regions:', response.statusText)
        throw new Error('Failed to fetch regions')
    }

    const data = await response.json()
    return data || []
}
