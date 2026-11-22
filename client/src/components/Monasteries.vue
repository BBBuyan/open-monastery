<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { getMonasteries } from '../api'
import { createCatalog } from '../catalog'
import { monasteryCatalogLink } from '../monasteries'

const {
  books: monasteries,
  loading,
  error,
  hasMore,
  loadMore,
  reset,
  dispose,
} = createCatalog(getMonasteries)
const query = ref('')
const searchPending = ref(false)
let searchTimer: ReturnType<typeof setTimeout> | undefined

watch(query, (value) => {
  clearTimeout(searchTimer)
  searchPending.value = true
  void reset(value, false)
  searchTimer = setTimeout(() => {
    searchPending.value = false
    void loadMore()
  }, 300)
})

onMounted(() => void reset())
onUnmounted(() => {
  clearTimeout(searchTimer)
  dispose()
})
</script>

<template>
  <div class="monastery-catalog-heading">
    <h1 id="catalog-heading">Choose a monastery</h1>
    <p>Find a monastery and explore its books.</p>
  </div>
  <div class="monastery-search library-toolbar">
    <label for="monastery-search">
      Search monasteries
      <div class="search-field">
        <svg
          viewBox="0 0 24 24"
          width="22"
          height="22"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          aria-hidden="true"
        >
          <circle cx="10.5" cy="10.5" r="6.5" />
          <path d="m16 16 5 5" />
        </svg>
        <input
          id="monastery-search"
          v-model="query"
          type="search"
          placeholder="Search by monastery name"
          autocomplete="off"
        />
      </div>
    </label>
  </div>
  <p v-if="loading || searchPending" role="status">Loading monasteries.</p>
  <div v-if="error && !monasteries.length" class="state-panel" role="status">
    <p>The monasteries couldn't be loaded. Please try again.</p>
    <button class="button secondary" type="button" @click="loadMore">
      Try again
    </button>
  </div>
  <div
    v-else-if="!loading && !searchPending && !monasteries.length"
    class="state-panel"
    role="status"
  >
    <p>
      {{
        query.trim()
          ? 'No matching monasteries.'
          : 'No monasteries are listed yet.'
      }}
    </p>
    <button
      v-if="query.trim()"
      class="button secondary"
      type="button"
      @click="query = ''"
    >
      Clear search
    </button>
  </div>
  <ul class="monastery-directory" :aria-busy="loading || searchPending">
    <li v-for="entry in monasteries" :key="entry.id">
      <a :href="monasteryCatalogLink(entry.id)" class="monastery-card">
        <img
          class="monastery-logo"
          src="/monastery-placeholder.png"
          width="56"
          height="56"
          alt=""
        />
        <span class="monastery-card-copy">
          <span class="monastery-card-name">{{ entry.name }}</span>
          <span class="monastery-card-action">Browse books</span>
        </span>
      </a>
    </li>
  </ul>
  <div v-if="monasteries.length" class="catalog-pagination">
    <p v-if="error" role="status">
      Couldn't load more monasteries. Please try again.
    </p>
    <button
      v-if="hasMore"
      class="button secondary"
      type="button"
      :disabled="loading || searchPending"
      @click="loadMore"
    >
      {{
        loading
          ? 'Loading monasteries.'
          : error
            ? 'Try again'
            : 'Load more monasteries'
      }}
    </button>
  </div>
</template>
