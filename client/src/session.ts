export type Session = {
  accessToken: string
  email: string
  userId: number
  expiresAt: number
}

const storageKey = 'open-monastery-session'

export function sessionFromToken(accessToken: string, email: string): Session {
  // The API verifies the signed token. The browser reads these claims for routing and expiry.
  const [, payload] = accessToken.split('.')
  if (!payload) throw new Error('The server returned an invalid session.')
  const claims = JSON.parse(
    atob(payload.replace(/-/g, '+').replace(/_/g, '/')),
  ) as { sub?: string; exp?: number }
  const userId = Number(claims.sub)
  const expiresAt = Number(claims.exp) * 1000
  if (
    !Number.isSafeInteger(userId) ||
    userId <= 0 ||
    !Number.isFinite(expiresAt) ||
    expiresAt <= Date.now()
  ) {
    throw new Error('The server returned an invalid or expired session.')
  }
  return { accessToken, email, userId, expiresAt }
}

export function readSession(): Session | null {
  try {
    const saved = sessionStorage.getItem(storageKey)
    if (!saved) return null
    const { accessToken, email } = JSON.parse(saved)
    if (typeof accessToken === 'string' && typeof email === 'string') {
      return sessionFromToken(accessToken, email)
    }
  } catch {
    // An expired token or unavailable storage starts a signed-out visit.
  }
  saveSession(null)
  return null
}

export function saveSession(session: Session | null) {
  try {
    if (session) {
      const { accessToken, email } = session
      sessionStorage.setItem(storageKey, JSON.stringify({ accessToken, email }))
    } else {
      sessionStorage.removeItem(storageKey)
    }
  } catch {
    // The active tab can still use the in-memory session.
  }
}
