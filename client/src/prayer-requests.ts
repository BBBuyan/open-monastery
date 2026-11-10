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

// Fictional records only. Never merge these with server responses or checkout.
export const demoRequests: PrayerRequest[] = [
  {
    id: 124,
    book_id: 158,
    book_name: 'Манал',
    customer_id: 1,
    monk_id: null,
    status: 'pending',
    created_at: '2026-11-09T08:00:00+08:00',
  },
  {
    id: 118,
    book_id: 184,
    book_name: 'Ногоон дарь эх',
    customer_id: 1,
    monk_id: 7,
    status: 'scheduled',
    created_at: '2026-11-07T08:00:00+08:00',
    assigned_at: '2026-11-08T08:00:00+08:00',
  },
  {
    id: 105,
    book_id: 306,
    book_name: 'Алтангэрэл',
    customer_id: 1,
    monk_id: null,
    status: 'cancelled',
    created_at: '2026-11-03T08:00:00+08:00',
  },
  {
    id: 99,
    book_id: 158,
    book_name: 'Манал',
    customer_id: 1,
    monk_id: null,
    status: 'refunded',
    created_at: '2026-11-01T08:00:00+08:00',
  },
]
