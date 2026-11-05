<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import Readings from './components/Readings.vue'
import About from './components/About.vue'

const sections = [
  { id: 'home', label: 'Home' },
  { id: 'readings', label: 'Books' },
  { id: 'about', label: 'About' },
] as const
const active = ref<string>('home')
function sync() {
  const section =
    sections.find((section) => section.id === location.hash.slice(1)) ??
    sections[0]
  active.value = section.id
  document.title = section.label + ' · Open Monastery'
}
onMounted(() => {
  sync()
  window.addEventListener('hashchange', sync)
})
onUnmounted(() => window.removeEventListener('hashchange', sync))
</script>

<template>
  <a class="skip-link" href="#main-content">Skip to content</a>
  <div class="app-shell">
    <header>
      <div class="header-inner">
        <a class="brand" href="#home" aria-label="Open Monastery home">
          <img
            src="/open-monastery-icon-yellow.svg"
            width="44"
            height="44"
            alt=""
          />
          <span>Open Monastery</span>
        </a>
        <nav aria-label="Main navigation">
          <a
            v-for="section in sections"
            :key="section.id"
            :href="'#' + section.id"
            :aria-current="active === section.id ? 'page' : undefined"
            >{{ section.label }}</a
          >
        </nav>
      </div>
    </header>
    <main id="main-content" tabindex="-1">
      <template v-if="active === 'home'">
        <section class="welcome" aria-labelledby="welcome-title">
          <div>
            <p class="eyebrow">Buddhist prayers in Mongolia</p>
            <h1 id="welcome-title">
              Prayers for you.<br />For those you hold dear.
            </h1>
            <p class="intro">
              In Mongolia’s Buddhist tradition, monks recite sacred texts and
              perform rituals on behalf of people and their loved ones. Find a
              prayer for the intention you hold.
            </p>
            <div class="actions">
              <a class="button primary" href="#readings"
                >Browse books <span aria-hidden="true">→</span></a
              >
              <a class="text-link" href="#about">About prayer requests</a>
            </div>
          </div>
          <aside
            class="welcome-note"
            aria-label="The tradition of prayer recitation"
          >
            <p class="eyebrow" lang="mn">Ном хуруулах</p>
            <h2>Recited by monks.<br />Dedicated to others.</h2>
            <p>
              Prayers are requested with intentions such as wellbeing,
              protection, removing obstacles, and remembrance of those who have
              died.
            </p>
          </aside>
        </section>
        <section class="overview" aria-labelledby="discover-title">
          <div class="overview-heading">
            <p class="eyebrow">Prayer offerings</p>
            <h2 id="discover-title">Choose a prayer with care</h2>
          </div>
          <a class="destination" href="#readings">
            <div>
              <h3>Books</h3>
              <p>
                Browse the Mongolian names and traditional purposes of the
                prayers monks recite.
              </p>
            </div>
            <span aria-hidden="true">→</span>
          </a>
          <a class="destination" href="#about">
            <div>
              <h3>Praying for someone</h3>
              <p>
                Learn about requesting a recitation for yourself, your family,
                or someone you wish to remember.
              </p>
            </div>
            <span aria-hidden="true">→</span>
          </a>
        </section>
      </template>
      <Readings v-else-if="active === 'readings'" />
      <About v-else />
    </main>
    <footer>
      <span>Open Monastery</span
      ><span>Prayers for people and their loved ones.</span>
    </footer>
  </div>
</template>
