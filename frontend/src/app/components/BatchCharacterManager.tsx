"use client"

import { useState } from 'react'
import { Button } from '@/app/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/app/components/ui/card'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
    Plus,
    Upload,
    Check,
    X,
    AlertCircle,
    Users,
    Loader2,
    Hash,
    User,
    Search
} from 'lucide-react'
import { batchAddCharacters, searchCharactersByName } from '@/app/lib/eveApi'

interface ProcessingResult {
    input: string
    characterId?: number
    characterName?: string
    status: 'success' | 'error' | 'pending' | 'resolving'
    message?: string
}

interface BatchCharacterManagerProps {
    onBatchComplete?: () => void
}

export function BatchCharacterManager({ onBatchComplete }: BatchCharacterManagerProps = {}) {
    const [inputType, setInputType] = useState<'ids' | 'names'>('names')
    const [characterInput, setCharacterInput] = useState('')
    const [isProcessing, setIsProcessing] = useState(false)
    const [results, setResults] = useState<ProcessingResult[]>([])
    const [progress, setProgress] = useState(0)

    const handleBatchAddByIds = async (ids: number[]) => {
        try {
            const response = await batchAddCharacters(ids)
            const processResults: ProcessingResult[] = []

            // Mark successful characters
            for (let i = 0; i < response.success_count; i++) {
                if (i < ids.length) {
                    processResults.push({
                        input: ids[i].toString(),
                        characterId: ids[i],
                        status: 'success',
                        message: 'Character added successfully'
                    })
                }
            }

            // Mark failed characters
            const failedStartIndex = response.success_count
            for (let i = 0; i < response.failure_count; i++) {
                const characterIndex = failedStartIndex + i
                if (characterIndex < ids.length) {
                    const errorMessage = response.errors && response.errors[i]
                        ? response.errors[i]
                        : 'Failed to add character'

                    processResults.push({
                        input: ids[characterIndex].toString(),
                        characterId: ids[characterIndex],
                        status: 'error',
                        message: errorMessage
                    })
                }
            }

            return processResults
        } catch (error) {
            return ids.map(id => ({
                input: id.toString(),
                characterId: id,
                status: 'error' as const,
                message: error instanceof Error ? error.message : 'Unknown error occurred'
            }))
        }
    }

    const handleBatchAddByNames = async (names: string[]) => {
        const processResults: ProcessingResult[] = []
        const characterIdsToAdd: number[] = []

        // First, resolve all names to IDs
        for (let i = 0; i < names.length; i++) {
            const name = names[i]
            setProgress((i / (names.length * 2)) * 100) // First half for name resolution

            try {
                // Update status to resolving
                const currentResults = [...processResults]
                currentResults.push({
                    input: name,
                    status: 'resolving',
                    message: 'Resolving character name...'
                })
                setResults(currentResults)

                const searchResults = await searchCharactersByName(name)
                
                if (searchResults.length === 0) {
                    processResults.push({
                        input: name,
                        status: 'error',
                        message: 'Character not found'
                    })
                } else if (searchResults.length > 1) {
                    // Multiple matches - use exact match if available
                    const exactMatch = searchResults.find(char => 
                        char.name.toLowerCase() === name.toLowerCase()
                    )
                    
                    if (exactMatch) {
                        characterIdsToAdd.push(exactMatch.id)
                        processResults.push({
                            input: name,
                            characterId: exactMatch.id,
                            characterName: exactMatch.name,
                            status: 'pending',
                            message: 'Character found, preparing to add...'
                        })
                    } else {
                        processResults.push({
                            input: name,
                            status: 'error',
                            message: `Multiple characters found (${searchResults.length}). Please use exact name.`
                        })
                    }
                } else {
                    // Single match found
                    const character = searchResults[0]
                    characterIdsToAdd.push(character.id)
                    processResults.push({
                        input: name,
                        characterId: character.id,
                        characterName: character.name,
                        status: 'pending',
                        message: 'Character found, preparing to add...'
                    })
                }
            } catch (error) {
                processResults.push({
                    input: name,
                    status: 'error',
                    message: error instanceof Error ? error.message : 'Failed to resolve character name'
                })
            }

            // Update results after each name resolution
            setResults([...processResults])
        }

        // Now add all resolved characters
        if (characterIdsToAdd.length > 0) {
            try {
                const response = await batchAddCharacters(characterIdsToAdd)
                
                // Update results with add status
                let successIndex = 0
                let failureIndex = 0
                
                for (let i = 0; i < processResults.length; i++) {
                    if (processResults[i].status === 'pending' && processResults[i].characterId) {
                        if (successIndex < response.success_count) {
                            processResults[i].status = 'success'
                            processResults[i].message = 'Character added successfully'
                            successIndex++
                        } else if (failureIndex < response.failure_count) {
                            const errorMessage = response.errors && response.errors[failureIndex]
                                ? response.errors[failureIndex]
                                : 'Failed to add character'
                            processResults[i].status = 'error'
                            processResults[i].message = errorMessage
                            failureIndex++
                        }
                    }
                    setProgress(50 + ((i + 1) / processResults.length) * 50) // Second half for adding
                }
            } catch (error) {
                // Mark all pending characters as failed
                for (let i = 0; i < processResults.length; i++) {
                    if (processResults[i].status === 'pending') {
                        processResults[i].status = 'error'
                        processResults[i].message = error instanceof Error ? error.message : 'Failed to add character'
                    }
                }
            }
        }

        return processResults
    }

    const handleBatchAdd = async () => {
        if (!characterInput.trim()) return

        const inputs = characterInput
            .split('\n')
            .map(input => input.trim())
            .filter(input => input.length > 0)

        if (inputs.length === 0) {
            setResults([{
                input: 'Invalid',
                status: 'error',
                message: 'No valid input found'
            }])
            return
        }

        setIsProcessing(true)
        setResults([])
        setProgress(0)

        let processResults: ProcessingResult[] = []

        if (inputType === 'ids') {
            // Validate that all inputs are numbers
            const ids = inputs
                .filter(input => !isNaN(Number(input)))
                .map(input => Number(input))

            if (ids.length === 0) {
                setResults([{
                    input: 'Invalid',
                    status: 'error',
                    message: 'No valid character IDs found'
                }])
                setIsProcessing(false)
                return
            }

            processResults = await handleBatchAddByIds(ids)
        } else {
            // Process names
            processResults = await handleBatchAddByNames(inputs)
        }

        setResults(processResults)
        setProgress(100)
        setIsProcessing(false)
        
        // Call the completion callback if provided
        onBatchComplete?.()
    }

    const clearResults = () => {
        setResults([])
        setCharacterInput('')
        setProgress(0)
    }

    const getStatusIcon = (status: ProcessingResult['status']) => {
        switch (status) {
            case 'success':
                return <Check className="size-4 text-green-500" />
            case 'error':
                return <X className="size-4 text-red-500" />
            case 'pending':
                return <Loader2 className="size-4 text-blue-500 animate-spin" />
            case 'resolving':
                return <Search className="size-4 text-yellow-500 animate-pulse" />
        }
    }

    const getStatusBadge = (status: ProcessingResult['status']) => {
        switch (status) {
            case 'success':
                return <Badge variant="default" className="bg-green-500">Success</Badge>
            case 'error':
                return <Badge variant="destructive">Error</Badge>
            case 'pending':
                return <Badge variant="secondary">Adding...</Badge>
            case 'resolving':
                return <Badge variant="outline">Resolving...</Badge>
        }
    }

    const successCount = results.filter(r => r.status === 'success').length
    const errorCount = results.filter(r => r.status === 'error').length

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold">Batch Character Management</h1>
                <p className="text-muted-foreground">
                    Add multiple characters by their EVE Online character names or IDs
                </p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Input Section */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Plus className="size-5" />
                            Add Characters
                        </CardTitle>
                        <CardDescription>
                            Choose whether to add characters by their names or IDs, then enter them one per line.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <Tabs value={inputType} onValueChange={(value) => setInputType(value as 'ids' | 'names')}>
                            <TabsList className="grid w-full grid-cols-2">
                                <TabsTrigger value="names" className="flex items-center gap-2">
                                    <User className="size-4" />
                                    Character Names
                                </TabsTrigger>
                                <TabsTrigger value="ids" className="flex items-center gap-2">
                                    <Hash className="size-4" />
                                    Character IDs
                                </TabsTrigger>
                            </TabsList>
                            
                            <TabsContent value="names" className="space-y-4">
                                <div className="space-y-2">
                                    <Label htmlFor="character-names">Character Names</Label>
                                    <Textarea
                                        id="character-names"
                                        placeholder="Enter character names, one per line:&#10;Tadeas CZ&#10;John Doe&#10;Jane Smith"
                                        value={characterInput}
                                        onChange={(e) => setCharacterInput(e.target.value)}
                                        rows={8}
                                        disabled={isProcessing}
                                    />
                                    <p className="text-xs text-muted-foreground">
                                        Enter exact character names. Names will be resolved to IDs automatically.
                                    </p>
                                </div>
                            </TabsContent>
                            
                            <TabsContent value="ids" className="space-y-4">
                                <div className="space-y-2">
                                    <Label htmlFor="character-ids">Character IDs</Label>
                                    <Textarea
                                        id="character-ids"
                                        placeholder="Enter character IDs, one per line:&#10;1772807647&#10;987654321&#10;555666777"
                                        value={characterInput}
                                        onChange={(e) => setCharacterInput(e.target.value)}
                                        rows={8}
                                        disabled={isProcessing}
                                    />
                                    <p className="text-xs text-muted-foreground">
                                        Paste character IDs from your spreadsheet or text file
                                    </p>
                                </div>
                            </TabsContent>
                        </Tabs>

                        <div className="flex gap-2">
                            <Button
                                onClick={handleBatchAdd}
                                disabled={!characterInput.trim() || isProcessing}
                                className="flex items-center gap-2"
                            >
                                {isProcessing ? (
                                    <Loader2 className="size-4 animate-spin" />
                                ) : (
                                    <Upload className="size-4" />
                                )}
                                {isProcessing ? 'Processing...' : `Process ${inputType === 'names' ? 'Names' : 'IDs'}`}
                            </Button>

                            {results.length > 0 && (
                                <Button
                                    variant="outline"
                                    onClick={clearResults}
                                    disabled={isProcessing}
                                >
                                    Clear Results
                                </Button>
                            )}
                        </div>

                        {isProcessing && (
                            <div className="space-y-2">
                                <div className="flex justify-between text-sm">
                                    <span>
                                        {inputType === 'names' 
                                            ? progress < 50 
                                                ? 'Resolving character names...' 
                                                : 'Adding characters...'
                                            : 'Processing characters...'
                                        }
                                    </span>
                                    <span>{Math.round(progress)}%</span>
                                </div>
                                <Progress value={progress} className="h-2" />
                            </div>
                        )}
                    </CardContent>
                </Card>

                {/* Results Section */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Users className="size-5" />
                            Processing Results
                        </CardTitle>
                        <CardDescription>
                            Status of character additions
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        {results.length === 0 ? (
                            <div className="text-center text-muted-foreground py-8">
                                <AlertCircle className="size-12 mx-auto mb-4 opacity-50" />
                                <p>No processing results yet</p>
                                <p className="text-sm">Enter character IDs and click &ldquo;Process Characters&rdquo; to begin</p>
                            </div>
                        ) : (
                            <div className="space-y-4">
                                {/* Summary */}
                                <div className="flex items-center justify-between p-4 bg-muted rounded-lg">
                                    <div className="text-sm">
                                        <span className="font-medium">Total: {results.length}</span>
                                        {successCount > 0 && (
                                            <span className="ml-4 text-green-600">✓ {successCount} successful</span>
                                        )}
                                        {errorCount > 0 && (
                                            <span className="ml-4 text-red-600">✗ {errorCount} failed</span>
                                        )}
                                    </div>
                                </div>

                                <Separator />

                                {/* Results List */}
                                <div className="space-y-2 max-h-96 overflow-y-auto">
                                    {results.map((result, index) => (
                                        <div
                                            key={index}
                                            className="flex items-center justify-between p-3 border rounded-lg"
                                        >
                                            <div className="flex items-center gap-3">
                                                {getStatusIcon(result.status)}
                                                <div className="min-w-0 flex-1">
                                                    <div className="flex items-center gap-2">
                                                        <p className="font-mono text-sm truncate">{result.input}</p>
                                                        {result.characterName && result.input !== result.characterName && (
                                                            <span className="text-xs text-muted-foreground">
                                                                → {result.characterName}
                                                            </span>
                                                        )}
                                                    </div>
                                                    {result.characterId && (
                                                        <p className="text-xs text-muted-foreground font-mono">
                                                            ID: {result.characterId}
                                                        </p>
                                                    )}
                                                    {result.message && (
                                                        <p className="text-xs text-muted-foreground">{result.message}</p>
                                                    )}
                                                </div>
                                            </div>
                                            {getStatusBadge(result.status)}
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>

            {/* Tips Section */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <AlertCircle className="size-5" />
                        Tips for Batch Adding Characters
                    </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                        <div>
                            <h4 className="font-medium mb-2 flex items-center gap-2">
                                <User className="size-4" />
                                Using Character Names:
                            </h4>
                            <ul className="space-y-1 text-muted-foreground">
                                <li>• Use exact character names</li>
                                <li>• Names are case-insensitive</li>
                                <li>• Automatically resolves to character IDs</li>
                                <li>• Shows character info after resolution</li>
                            </ul>
                        </div>
                        <div>
                            <h4 className="font-medium mb-2 flex items-center gap-2">
                                <Hash className="size-4" />
                                Using Character IDs:
                            </h4>
                            <ul className="space-y-1 text-muted-foreground">
                                <li>• Find IDs in EVE Online profiles</li>
                                <li>• Use third-party tools like zkillboard</li>
                                <li>• Check corporation member lists</li>
                                <li>• Faster processing than names</li>
                            </ul>
                        </div>
                        <div>
                            <h4 className="font-medium mb-2">Best Practices:</h4>
                            <ul className="space-y-1 text-muted-foreground">
                                <li>• Process in batches of 50 or fewer</li>
                                <li>• Verify names/IDs before processing</li>
                                <li>• Keep a backup of your character lists</li>
                                <li>• Review failed additions manually</li>
                            </ul>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
