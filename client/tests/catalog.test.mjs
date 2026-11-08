import assert from 'node:assert/strict'
import { test } from 'node:test'
import { createCatalog } from '../src/catalog.ts'

test('category pagination retains the category and search on subsequent pages', async () => {
  const calls = []
  const catalog = createCatalog(async (options) => {
    calls.push(options)
    return calls.length === 1 ? page([5, 12], 12) : page([22])
  })
  await catalog.reset('Манал', true, 3)
  await catalog.loadMore()
  assert.deepEqual(
    calls.map(({ categoryId, query, cursor }) => ({
      categoryId,
      query,
      cursor,
    })),
    [
      { categoryId: 3, query: 'Манал', cursor: undefined },
      { categoryId: 3, query: 'Манал', cursor: 'opaque/12+=' },
    ],
  )
  assert.deepEqual(
    catalog.books.value.map((item) => item.id),
    [5, 12, 22],
  )
})

test('switching category ignores a late old page and starts with no cursor', async () => {
  const oldPage = deferred()
  const calls = []
  const catalog = createCatalog((options) => {
    calls.push(options)
    if (calls.length === 1) return Promise.resolve(page([2], 2))
    if (calls.length === 2) return oldPage.promise
    return Promise.resolve(page([8]))
  })
  await catalog.reset('search', true, 1)
  const oldLoad = catalog.loadMore()
  await catalog.reset('search', true, 2)
  assert.equal(calls[1].signal.aborted, true)
  oldPage.resolve(page([4]))
  await oldLoad
  assert.equal(calls[2].categoryId, 2)
  assert.equal(calls[2].cursor, undefined)
  assert.deepEqual(
    catalog.books.value.map((item) => item.id),
    [8],
  )
})

test('clearing a category returns to the full catalog and resets its cursor', async () => {
  const calls = []
  const catalog = createCatalog(async (options) => {
    calls.push(options)
    return page([1], 1)
  })
  await catalog.reset('', true, 7)
  await catalog.reset()
  assert.equal(calls[1].categoryId, undefined)
  assert.equal(calls[1].cursor, undefined)
})

const book = (id) => ({ id, name: `Book ${id}` })
const page = (ids, next = null) => ({
  data: ids.map(book),
  pagination: {
    nextCursor: next === null ? null : `opaque/${next}+=`,
    hasMore: next !== null,
  },
})
function deferred() {
  let resolve
  const promise = new Promise((done) => {
    resolve = done
  })
  return { promise, resolve }
}

test('appends pages using the returned cursor and stops at the final page', async () => {
  const calls = []
  const catalog = createCatalog(async (options) => {
    calls.push(options)
    return calls.length === 1 ? page([2, 7], 7) : page([10])
  })
  await catalog.reset('  Манал  ')
  await catalog.loadMore()
  await catalog.loadMore()
  assert.deepEqual(
    catalog.books.value.map((item) => item.id),
    [2, 7, 10],
  )
  assert.equal(calls.length, 2)
  assert.equal(calls[0].cursor, undefined)
  assert.equal(calls[1].cursor, 'opaque/7+=')
  assert.equal(calls[1].query, 'Манал')
  assert.equal(catalog.hasMore.value, false)
})

test('repeated scroll triggers do not start concurrent page requests', async () => {
  const pending = deferred()
  let calls = 0
  const catalog = createCatalog(() => {
    calls++
    return pending.promise
  })
  const first = catalog.reset()
  await catalog.loadMore()
  await catalog.loadMore()
  assert.equal(calls, 1)
  pending.resolve(page([1]))
  await first
  assert.equal(catalog.loading.value, false)
})

test('a failed next page keeps earlier books and retries the same cursor', async () => {
  const cursors = []
  const catalog = createCatalog(async ({ cursor }) => {
    cursors.push(cursor)
    if (cursors.length === 1) return page([1, 3], 3)
    if (cursors.length === 2) throw new Error('Network unavailable')
    return page([3, 9])
  })
  await catalog.reset()
  await catalog.loadMore()
  assert.equal(catalog.error.value, true)
  assert.deepEqual(
    catalog.books.value.map((item) => item.id),
    [1, 3],
  )
  await catalog.loadMore()
  assert.deepEqual(cursors, [undefined, 'opaque/3+=', 'opaque/3+='])
  assert.deepEqual(
    catalog.books.value.map((item) => item.id),
    [1, 3, 9],
  )
  assert.equal(catalog.error.value, false)
})

test('search reset aborts and ignores stale responses even if fetch ignores abort', async () => {
  const old = deferred()
  const current = deferred()
  const calls = []
  const catalog = createCatalog((options) => {
    calls.push(options)
    return calls.length === 1 ? old.promise : current.promise
  })
  const oldLoad = catalog.reset('old')
  const currentLoad = catalog.reset('new')
  assert.equal(calls[0].signal.aborted, true)
  old.resolve(page([1], 1))
  await oldLoad
  assert.deepEqual(catalog.books.value, [])
  assert.equal(catalog.loading.value, true)
  current.resolve(page([20]))
  await currentLoad
  assert.deepEqual(
    catalog.books.value.map((item) => item.id),
    [20],
  )
  assert.equal(calls[1].query, 'new')
  assert.equal(calls[1].cursor, undefined)
})

test('disposing aborts the request and prevents late updates', async () => {
  const pending = deferred()
  let signal
  const catalog = createCatalog((options) => {
    signal = options.signal
    return pending.promise
  })
  const load = catalog.reset()
  catalog.dispose()
  pending.resolve(page([1]))
  await load
  assert.equal(signal.aborted, true)
  assert.deepEqual(catalog.books.value, [])
  assert.equal(catalog.loading.value, false)
})

test('a non-advancing cursor becomes a retryable error instead of a fetch loop', async () => {
  let calls = 0
  const catalog = createCatalog(async () =>
    ++calls === 1 ? page([1], 1) : page([1], 1),
  )
  await catalog.reset()
  await catalog.loadMore()
  assert.equal(catalog.error.value, true)
  assert.deepEqual(
    catalog.books.value.map((item) => item.id),
    [1],
  )
})

test('a search can invalidate the old page before its debounced load begins', async () => {
  const pending = deferred()
  const catalog = createCatalog(() => pending.promise)
  const oldLoad = catalog.reset()
  await catalog.reset('new search', false)
  pending.resolve(page([1]))
  await oldLoad
  assert.deepEqual(catalog.books.value, [])
  assert.equal(catalog.loading.value, false)
  assert.equal(catalog.hasMore.value, true)
})
