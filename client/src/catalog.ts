import type { Paginated } from './api'
import { ref, shallowRef } from 'vue'

type PageRequest = {
  cursor?: string
  query: string
  categoryId?: number
  signal: AbortSignal
}

/** Keeps cursor, retry, and search cancellation consistent for every catalog load. */
export function createCatalog<T extends { id: number }>(
  fetchPage: (options: PageRequest) => Promise<Paginated<T>>,
) {
  const books = shallowRef<T[]>([])
  const loading = ref(false)
  const error = ref(false)
  const hasMore = ref(true)
  let cursor: string | undefined
  let query = ''
  let categoryId: number | undefined
  let generation = 0
  let controller: AbortController | undefined

  async function loadMore() {
    if (loading.value || !hasMore.value) return
    const current = generation
    const request = new AbortController()
    controller = request
    loading.value = true
    error.value = false
    try {
      const page = await fetchPage({
        cursor,
        query,
        categoryId,
        signal: request.signal,
      })
      if (current !== generation || request.signal.aborted) return
      if (
        page.pagination.hasMore &&
        (!page.data.length ||
          typeof page.pagination.nextCursor !== 'string' ||
          !page.pagination.nextCursor.length ||
          page.pagination.nextCursor === cursor)
      ) {
        throw new Error('The catalog cursor did not advance.')
      }
      const existing = new Set(books.value.map((book) => book.id))
      books.value = [
        ...books.value,
        ...page.data.filter((book) => {
          if (existing.has(book.id)) return false
          existing.add(book.id)
          return true
        }),
      ]
      cursor = page.pagination.nextCursor ?? undefined
      hasMore.value = page.pagination.hasMore
    } catch {
      if (current === generation && !request.signal.aborted) error.value = true
    } finally {
      if (current === generation) loading.value = false
    }
  }

  function dispose() {
    generation++
    controller?.abort()
    loading.value = false
  }

  function reset(search = '', autoLoad = true, category?: number) {
    dispose()
    query = search.trim()
    categoryId = category
    cursor = undefined
    books.value = []
    error.value = false
    hasMore.value = true
    return autoLoad ? loadMore() : Promise.resolve()
  }

  return { books, loading, error, hasMore, loadMore, reset, dispose }
}
