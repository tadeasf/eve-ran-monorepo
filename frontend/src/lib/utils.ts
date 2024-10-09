

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
