import assert from 'node:assert/strict'
import { test } from 'node:test'
import {
  getMonasteries,
  getMonastery,
  getCategories,
  getBooks,
} from '../src/api.ts'
import {
  monasteryCatalogLink,
  monasteryIdFromHash,
} from '../src/monasteries.ts'

test('monastery list maps Go fields and preserves the next page cursor', async (t) => {
  const urls = []
  t.mock.method(globalThis, 'fetch', async (url) => {
    urls.push(new URL(url, 'http://localhost'))
    return Response.json(
      urls.length === 1
        ? {
            data: [{ Id: 42, Name: 'Monastery' }],
            pagination: { hasMore: true, nextCursor: 'opaque/42+=' },
          }
        : { data: [], pagination: { hasMore: false, nextCursor: null } },
    )
  })
  const page = await getMonasteries()
  assert.deepEqual(page.data, [{ id: 42, name: 'Monastery' }])
  await getMonasteries({ cursor: page.pagination.nextCursor })
  assert.equal(urls[0].pathname, '/api/v1/monasteries')
  assert.equal(urls[1].searchParams.get('cursor'), 'opaque/42+=')
})

test('selected monastery links survive reload and load scoped books and categories', async (t) => {
  const urls = []
  t.mock.method(globalThis, 'fetch', async (url) => {
    urls.push(url)
    if (url.endsWith('/categories'))
      return Response.json([{ Id: 9, Name: 'Protection', MonasteryId: 42 }])
    return Response.json({ Id: 42, Name: 'Monastery' })
  })
  const id = monasteryIdFromHash(monasteryCatalogLink(42))
  assert.equal(id, 42)
  assert.deepEqual(await getMonastery(id), { id: 42, name: 'Monastery' })
  assert.deepEqual(await getCategories(id), [{ id: 9, name: 'Protection' }])
  assert.deepEqual(urls, [
    '/api/v1/monasteries/42',
    '/api/v1/monasteries/42/categories',
  ])
  for (const hash of [
    '#readings?monastery=amber-valley',
    '#readings?monastery=-1',
    '#readings?monastery=0',
    '#home?monastery=42',
  ]) {
    assert.equal(monasteryIdFromHash(hash), undefined)
  }
})

test('catalog rejects data from a different monastery', async (t) => {
  t.mock.method(globalThis, 'fetch', async (url) =>
    url.endsWith('/categories')
      ? Response.json([{ Id: 9, Name: 'Protection', MonasteryId: 43 }])
      : Response.json({
          data: [
            {
              Id: 1,
              Name: 'Book',
              Description: '',
              Price: 10,
              MonasteryId: 43,
            },
          ],
          pagination: { hasMore: false, nextCursor: null },
        }),
  )
  await assert.rejects(getCategories(42), /invalid category/)
  await assert.rejects(getBooks(42), /invalid book/)
})

test('empty categories are displayed as an empty list', async (t) => {
  t.mock.method(globalThis, 'fetch', async () => Response.json(null))
  assert.deepEqual(await getCategories(42), [])
})
