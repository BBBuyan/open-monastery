export type UserRole = 'customer' | 'monk'

export type User = {
  id: number
  first_name: string
  last_name: string
  user_type: UserRole
  monastery_id?: number | null
}

// A catalog prayer or ritual; Book is the existing API name.
export type Book = {
  id: number
  name: string
  description: string
  price: number
}

export type Category = {
  id: number
  name: string
}

export type ReadingStatus = 'pending' | 'scheduled' | 'cancelled' | 'refunded'

// A request for monks to recite a prayer, rather than a text-reading session.
export type Reading = {
  id: number
  book_id: number
  monk_id: number | null
  customer_id: number
  status: ReadingStatus | 'unknown'
  created_at?: string
  assigned_at?: string
}

// Proposed list route: confirm with the server contract before enabling live history.
export async function getUserReadings(
  userId: number,
  options: { cursor?: string; signal?: AbortSignal } = {},
): Promise<Paginated<Reading>> {
  if (!Number.isSafeInteger(userId) || userId <= 0)
    throw new Error('Invalid user ID.')
  const params = new URLSearchParams({ user_id: String(userId), limit: '20' })
  if (options.cursor) params.set('cursor', options.cursor)
  const page = await request<Paginated<Reading>>(`/readings?${params}`, {
    signal: options.signal,
  })
  if (
    !Array.isArray(page?.data) ||
    typeof page.pagination?.hasMore !== 'boolean' ||
    !(
      page.pagination.nextCursor === null ||
      (typeof page.pagination.nextCursor === 'string' &&
        page.pagination.nextCursor.length > 0)
    ) ||
    (page.pagination.hasMore &&
      (!page.data.length ||
        page.pagination.nextCursor === null ||
        page.pagination.nextCursor === options.cursor))
  ) {
    throw new Error('The server returned an unexpected prayer request page.')
  }
  const statuses: string[] = ['pending', 'scheduled', 'cancelled', 'refunded']
  return {
    ...page,
    data: page.data.map((reading) => {
      if (
        !reading ||
        reading.customer_id !== userId ||
        !Number.isSafeInteger(reading.id) ||
        reading.id <= 0 ||
        !Number.isSafeInteger(reading.book_id) ||
        reading.book_id <= 0
      ) {
        throw new Error('The server returned an invalid prayer request.')
      }
      return {
        ...reading,
        status: statuses.includes(reading.status) ? reading.status : 'unknown',
      }
    }),
  }
}

const API_BASE_URL = (import.meta.env?.VITE_API_BASE_URL || '/api/v1').replace(
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

export type Paginated<T> = {
  data: T[]
  pagination: {
    nextCursor: string | null
    hasMore: boolean
  }
}

export type BookPage = Paginated<Book>

export type BookListOptions = {
  limit?: number
  cursor?: string
  query?: string
  signal?: AbortSignal
}

// Omit categoryId to search across the catalog. Both endpoints are paginated.
export async function getBooks(
  categoryId?: number,
  options: BookListOptions = {},
): Promise<BookPage> {
  const path =
    categoryId === undefined ? '/books' : `/categories/${categoryId}/books`
  const params = new URLSearchParams({ limit: String(options.limit ?? 20) })
  if (options.cursor !== undefined) params.set('cursor', options.cursor)
  if (options.query?.trim()) params.set('q', options.query.trim())
  const payload = await request<BookPage>(`${path}?${params}`, {
    signal: options.signal,
  })
  if (
    !payload ||
    !Array.isArray(payload.data) ||
    !payload.pagination ||
    typeof payload.pagination.hasMore !== 'boolean' ||
    !(
      payload.pagination.nextCursor === null ||
      (typeof payload.pagination.nextCursor === 'string' &&
        payload.pagination.nextCursor.length > 0)
    ) ||
    (payload.pagination.hasMore && payload.pagination.nextCursor === null)
  ) {
    throw new Error('The server returned an unexpected book page.')
  }
  const books = payload.data.map((value) => {
    const book = value as Partial<Book> | null
    if (!book || !Number.isSafeInteger(book.id) || Number(book.id) <= 0) {
      throw new Error('The server returned an invalid book ID.')
    }
    return {
      id: Number(book.id),
      name: String(book.name ?? 'Unnamed book'),
      description: String(book.description ?? ''),
      price: Number(book.price ?? 0),
    }
  })
  return { data: books, pagination: payload.pagination }
}

export async function getCategories(signal?: AbortSignal): Promise<Category[]> {
  const payload = await request<unknown>('/categories', { signal })
  return pickCollection(payload, 'categories').map((value) => {
    const category = value as Partial<Category> | null
    const id = Number(category?.id)
    if (
      !Number.isSafeInteger(id) ||
      id <= 0 ||
      typeof category?.name !== 'string'
    ) {
      throw new Error('The server returned an invalid category.')
    }
    return { id, name: category.name }
  })
}

export async function getCategoriesByBookId(
  bookId: number,
): Promise<Category[]> {
  const payload = await request<unknown>(`/books/${bookId}/categories`)
  return pickCollection(payload, 'categories') as Category[]
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
