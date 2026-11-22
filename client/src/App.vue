<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import Monasteries from './components/Monasteries.vue'
import Basket from './components/Basket.vue'
import UserProfile from './components/UserProfile.vue'
import PrayerRequests from './components/PrayerRequests.vue'
import ProfileMenu from './components/ProfileMenu.vue'
import LoginPage from './components/LoginPage.vue'
import { readSession, saveSession, type Session } from './session'
import { basket } from './basket'
import { formatPrice } from './money'
import { useTheme } from './theme'
import { monasteryIdFromHash } from './monasteries'
const { theme, toggle: toggleTheme } = useTheme()
const session = ref(readSession())
const registeredEmail = ref('')
function finishRegistration(email: string) {
  registeredEmail.value = email
  location.hash = '#login'
}
function signIn(value: Session) {
  session.value = value
  saveSession(value)
  registeredEmail.value = ''
  location.hash = '#readings'
}
function clearSession() {
  session.value = null
  saveSession(null)
}
function signOut() {
  clearSession()
  location.hash = '#login'
}
function expireSession() {
  clearSession()
  if (active.value === 'profile' || active.value === 'requests') {
    location.hash = '#login'
  }
}

const { count: basketCount, total: basketTotal, announcement } = basket
function pageFromHash() {
  if (location.hash === '#login') return 'login'
  if (location.hash === '#signup') return 'signup'
  if (/^#(?:requests(?:\/\d+)?|profile\/\d+)$/.test(location.hash))
    return 'requests'
  if (location.hash === '#profile') return 'profile'
  return location.hash === '#basket'
    ? 'basket'
    : location.hash === '#checkout'
      ? 'checkout'
      : 'home'
}
const active = ref(pageFromHash())
const monasteryId = ref(monasteryIdFromHash(location.hash))
const requestId = ref<number>()
watch(
  session,
  (value, _, onCleanup) => {
    if (!value) return
    const duration = Math.max(0, value.expiresAt - Date.now())
    const timer = window.setTimeout(expireSession, duration)
    onCleanup(() => window.clearTimeout(timer))
  },
  { immediate: true },
)

async function sync() {
  if (location.hash === '#main-content') return
  const previous = active.value
  active.value = pageFromHash()
  monasteryId.value = monasteryIdFromHash(location.hash)
  if (
    !session.value &&
    (active.value === 'profile' || active.value === 'requests')
  ) {
    active.value = 'login'
    history.replaceState(null, '', '#login')
  }
  requestId.value = /^#(?:profile|requests)\/\d+$/.test(location.hash)
    ? Number(location.hash.split('/')[1])
    : undefined
  document.title =
    active.value === 'home'
      ? 'Open Monastery'
      : `${active.value === 'login' ? 'Log in' : active.value === 'signup' ? 'Sign up' : active.value === 'requests' ? (requestId.value ? `Prayer request #${requestId.value}` : 'Readings') : active.value === 'profile' ? 'Profile' : active.value === 'basket' ? 'Basket' : 'Checkout'} · Open Monastery`
  await nextTick()
  if (location.hash.startsWith('#readings')) {
    const catalog = document.getElementById('readings')
    catalog?.scrollIntoView({ block: 'start' })
    catalog?.focus({ preventScroll: true })
  } else if (previous !== active.value || location.hash === '#home') {
    window.scrollTo({ top: 0, behavior: 'instant' })
    document.getElementById('main-content')?.focus({ preventScroll: true })
  }
}
function focusMain() {
  document.getElementById('main-content')?.focus()
}
onMounted(() => {
  void sync()
  window.addEventListener('hashchange', sync)
})
onUnmounted(() => window.removeEventListener('hashchange', sync))
</script>

<template>
  <a class="skip-link" href="#main-content" @click.prevent="focusMain">
    Skip to content
  </a>
  <div class="app-shell">
    <p
      class="visually-hidden"
      role="status"
      aria-live="polite"
      aria-atomic="true"
    >
      {{ announcement }}
    </p>
    <header class="site-header">
      <div class="site-header-inner">
        <a class="site-brand" href="#home" aria-label="Open Monastery home">
          <img
            src="/open-monastery-icon-yellow.svg"
            width="44"
            height="44"
            alt=""
          />
          <span>Open Monastery</span>
        </a>
        <div class="header-actions">
          <a
            class="basket-link"
            href="#basket"
            :aria-label="`Basket, ${basketCount} ${basketCount === 1 ? 'book' : 'books'}, ${formatPrice(basketTotal)}`"
            :aria-current="active === 'basket' ? 'page' : undefined"
          >
            <span class="basket-link-total">{{
              formatPrice(basketTotal)
            }}</span>
            <svg
              viewBox="0 0 24 24"
              width="20"
              height="20"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
              aria-hidden="true"
            >
              <path d="M3 9h18l-2 11H5L3 9Zm4 0 5-6 5 6M9 12v5m6-5v5" />
            </svg>
            <span class="basket-link-label">Basket</span>
            <span class="basket-count">{{ basketCount }}</span>
          </a>
          <ProfileMenu
            :theme="theme"
            :active="active"
            :logged-in="session !== null"
            @toggle-theme="toggleTheme"
            @sign-out="signOut"
          />
        </div>
      </div>
    </header>
    <main
      id="main-content"
      tabindex="-1"
      :class="{
        'home-main': active === 'home',
        'transaction-main': active !== 'home',
      }"
    >
      <template v-if="active === 'home'">
        <section class="monastery-banner" aria-label="Monastery illustration">
          <img
            class="banner-atmosphere"
            src="/home-atmosphere.png"
            alt=""
            fetchpriority="high"
          />
          <img
            class="banner-monastery"
            src="/home-monastery.png"
            alt="A fictional Mongolian Buddhist monastery with amber-lit roofs against a dark mountain landscape."
            fetchpriority="high"
          />
          <div class="banner-shade" aria-hidden="true"></div>
        </section>
        <section
          id="readings"
          class="home-catalog"
          tabindex="-1"
          aria-labelledby="catalog-heading"
        >
          <Monasteries :selected-id="monasteryId" />
        </section>
      </template>
      <div
        v-else-if="active === 'login' || active === 'signup'"
        class="catalog-view"
      >
        <LoginPage
          :key="active"
          :signup="active === 'signup'"
          :registered-email="active === 'login' ? registeredEmail : ''"
          @registered="finishRegistration"
          @signed-in="signIn"
        />
      </div>
      <div v-else-if="active === 'profile' && session" class="catalog-view">
        <UserProfile :email="session.email" />
      </div>
      <div v-else-if="active === 'requests' && session" class="catalog-view">
        <a class="return-link" href="#readings">← Browse monasteries</a>
        <PrayerRequests
          :request-id="requestId"
          :session="session"
          @session-expired="expireSession"
        />
      </div>
      <div
        v-else-if="active === 'basket' || active === 'checkout'"
        class="catalog-view"
      >
        <a class="return-link" href="#readings"
          ><span aria-hidden="true">←</span> Browse monasteries</a
        >
        <Basket :checkout="active === 'checkout'" />
      </div>
    </main>
  </div>
</template>
