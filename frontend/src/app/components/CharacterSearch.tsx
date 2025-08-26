"use client"

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
    Search,
    Loader2,
    User,
    Shield,
    Star,
    Plus
} from 'lucide-react'
import { searchCharactersWithInfo, BackendCharacterInfo, batchAddCharacters } from '@/app/lib/eveApi'

interface CharacterSearchProps {
    onCharacterSelect?: (character: BackendCharacterInfo) => void
    showAddButton?: boolean
}

export function CharacterSearch({ onCharacterSelect, showAddButton = false }: CharacterSearchProps) {
    const [searchTerm, setSearchTerm] = useState('')
    const [isSearching, setIsSearching] = useState(false)
    const [searchResults, setSearchResults] = useState<BackendCharacterInfo[]>([])
    const [error, setError] = useState<string | null>(null)
    const [addingCharacters, setAddingCharacters] = useState<Set<number>>(new Set())

    const handleSearch = async (e: React.FormEvent) => {
        e.preventDefault()

        if (!searchTerm.trim()) {
            return
        }

        setIsSearching(true)
        setError(null)
        setSearchResults([])

        try {
            const results = await searchCharactersWithInfo(searchTerm.trim())
            setSearchResults(results)

            if (results.length === 0) {
                setError('No characters found matching your search')
            }
        } catch (err) {
            setError('Failed to search characters. Please try again.')
            console.error('Search error:', err)
        } finally {
            setIsSearching(false)
        }
    }

    const getCharacterAvatarUrl = (characterId: number, size: number = 64) => {
        return `https://images.evetech.net/characters/${characterId}/portrait?size=${size}`
    }

    const formatSecurityStatus = (status?: number) => {
        if (status === undefined) return 'Unknown'
        return status.toFixed(1)
    }

    const getSecurityStatusColor = (status?: number) => {
        if (status === undefined) return 'text-muted-foreground'
        if (status >= 0) return 'text-green-500'
        return 'text-red-500'
    }

    const handleAddCharacter = async (character: BackendCharacterInfo) => {
        setAddingCharacters(prev => new Set(prev).add(character.id))
        
        try {
            await batchAddCharacters([character.id])
            onCharacterSelect?.(character)
        } catch (error) {
            console.error('Error adding character:', error)
            setError(`Failed to add ${character.name}. Please try again.`)
        } finally {
            setAddingCharacters(prev => {
                const newSet = new Set(prev)
                newSet.delete(character.id)
                return newSet
            })
        }
    }

    return (
        <div className="space-y-6">
            {/* Search Form */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Search className="size-5" />
                        Character Search
                    </CardTitle>
                    <CardDescription>
                        Search for EVE Online characters by name using the official ESI API
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSearch} className="flex gap-2">
                        <Input
                            placeholder="Enter character name..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            disabled={isSearching}
                            className="flex-1"
                        />
                        <Button type="submit" disabled={isSearching || !searchTerm.trim()}>
                            {isSearching ? (
                                <Loader2 className="size-4 animate-spin" />
                            ) : (
                                <Search className="size-4" />
                            )}
                            {isSearching ? 'Searching...' : 'Search'}
                        </Button>
                    </form>

                    {error && (
                        <div className="mt-4 p-4 border border-destructive rounded-lg">
                            <p className="text-destructive text-sm">{error}</p>
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Search Results */}
            {searchResults.length > 0 && (
                <Card>
                    <CardHeader>
                        <CardTitle>Search Results</CardTitle>
                        <CardDescription>
                            Found {searchResults.length} character{searchResults.length !== 1 ? 's' : ''} matching &ldquo;{searchTerm}&rdquo;
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            {searchResults.map((character) => (
                                <div
                                    key={character.id}
                                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors"
                                >
                                    <div className="flex items-center space-x-4">
                                        <Avatar className="size-16">
                                            <AvatarImage
                                                src={getCharacterAvatarUrl(character.id, 128)}
                                                alt={character.name}
                                            />
                                            <AvatarFallback>
                                                <User className="size-8" />
                                            </AvatarFallback>
                                        </Avatar>

                                        <div className="space-y-2">
                                            <div>
                                                <h3 className="text-lg font-semibold">{character.name}</h3>
                                                <p className="text-sm text-muted-foreground font-mono">
                                                    ID: {character.id}
                                                </p>
                                            </div>

                                            <div className="flex flex-wrap gap-2">
                                                <Badge
                                                    variant="outline"
                                                    className={`flex items-center gap-1 ${getSecurityStatusColor(character.security_status)}`}
                                                >
                                                    <Star className="size-3" />
                                                    Sec: {formatSecurityStatus(character.security_status)}
                                                </Badge>

                                                <Badge variant="outline" className="flex items-center gap-1">
                                                    <User className="size-3" />
                                                    Race ID: {character.race_id}
                                                </Badge>

                                                {character.title && (
                                                    <Badge variant="outline" className="flex items-center gap-1">
                                                        <Shield className="size-3" />
                                                        {character.title}
                                                    </Badge>
                                                )}
                                            </div>
                                        </div>
                                    </div>

                                    <div className="flex gap-2">
                                        {showAddButton && (
                                            <Button
                                                size="sm"
                                                onClick={() => handleAddCharacter(character)}
                                                disabled={addingCharacters.has(character.id)}
                                                className="flex items-center gap-2"
                                            >
                                                {addingCharacters.has(character.id) ? (
                                                    <Loader2 className="size-4 animate-spin" />
                                                ) : (
                                                    <Plus className="size-4" />
                                                )}
                                                {addingCharacters.has(character.id) ? 'Adding...' : 'Add Character'}
                                            </Button>
                                        )}

                                        <Button
                                            variant="outline"
                                            size="sm"
                                            onClick={() => window.open(`https://zkillboard.com/character/${character.id}/`, '_blank')}
                                        >
                                            View on zKillboard
                                        </Button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Instructions */}
            <Card>
                <CardHeader>
                    <CardTitle>How to Use Character Search</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                        <div>
                            <h4 className="font-medium mb-2">Search Tips:</h4>
                            <ul className="space-y-1 text-muted-foreground">
                                <li>• Enter partial or full character names</li>
                                <li>• Search is case-insensitive</li>
                                <li>• Results limited to 10 characters</li>
                                <li>• Uses official EVE ESI API</li>
                            </ul>
                        </div>
                        <div>
                            <h4 className="font-medium mb-2">Character Information:</h4>
                            <ul className="space-y-1 text-muted-foreground">
                                <li>• Avatar and basic character details</li>
                                <li>• Corporation and alliance IDs</li>
                                <li>• Security status and creation date</li>
                                <li>• Direct links to external tools</li>
                            </ul>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
