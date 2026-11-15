import assert from 'node:assert/strict'
import { test } from 'node:test'
import { ApiError, logIn, registerUser } from '../src/api.ts'
import { sessionFromToken } from '../src/session.ts'

const user = {
  first_name: 'Ariun',
  last_name: 'Bat',
  email: 'ariun@example.com',
  password: 'example password',
}

function tokenWith(claims) {
  return `header.${Buffer.from(JSON.stringify(claims)).toString('base64url')}.signature`
}

test('email signup and login use the existing auth routes and payloads', async (t) => {
  const expires = Math.floor(Date.now() / 1000) + 900
  const token = tokenWith({ sub: '7', exp: expires })
  const fetch = t.mock.method(globalThis, 'fetch', async (url) =>
    url.endsWith('/users')
      ? Response.json({ id: 7 }, { status: 201 })
      : Response.json({ access_token: token, token_type: 'Bearer' }),
  )
  await registerUser(user)
  const accessToken = await logIn(user.email, user.password)
  const [registration, login] = fetch.mock.calls.map(
    ({ arguments: [path, options] }) => [
      path,
      options.method,
      JSON.parse(options.body),
    ],
  )
  assert.deepEqual(registration, ['/api/v1/auth/register/users', 'POST', user])
  assert.deepEqual(login, [
    '/api/v1/auth/login',
    'POST',
    { email: user.email, password: user.password },
  ])
  assert.deepEqual(sessionFromToken(accessToken, user.email), {
    accessToken: token,
    email: user.email,
    userId: 7,
    expiresAt: expires * 1000,
  })
})

test('invalid login and expired tokens cannot start a session', async (t) => {
  t.mock.method(globalThis, 'fetch', async () =>
    Response.json({ title: 'invalid.credentials' }, { status: 401 }),
  )
  await assert.rejects(
    logIn(user.email, user.password),
    (error) => error instanceof ApiError && error.status === 401,
  )
  assert.throws(
    () => sessionFromToken(tokenWith({ sub: '7', exp: 1 }), user.email),
    /expired session/,
  )
})
