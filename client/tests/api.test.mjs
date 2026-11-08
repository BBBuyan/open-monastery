import assert from 'node:assert/strict'
import { test } from 'node:test'
import { getBooks } from '../src/api.ts'

test('book endpoints preserve opaque cursors and use the new response envelope', async (t) => {
  const urls = []
  const page = {
    data: [{ id: 7, name: 'Book', description: '', price: 0 }],
    pagination: { hasMore: true, nextCursor: 'opaque/7+=' },
  }
  t.mock.method(globalThis, 'fetch', async (url) => {
    urls.push(new URL(url, 'http://localhost'))
    return Response.json(page)
  })
  const first = await getBooks()
  assert.deepEqual(first, page)
  await getBooks(3, { cursor: first.pagination.nextCursor, query: ' search ' })
  assert.equal(urls[0].pathname, '/api/v1/books')
  assert.equal(urls[0].searchParams.has('cursor'), false)
  assert.equal(urls[1].pathname, '/api/v1/categories/3/books')
  assert.equal(urls[1].searchParams.get('cursor'), 'opaque/7+=')
  assert.equal(urls[1].searchParams.has('after_id'), false)
  assert.equal(urls[1].searchParams.get('q'), 'search')
})

test('empty final pages are accepted and invalid continuation metadata is rejected', async (t) => {
  let page = { data: [], pagination: { hasMore: false, nextCursor: null } }
  t.mock.method(globalThis, 'fetch', async () => Response.json(page))
  assert.deepEqual(await getBooks(), page)
  for (const pagination of [
    null,
    { hasMore: true, nextCursor: null },
    { hasMore: true, nextCursor: 7 },
    { hasMore: true, nextCursor: '' },
  ]) {
    page = { data: [], pagination }
    await assert.rejects(getBooks(), /unexpected book page/)
  }
})
