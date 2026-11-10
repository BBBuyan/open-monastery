<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import Readings from './components/Readings.vue'
import Basket from './components/Basket.vue'
import UserProfile from './components/UserProfile.vue'
import PrayerRequests from './components/PrayerRequests.vue'
import ProfileMenu from './components/ProfileMenu.vue'
import { basket } from './basket'
import { formatPrice } from './money'
import { useTheme } from './theme'
const { theme, toggle: toggleTheme } = useTheme()

const { count: basketCount, total: basketTotal, announcement } = basket
function pageFromHash() {
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
const requestId = ref<number>()
async function sync() {
  if (location.hash === '#main-content') return
  const previous = active.value
  active.value = pageFromHash()
  requestId.value = /^#(?:profile|requests)\/\d+$/.test(location.hash)
    ? Number(location.hash.split('/')[1])
    : undefined
  document.title =
    active.value === 'home'
      ? 'Open Monastery'
      : `${active.value === 'requests' ? (requestId.value ? `Prayer request #${requestId.value}` : 'Readings') : active.value === 'profile' ? 'Profile Settings' : active.value === 'basket' ? 'Basket' : 'Checkout'} · Open Monastery`
  await nextTick()
  if (location.hash === '#readings') {
    const catalog = document.getElementById('readings')
    catalog?.scrollIntoView({ block: 'start' })
    catalog?.focus({ preventScroll: true })
  } else if (previous !== active.value || location.hash === '#home') {
    window.scrollTo({ top: 0, behavior: 'instant' })
    document.getElementById('main-content')?.focus({ preventScroll: true })
  }
}
onMounted(() => {
  void sync()
  window.addEventListener('hashchange', sync)
})
onUnmounted(() => window.removeEventListener('hashchange', sync))
</script>

<template>
  <a class="skip-link" href="#main-content">Skip to content</a>
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
            :aria-current="active === 'basket' ? 'page' : undefined"
          >
            <span class="basket-link-heading">
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
              Basket <span class="basket-count">{{ basketCount }}</span>
            </span>
            <span class="basket-link-total">{{
              formatPrice(basketTotal)
            }}</span>
          </a>
          <ProfileMenu
            :theme="theme"
            :active="active"
            @toggle-theme="toggleTheme"
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
        <h1 class="visually-hidden">Open Monastery</h1>
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
          aria-label="Books"
        >
          <Readings />
        </section>
      </template>
      <div v-else-if="active === 'profile'" class="catalog-view">
        <UserProfile />
      </div>
      <div v-else-if="active === 'requests'" class="catalog-view">
        <a class="return-link" href="#readings">← Back to books</a>
        <PrayerRequests :request-id="requestId" />
      </div>
      <div v-else class="catalog-view">
        <a class="return-link" href="#readings"
          ><span aria-hidden="true">←</span> Back to books</a
        >
        <Basket :checkout="active === 'checkout'" />
      </div>
    </main>
  </div>
</template>
