<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { getMonasteries, type Monastery } from '../api'
import { monasteryCatalogLink } from '../monasteries'
import Readings from './Readings.vue'
import SelectDropdown from './SelectDropdown.vue'

const props = defineProps<{ selectedId?: number }>()
const monasteries = ref<Monastery[]>([])
const loading = ref(true)
const error = ref(false)
const selected = computed(() =>
  props.selectedId === undefined
    ? monasteries.value[0]
    : monasteries.value.find((entry) => entry.id === props.selectedId),
)
let controller: AbortController | undefined

async function loadMonasteries() {
  controller?.abort()
  const request = new AbortController()
  controller = request
  loading.value = true
  error.value = false
  try {
    const page = await getMonasteries({ signal: request.signal })
    if (!request.signal.aborted) monasteries.value = page.data
  } catch {
    if (!request.signal.aborted) error.value = true
  } finally {
    if (!request.signal.aborted) loading.value = false
  }
}

function selectMonastery(id: number | undefined) {
  if (id !== undefined) location.hash = monasteryCatalogLink(id)
}

onMounted(() => void loadMonasteries())
onUnmounted(() => controller?.abort())
</script>

<template>
  <h1 class="visually-hidden">Books</h1>
  <div class="monastery-selector">
    <label id="monastery-select-label">Monastery</label>
    <SelectDropdown
      id="monastery-select"
      label-id="monastery-select-label"
      :model-value="selected?.id"
      :options="monasteries"
      :loading="loading"
      placeholder="Select a monastery"
      loading-label="Loading monasteries…"
      @update:model-value="selectMonastery"
    />
    <p v-if="error" role="status">
      The monasteries couldn't be loaded.
      <button
        class="button secondary"
        type="button"
        :disabled="loading"
        @click="loadMonasteries"
      >
        Try again
      </button>
    </p>
    <p v-else-if="!loading && !monasteries.length" role="status">
      No monasteries are listed yet.
    </p>
    <p v-else-if="!loading && !selected" role="status">
      This monastery is unavailable. Select another monastery.
    </p>
  </div>
  <Readings v-if="selected" :key="selected.id" :monastery="selected" />
</template>
