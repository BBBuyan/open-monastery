import type { Reading } from './api'

export const requestStatuses = {
  pending: { label: 'Pending', description: 'Awaiting assignment to a monk.' },
  scheduled: {
    label: 'Assigned',
    description: 'A monk has been assigned to your request.',
  },
  cancelled: {
    label: 'Cancelled',
    description: 'This prayer request has been cancelled.',
  },
  refunded: {
    label: 'Refunded',
    description: 'This prayer request has been refunded.',
  },
  unknown: {
    label: 'Status unavailable',
    description:
      'The current status is not available. Please check again later.',
  },
}

export function requestStatus(status: string) {
  return (
    requestStatuses[status as keyof typeof requestStatuses] ??
    requestStatuses.unknown
  )
}

export function requestDate(value?: string) {
  if (!value || value.startsWith('0001-')) return null
  const date = new Date(value)
  return Number.isNaN(date.getTime())
    ? null
    : new Intl.DateTimeFormat('en-GB', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      }).format(date)
}

export type PrayerRequest = Reading & { book_name: string; monk_name?: string }
