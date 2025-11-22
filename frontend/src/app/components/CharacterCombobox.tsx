'use client'

import * as React from "react"
import { Check, ChevronsUpDown } from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/app/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { Character } from "@/lib/types"

interface CharacterComboboxProps {
  characters: Character[]
  value: number | null
  onValueChange: (value: number | null) => void
  placeholder?: string
  className?: string
}

export function CharacterCombobox({
  characters,
  value,
  onValueChange,
  placeholder = "Select a character...",
  className,
}: CharacterComboboxProps) {
  const [open, setOpen] = React.useState(false)

  // Sort characters alphabetically by name
  const sortedCharacters = React.useMemo(() => {
    return [...characters].sort((a, b) =>
      a.name.localeCompare(b.name, undefined, { sensitivity: 'base' })
    )
  }, [characters])

  const selectedCharacter = sortedCharacters.find(char => char.id === value)

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn("w-full sm:w-[250px] justify-between", className)}
        >
          {selectedCharacter ? selectedCharacter.name : placeholder}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full sm:w-[250px] p-0" align="start">
        <Command>
          <CommandInput placeholder="Search characters..." className="h-9" />
          <CommandEmpty>No character found.</CommandEmpty>
          <CommandList>
            <CommandGroup>
              {sortedCharacters.map((character) => (
                <CommandItem
                  key={character.id}
                  value={character.name}
                  onSelect={() => {
                    onValueChange(character.id === value ? null : character.id)
                    setOpen(false)
                  }}
                >
                  {character.name}
                  <Check
                    className={cn(
                      "ml-auto h-4 w-4",
                      value === character.id ? "opacity-100" : "opacity-0"
                    )}
                  />
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
