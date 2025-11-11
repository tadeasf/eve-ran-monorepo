"use client"

import { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { 
  Users, 
  Trash2, 
  Search, 
  Loader2, 
  User, 
  Shield, 
  Star,
  ExternalLink,
  RefreshCw
} from "lucide-react"
import { getAllCharacters, removeCharacter, BackendCharacterInfo } from '@/app/lib/eveApi'
import { CharacterSearch } from './CharacterSearch'
import { toast } from 'sonner'

export function AdminCharacterManager() {
  const [characters, setCharacters] = useState<BackendCharacterInfo[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchCharacters = async () => {
    try {
      setError(null)
      const data = await getAllCharacters()
      setCharacters(data)
    } catch (err) {
      setError('Failed to fetch characters')
      console.error('Error fetching characters:', err)
    } finally {
      setIsLoading(false)
      setIsRefreshing(false)
    }
  }

  const handleRefresh = async () => {
    setIsRefreshing(true)
    await fetchCharacters()
  }

  const handleRemoveCharacter = async (characterId: number, characterName: string) => {
    if (!confirm(`Are you sure you want to remove ${characterName}?`)) {
      return
    }

    try {
      await removeCharacter(characterId)
      setCharacters(prev => prev.filter(char => char.id !== characterId))
      toast.success(`Successfully removed ${characterName}`)
    } catch (err) {
      console.error('Error removing character:', err)
      toast.error(`Failed to remove ${characterName}`)
    }
  }

  const handleCharacterAdded = (character: BackendCharacterInfo) => {
    // Add the character to the list if it's not already there
    setCharacters(prev => {
      const exists = prev.some(char => char.id === character.id)
      if (exists) return prev
      return [...prev, character]
    })
    toast.success(`Added ${character.name} to your character list`)
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

  useEffect(() => {
    fetchCharacters()
  }, [])

  if (error && isLoading) {
    return (
      <div className="text-center">
        <p className="text-destructive mb-4">{error}</p>
        <Button onClick={fetchCharacters}>Try Again</Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Character Management</h2>
          <p className="text-muted-foreground">
            Manage characters in your tracking system
          </p>
        </div>
        <Button 
          onClick={handleRefresh} 
          disabled={isRefreshing}
          variant="outline"
          className="flex items-center gap-2"
        >
          <RefreshCw className={`size-4 ${isRefreshing ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <Tabs defaultValue="existing" className="space-y-6">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="existing" className="flex items-center gap-2">
            <Users className="size-4" />
            Existing Characters ({characters.length})
          </TabsTrigger>
          <TabsTrigger value="search" className="flex items-center gap-2">
            <Search className="size-4" />
            Search & Add
          </TabsTrigger>
        </TabsList>

        <TabsContent value="existing" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="size-5" />
                Tracked Characters
              </CardTitle>
              <CardDescription>
                Characters currently being tracked in your system
              </CardDescription>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="size-8 animate-spin" />
                  <span className="ml-2">Loading characters...</span>
                </div>
              ) : characters.length === 0 ? (
                <div className="text-center py-8">
                  <Users className="size-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">No characters found</p>
                  <p className="text-sm text-muted-foreground">Add some characters using the Search & Add tab</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {characters.map((character) => (
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
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => window.open(`https://zkillboard.com/character/${character.id}/`, '_blank')}
                          className="flex items-center gap-2"
                        >
                          <ExternalLink className="size-4" />
                          zKillboard
                        </Button>

                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={() => handleRemoveCharacter(character.id, character.name)}
                          className="flex items-center gap-2"
                        >
                          <Trash2 className="size-4" />
                          Remove
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="search" className="space-y-6">
          <CharacterSearch 
            showAddButton={true}
            onCharacterSelect={handleCharacterAdded}
          />
        </TabsContent>
      </Tabs>
    </div>
  )
}