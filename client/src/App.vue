<script setup lang="ts">
import { ref } from 'vue'

const sections = [
  { id: 'home', label: 'Home' },
  { id: 'readings', label: 'Readings' },
  { id: 'about', label: 'About' },
] as const

type Section = (typeof sections)[number]['id']

const activeSection = ref<Section>('home')
</script>

<template>
  <div class="min-h-screen bg-cream text-ink">
    <header class="border-b border-[#344c3b]/20 bg-cream-dark">
      <div
        class="mx-auto flex min-h-[88px] max-w-[1600px] flex-wrap items-center gap-4 px-6 py-4 lg:flex-nowrap lg:px-12"
      >
        <div class="inline-flex shrink-0 items-center gap-3">
          <span class="grid size-[50px] place-items-center" aria-hidden="true">
            <img
              class="size-full"
              src="/open-monastery-icon-yellow.svg"
              alt=""
            />
          </span>
          <strong class="font-brand text-2xl font-semibold tracking-wide">
            Open Monastery
          </strong>
        </div>

        <nav
          class="order-3 flex w-full justify-center lg:order-none lg:w-auto lg:flex-1"
          aria-label="Main sections"
          role="tablist"
        >
          <div
            class="inline-flex rounded-full border border-monastery/15 bg-cream p-1"
          >
            <button
              v-for="section in sections"
              :id="`${section.id}-tab`"
              :key="section.id"
              type="button"
              role="tab"
              :aria-selected="activeSection === section.id"
              :aria-controls="`${section.id}-panel`"
              class="rounded-full px-5 py-2 text-sm font-medium transition-colors"
              :class="
                activeSection === section.id
                  ? 'bg-monastery text-cream'
                  : 'text-muted hover:bg-cream-dark hover:text-ink'
              "
              @click="activeSection = section.id"
            >
              {{ section.label }}
            </button>
          </div>
        </nav>

        <div class="ml-auto flex items-center gap-3">
          <button
            type="button"
            class="h-10 shrink-0 rounded-full bg-monastery px-5 text-sm font-semibold text-cream transition-opacity hover:opacity-85"
          >
            Log in
          </button>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-6 py-16 lg:px-12">
      <section
        v-if="activeSection === 'home'"
        id="home-panel"
        role="tabpanel"
        aria-labelledby="home-tab"
      >
        <p class="text-sm font-semibold tracking-widest text-muted uppercase">
          Home
        </p>
        <h1 class="mt-3 text-4xl font-semibold">Welcome to Open Monastery</h1>
        <p class="mt-4 max-w-2xl text-muted">
          The home view will be developed here.
        </p>
      </section>

      <section
        v-else-if="activeSection === 'readings'"
        id="readings-panel"
        role="tabpanel"
        aria-labelledby="readings-tab"
      >
        <p class="text-sm font-semibold tracking-widest text-muted uppercase">
          Readings
        </p>
        <h1 class="mt-3 text-4xl font-semibold">Reading library</h1>
        <p class="mt-4 max-w-2xl text-muted">
          The readings view will be developed here.
        </p>
      </section>

      <section
        v-else
        id="about-panel"
        role="tabpanel"
        aria-labelledby="about-tab"
      >
        <p class="text-sm font-semibold tracking-widest text-muted uppercase">
          About
        </p>
        <h1 class="mt-3 text-4xl font-semibold">About the monastery</h1>
        <p class="mt-4 max-w-2xl text-muted">
          The about view will be developed here.
        </p>
      </section>
    </main>
  </div>
</template>
