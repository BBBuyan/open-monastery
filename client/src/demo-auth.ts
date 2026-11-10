import { computed, ref } from 'vue'

export const loginMethods = ['email', 'phone', 'google', 'facebook'] as const
export type LoginMethod = (typeof loginMethods)[number]

/** UI preview only: never authenticates, sends credentials, or persists a session. */
export function createDemoAuth() {
  const method = ref<LoginMethod | null>(null)
  function enter(value: LoginMethod) {
    if (!loginMethods.includes(value)) return false
    method.value = value
    return true
  }
  function signOut() {
    method.value = null
  }
  return {
    loggedIn: computed(() => method.value !== null),
    method: computed(() => method.value),
    enter,
    signOut,
  }
}
