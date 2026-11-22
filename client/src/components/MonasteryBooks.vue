<script setup lang="ts">
import { ref, watch } from 'vue'
import { ApiError, getMonastery, type Monastery } from '../api'
import Readings from './Readings.vue'

const props = defineProps<{ monasteryId?: number }>()
const monastery = ref<Monastery>()
const loading = ref(false)
const error = ref(false)
const notFound = ref(false)
const retry = ref(0)

watch(
  [() => props.monasteryId, retry],
  async ([id], _, onCleanup) => {
    monastery.value = undefined
    error.value = false
    notFound.value = id === undefined
    loading.value = false
    if (id === undefined) return
    const controller = new AbortController()
    onCleanup(() => controller.abort())
    loading.value = true
    try {
      const result = await getMonastery(id, controller.signal)
      if (!controller.signal.aborted) {
        monastery.value = result
        document.title = `${result.name} – Open Monastery`
      }
    } catch (cause) {
      if (!controller.signal.aborted) {
        notFound.value = cause instanceof ApiError && cause.status === 404
        error.value = !notFound.value
      }
    } finally {
      if (!controller.signal.aborted) loading.value = false
    }
  },
  { immediate: true },
)
</script>

<template>
  <a class="return-link" href="#readings">All monasteries</a>
  <p v-if="loading" role="status">Loading monastery.</p>
  <div v-else-if="notFound" class="state-panel" role="alert">
    <h1>Monastery not found</h1>
    <p>Return to all monasteries to choose another.</p>
  </div>
  <div v-else-if="error" class="state-panel" role="alert">
    <h1>This monastery couldn't be loaded</h1>
    <p>Please try again.</p>
    <button class="button secondary" type="button" @click="retry++">
      Try again
    </button>
  </div>
  <template v-else-if="monastery">
    <div class="monastery-page-heading">
      <img
        class="monastery-logo"
        src="/monastery-placeholder.png"
        width="56"
        height="56"
        alt=""
      />
      <div>
        <h1>{{ monastery.name }}</h1>
        <p>Browse books and choose a reading.</p>
      </div>
    </div>
    <Readings :key="monastery.id" :monastery="monastery" />
  </template>
</template>
