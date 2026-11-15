import assert from 'node:assert/strict'
import { test } from 'node:test'
import { getUserReadings } from '../src/api.ts'
import { requestStatus, requestDate } from '../src/prayer-requests.ts'

test('request status does not imply a scheduled date or silently invent pending', () => {
  assert.equal(requestStatus('scheduled').label, 'Assigned')
  assert.equal(requestStatus('').label, 'Status unavailable')
  assert.equal(requestStatus('completed').label, 'Status unavailable')
  assert.equal(requestDate('0001-01-01T00:00:00Z'), null)
  assert.equal(requestDate('invalid'), null)
  assert.equal(requestDate(), null)
})

test('user request pages use bearer auth, preserve cursors, and reject other users', async (t) => {
  let reading = {
    id: 5,
    book_id: 158,
    customer_id: 7,
    monk_id: null,
    status: '',
  }
  let url
  let authorization
  t.mock.method(globalThis, 'fetch', async (value, options) => {
    url = new URL(value, 'http://localhost')
    authorization = options.headers.Authorization
    return Response.json({
      data: [reading],
      pagination: { hasMore: false, nextCursor: null },
    })
  })
  const result = await getUserReadings(7, 'token', { cursor: 'opaque/+=' })
  assert.equal(authorization, 'Bearer token')
  assert.equal(url.searchParams.has('user_id'), false)
  assert.equal(url.searchParams.get('cursor'), 'opaque/+=')
  assert.equal(result.data[0].status, 'unknown')
  reading = { ...reading, customer_id: 8 }
  await assert.rejects(getUserReadings(7, 'token'), /invalid prayer request/)
  await assert.rejects(getUserReadings(0, 'token'), /Invalid user ID/)
})
