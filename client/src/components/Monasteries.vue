<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { getMonasteries, getMonastery, type Monastery } from '../api'
import { createCatalog } from '../catalog'
import { monasteryCatalogLink } from '../monasteries'
import Readings from './Readings.vue'

const props = defineProps<{ selectedId?: number }>()
const {
  books: monasteries,
  loading,
  error,
  hasMore,
  loadMore,
  reset,
  dispose,
} = createCatalog(getMonasteries)
const selected = ref<Monastery>()
const selectionLoading = ref(false)
const selectionError = ref(false)
const retry = ref(0)
watch(
  [() => props.selectedId, retry],
  async ([id], _, onCleanup) => {
    selected.value = undefined
    selectionError.value = false
    selectionLoading.value = false
    if (id === undefined) return
    const controller = new AbortController()
    onCleanup(() => controller.abort())
    selectionLoading.value = true
    try {
      const monastery = await getMonastery(id, controller.signal)
      if (!controller.signal.aborted) selected.value = monastery
    } catch {
      if (!controller.signal.aborted) selectionError.value = true
    } finally {
      if (!controller.signal.aborted) selectionLoading.value = false
    }
  },
  { immediate: true },
)
onMounted(() => void reset())
onUnmounted(dispose)
</script>

<template>
  <div class="monastery-catalog-heading">
    <h1 id="catalog-heading">Choose a monastery</h1>
    <p>Browse its books and categories.</p>
  </div>
  <p v-if="loading" role="status">Loading monasteries…</p>
  <div v-if="error" class="state-panel" role="status">
    <p>The monasteries couldn't be loaded. Please try again.</p>
    <button class="button secondary" type="button" @click="loadMore">
      Try again
    </button>
  </div>
  <p v-else-if="!loading && !monasteries.length" role="status">
    No monasteries are listed yet.
  </p>
  <ul class="monastery-directory">
    <li v-for="entry in monasteries" :key="entry.id">
      <a
        :href="monasteryCatalogLink(entry.id)"
        :aria-current="entry.id === selectedId ? 'true' : undefined"
      >
        <span>{{ entry.name }}</span>
        <span aria-hidden="true">{{
          entry.id === selectedId ? '✓' : '→'
        }}</span>
      </a>
    </li>
  </ul>
  <button
    v-if="hasMore && !error"
    class="button secondary"
    type="button"
    :disabled="loading"
    @click="loadMore"
  >
    Load more monasteries
  </button>
  <p v-if="selectionLoading" role="status">Loading the chosen monastery…</p>
  <div v-if="selectionError" class="state-panel" role="alert">
    <p>This monastery couldn't be loaded. Choose another or try again.</p>
    <button class="button secondary" type="button" @click="retry++">
      Try again
    </button>
  </div>
  <section
    v-if="selected"
    :key="selected.id"
    class="monastery-readings"
    :aria-label="selected.name"
  >
    <h2>{{ selected.name }}</h2>
    <Readings :monastery="selected" />
  </section>
</template>
