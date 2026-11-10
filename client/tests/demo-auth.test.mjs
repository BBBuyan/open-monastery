import assert from 'node:assert/strict'
import { test } from 'node:test'
import { createDemoAuth, loginMethods } from '../src/demo-auth.ts'

test('all four methods can enter and end an explicit demo session', () => {
  for (const method of loginMethods) {
    const auth = createDemoAuth()
    assert.equal(auth.loggedIn.value, false)
    assert.equal(auth.enter(method), true)
    assert.equal(auth.method.value, method)
    assert.equal(auth.loggedIn.value, true)
    auth.signOut()
    assert.equal(auth.loggedIn.value, false)
    assert.equal(auth.method.value, null)
  }
})

test('demo sessions are not restored on reload and unsupported methods are rejected', () => {
  const first = createDemoAuth()
  first.enter('google')
  assert.equal(createDemoAuth().loggedIn.value, false)
  const other = createDemoAuth()
  assert.equal(other.enter('unsupported'), false)
  assert.equal(other.loggedIn.value, false)
})
