

import React from 'react'
import { Region } from '@/lib/types'
import { Button } from "@/app/components/ui/button"
import { Skeleton } from "@/app/components/ui/skeleton"
import { Progress } from "@/app/components/ui/progress"
import { Check, ChevronsUpDown } from "lucide-react"
import { cn } from "@/lib/utils"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "@/app/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/app/components/ui/popover"

interface FilterControlsProps {
  regions: Region[]
  selectedRegions: Array<{ id: number, name: string }>
  setSelectedRegions: React.Dispatch<React.SetStateAction<Array<{ id: number, name: string }>>>
  startDate: string
  setStartDate: React.Dispatch<React.SetStateAction<string>>
  endDate: string
  setEndDate: React.Dispatch<React.SetStateAction<string>>
  onApplyFilters: () => void
  isLoading: boolean
}

export default function FilterControls({
  regions,
  selectedRegions,
  setSelectedRegions,
  startDate,
  setStartDate,
  endDate,
  setEndDate,
  onApplyFilters,
  isLoading
}: FilterControlsProps) {
  const [open, setOpen] = React.useState(false)

  const handleSelectRegion = (regionId: number, regionName: string) => {
    setSelectedRegions(prev => {
      const index = prev.findIndex(r => r.id === regionId)
      if (index > -1) {
        return prev.filter((_, i) => i !== index)
      } else {
        return [...prev, { id: regionId, name: regionName }]
      }
    })
  }

  return (
    <div className="flex flex-col space-y-4 mb-4">
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div>
          <label className="block text-sm font-medium mb-1">Regions</label>
          {isLoading || regions.length === 0 ? (
            <Skeleton className="w-full h-10" />
          ) : (
            <Popover open={open} onOpenChange={setOpen}>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  role="combobox"
                  aria-expanded={open}
                  className="w-full justify-between"
                >
                  {selectedRegions.length > 0
                    ? `${selectedRegions.length} selected`
                    : "Select regions..."}
                  <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-full p-0">
                <Command>
                  <CommandInput placeholder="Search regions..." />
                  <CommandEmpty>No region found.</CommandEmpty>
                  <CommandGroup className="max-h-[300px] overflow-y-auto">
                    {Array.isArray(regions) && regions.length > 0 ? (
                      regions.map((region) => (
                        <CommandItem
                          key={region.region_id}
                          onSelect={() => handleSelectRegion(region.region_id, region.name)}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              selectedRegions.some(r => r.id === region.region_id) ? "opacity-100" : "opacity-0"
                            )}
                          />
                          {region.name}
                        </CommandItem>
                      ))
                    ) : (
                      <CommandItem>No regions available</CommandItem>
                    )}
                  </CommandGroup>
                </Command>
              </PopoverContent>
            </Popover>
          )}
        </div>
        <div>
          <label htmlFor="start-date" className="block text-sm font-medium mb-1">Start Date</label>
          {isLoading ? (
            <Skeleton className="w-full h-10" />
          ) : (
            <input
              id="start-date"
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              className="w-full rounded-[0.5rem] p-2 bg-background text-foreground border border-input hover:border-ring focus:border-ring focus:ring-1 focus:ring-ring focus:outline-none"
            />
          )}
        </div>
        <div>
          <label htmlFor="end-date" className="block text-sm font-medium mb-1">End Date</label>
          {isLoading ? (
            <Skeleton className="w-full h-10" />
          ) : (
            <input
              id="end-date"
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              className="w-full rounded-[0.5rem] p-2 bg-background text-foreground border border-input hover:border-ring focus:border-ring focus:ring-1 focus:ring-ring focus:outline-none"
            />
          )}
        </div>
      </div>
      <div className="flex justify-end">
        <Button onClick={onApplyFilters} disabled={isLoading} className="w-full sm:w-auto">
          {isLoading ? <Progress value={33} className="w-16" /> : "Apply Filters"}
        </Button>
      </div>
    </div>
  )
}