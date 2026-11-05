import { ref, shallowRef } from 'vue'

type Page<T> = { books: T[]; next_cursor: number | null; has_more: boolean }
type PageRequest = {
  afterId?: number
  query: string
  categoryId?: number
  signal: AbortSignal
}

/** Keeps cursor, retry, and search cancellation consistent for every catalog load. */
export function createCatalog<T extends { id: number }>(
  fetchPage: (options: PageRequest) => Promise<Page<T>>,
) {
  const books = shallowRef<T[]>([])
  const loading = ref(false)
  const error = ref(false)
  const hasMore = ref(true)
  let cursor: number | undefined
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
        afterId: cursor,
        query,
        categoryId,
        signal: request.signal,
      })
      if (current !== generation || request.signal.aborted) return
      if (
        page.has_more &&
        (!page.books.length ||
          !Number.isSafeInteger(page.next_cursor) ||
          page.next_cursor! <= (cursor ?? 0))
      ) {
        throw new Error('The catalog cursor did not advance.')
      }
      const existing = new Set(books.value.map((book) => book.id))
      books.value = [
        ...books.value,
        ...page.books.filter((book) => {
          if (existing.has(book.id)) return false
          existing.add(book.id)
          return true
        }),
      ]
      cursor = page.next_cursor ?? undefined
      hasMore.value = page.has_more
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
