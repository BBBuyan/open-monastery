export type UserRole = 'customer' | 'monk'

export type User = {
  id: number
  first_name: string
  last_name: string
  user_type: UserRole
  monastery_id?: number | null
}

export type Book = {
  id: number
  name: string
  description: string
  price: number
}

export type ReadingStatus = 'pending' | 'scheduled' | 'cancelled' | 'refunded'

export type Reading = {
  id: number
  book_id: number
  monk_id: number | null
  customer_id: number
  status: ReadingStatus
  created_at?: string
}

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api/v1').replace(
  /\/$/,
  '',
)

function pickCollection(payload: unknown, key: string): unknown[] {
  if (Array.isArray(payload)) return payload
  const wrapped = payload as Record<string, unknown> | null
  const candidate = wrapped?.[key] ?? wrapped?.data ?? wrapped?.items
  if (!Array.isArray(candidate))
    throw new Error('The server returned an unexpected response.')
  return candidate
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      Accept: 'application/json',
      ...(options?.body ? { 'Content-Type': 'application/json' } : {}),
      ...options?.headers,
    },
  })

  if (!response.ok) {
    let message = `The server responded with ${response.status}.`
    try {
      const body = (await response.json()) as {
        error?: string
        message?: string
      }
      message = body.message || body.error || message
    } catch {
      // Keep the status-based fallback for non-JSON errors.
    }
    throw new Error(message)
  }

  return response.json() as Promise<T>
}

export async function getBooks(): Promise<Book[]> {
  const payload = await request<unknown>('/books')
  return pickCollection(payload, 'books').map((value, index) => {
    const book = value as Partial<Book>
    return {
      id: Number(book.id ?? index),
      name: String(book.name ?? 'Untitled sutra'),
      description: String(book.description ?? ''),
      price: Number(book.price ?? 0),
    }
  })
}

export function getUser(id: number): Promise<User> {
  return request<User>(`/users/${id}`)
}

export function createCustomer(
  firstName: string,
  lastName: string,
): Promise<User> {
  return request<User>('/users', {
    method: 'POST',
    body: JSON.stringify({
      first_name: firstName,
      last_name: lastName,
      user_type: 'customer',
    }),
  })
}

export async function getReadings(): Promise<Reading[]> {
  const payload = await request<unknown>('/readings')
  return pickCollection(payload, 'readings').map((value) => {
    const reading = value as Partial<Reading> & { CreatedAt?: string }
    return {
      id: Number(reading.id),
      book_id: Number(reading.book_id),
      monk_id: reading.monk_id == null ? null : Number(reading.monk_id),
      customer_id: Number(reading.customer_id),
      status: reading.status || 'pending',
      created_at: reading.created_at || reading.CreatedAt,
    }
  })
}

export function createReading(
  bookId: number,
  customerId: number,
): Promise<Reading> {
  return request<Reading>('/readings', {
    method: 'POST',
    body: JSON.stringify({ book_id: bookId, customer_id: customerId }),
  })
}

export function claimReading(
  readingId: number,
  monkId: number,
): Promise<Reading> {
  return request<Reading>(`/readings/${readingId}`, {
    method: 'PATCH',
    body: JSON.stringify({ monk_id: monkId, status: 'scheduled' }),
  })
}
