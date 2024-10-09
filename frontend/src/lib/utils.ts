// Copyright (C) 2024 Tadeáš Fořt
// 
// This file is part of EVE Ran Services.
// 
// EVE Ran Services is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// EVE Ran Services is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with EVE Ran Services.  If not, see <https://www.gnu.org/licenses/>.

import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatISK(value: number): string {
  const absValue = Math.abs(value);
  if (absValue >= 1000000000000) {
    return (value / 1000000000000).toFixed(2) + ' tril ISK';
  } else if (absValue >= 1000000000) {
    return (value / 1000000000).toFixed(2) + ' bil ISK';
  } else if (absValue >= 1000000) {
    return (value / 1000000).toFixed(2) + ' mil ISK';
  } else if (absValue >= 1000) {
    return (value / 1000).toFixed(2) + 'k ISK';
  } else {
    return value.toFixed(2) + ' ISK';
  }
}

export function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString()
}
