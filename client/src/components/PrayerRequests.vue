<script setup lang="ts">
import { computed, nextTick, onUnmounted, ref, watch } from 'vue'
import { getUserReadings } from '../api'
import {
  demoRequests,
  requestDate,
  requestStatus,
  type PrayerRequest,
} from '../prayer-requests'

const props = defineProps<{
  requestId?: number
  userId?: number
  embedded?: boolean
  demoAccount?: string
}>()
const account = ref('history')
watch(
  () => props.demoAccount,
  (value) => {
    account.value = value ?? 'history'
  },
  { immediate: true },
)
const records = ref<PrayerRequest[]>([])
const loading = ref(false)
const error = ref(false)
const hasMore = ref(false)
const loaded = ref(false)
let cursor: string | undefined
let controller: AbortController | undefined
let generation = 0
const selected = computed(() =>
  records.value.find((record) => record.id === props.requestId),
)

async function load(append = false) {
  if (loading.value) return
  const current = generation
  controller = new AbortController()
  const signal = controller.signal
  loading.value = true
  error.value = false
  try {
    let next: PrayerRequest[]
    let nextCursor: string | undefined
    let more = false
    if (props.userId !== undefined) {
      const page = await getUserReadings(props.userId, {
        cursor: append ? cursor : undefined,
        signal,
      })
      // A missing catalog entry must never hide an existing prayer request.
      const names = await Promise.all(
        [...new Set(page.data.map((item) => item.book_id))].map(async (id) => {
          try {
            const base = (
              import.meta.env.VITE_API_BASE_URL || '/api/v1'
            ).replace(/\/$/, '')
            const response = await fetch(`${base}/books/${id}`, {
              signal,
              headers: { Accept: 'application/json' },
            })
            if (!response.ok) return [id, `Book #${id}`] as const
            const book = await response.json()
            return [
              id,
              typeof book.name === 'string' && book.name
                ? book.name
                : `Book #${id}`,
            ] as const
          } catch {
            return [id, `Book #${id}`] as const
          }
        }),
      )
      const lookup = new Map(names)
      next = page.data.map((item) => ({
        ...item,
        book_name: lookup.get(item.book_id)!,
      }))
      more = page.pagination.hasMore
      nextCursor = page.pagination.nextCursor ?? undefined
    } else {
      if (account.value === 'unavailable')
        throw new Error('Demo unavailable state')
      next =
        account.value === 'empty'
          ? []
          : demoRequests.map((item) => ({ ...item }))
    }
    if (current !== generation || signal.aborted) return
    const existing = new Set(append ? records.value.map((item) => item.id) : [])
    records.value = [
      ...(append ? records.value : []),
      ...next.filter((item) => !existing.has(item.id)),
    ]
    cursor = nextCursor
    hasMore.value = more
    loaded.value = true
  } catch {
    if (current === generation && !signal.aborted) error.value = true
  } finally {
    if (current === generation) loading.value = false
  }
}
watch(
  [account, () => props.userId],
  () => {
    generation++
    controller?.abort()
    loading.value = false
    records.value = []
    cursor = undefined
    hasMore.value = false
    loaded.value = false
    void load()
  },
  { immediate: true },
)
watch(
  () => props.requestId,
  async () => {
    await nextTick()
    window.scrollTo({ top: 0, behavior: 'instant' })
    document.getElementById('main-content')?.focus({ preventScroll: true })
  },
)
onUnmounted(() => {
  generation++
  controller?.abort()
})
</script>

<template>
  <section aria-labelledby="requests-title" class="requests-page">
    <div v-if="userId === undefined && !embedded" class="request-demo">
      <p>
        <strong>Demo preview</strong> · Fictional requests. No account is signed
        in.
      </p>
      <label for="demo-account">Demo account</label>
      <select id="demo-account" v-model="account">
        <option value="history">Sample account · Request history</option>
        <option value="empty">Sample account · No requests</option>
        <option value="unavailable">
          Sample account · Service unavailable
        </option>
      </select>
    </div>
    <a v-if="requestId" class="return-link" href="#requests"
      >← All prayer requests</a
    >
    <div class="requests-heading">
      <div>
        <component :is="embedded ? 'h2' : 'h1'" id="requests-title">
          {{ requestId ? `Prayer request #${requestId}` : 'My Readings' }}
        </component>
        <p v-if="!requestId">Check the status of prayers you’ve requested.</p>
      </div>
      <button
        class="button secondary"
        type="button"
        :disabled="loading"
        @click="load()"
      >
        {{ loading ? 'Refreshing…' : 'Refresh' }}
      </button>
    </div>
    <p class="visually-hidden" role="status">
      {{
        loading
          ? 'Loading prayer requests…'
          : error
            ? 'Prayer requests could not be updated.'
            : loaded
              ? 'Prayer requests loaded.'
              : ''
      }}
    </p>
    <div v-if="error" class="state-panel" role="status">
      <h2>Prayer requests couldn’t be loaded</h2>
      <p>
        {{
          records.length
            ? 'Your previously loaded requests are still shown. Try refreshing again.'
            : 'Please try again in a moment.'
        }}
      </p>
      <button
        class="button secondary"
        :disabled="loading"
        @click="load(hasMore)"
      >
        Try again
      </button>
    </div>
    <div :aria-busy="loading">
      <ul
        v-if="loading && !records.length"
        class="reading-list catalog-skeleton"
        aria-hidden="true"
      >
        <li v-for="index in 3" :key="index">
          <article>
            <span class="skeleton-line skeleton-title"></span
            ><span class="skeleton-line skeleton-short"></span>
          </article>
        </li>
      </ul>
      <template v-else-if="requestId && selected">
        <article class="request-detail">
          <h2 lang="mn">{{ selected.book_name }}</h2>
          <span class="request-status">{{
            requestStatus(selected.status).label
          }}</span>
          <p>{{ requestStatus(selected.status).description }}</p>
          <dl>
            <div>
              <dt>Request reference</dt>
              <dd>#{{ selected.id }}</dd>
            </div>
            <div>
              <dt>Requested</dt>
              <dd>
                {{ requestDate(selected.created_at) ?? 'Date unavailable' }}
              </dd>
            </div>
            <div v-if="selected.monk_id">
              <dt>Monk</dt>
              <dd>{{ selected.monk_name ?? 'Assigned · Name unavailable' }}</dd>
            </div>
            <div v-if="requestDate(selected.assigned_at)">
              <dt>Assigned</dt>
              <dd>{{ requestDate(selected.assigned_at) }}</dd>
            </div>
          </dl>
        </article>
      </template>
      <div v-else-if="requestId && loaded" class="state-panel">
        <h2>This request isn’t in the loaded history</h2>
        <p>Return to your prayer requests to find it.</p>
        <a class="button secondary" href="#requests">View prayer requests</a>
      </div>
      <ul v-else-if="records.length" class="request-list">
        <li v-for="record in records" :key="record.id">
          <article>
            <div class="request-row-heading">
              <h2 lang="mn">{{ record.book_name }}</h2>
              <span class="request-status">{{
                requestStatus(record.status).label
              }}</span>
            </div>
            <p class="request-date">
              {{
                requestDate(record.created_at)
                  ? `Requested ${requestDate(record.created_at)}`
                  : 'Request date unavailable'
              }}
            </p>
            <p>{{ requestStatus(record.status).description }}</p>
            <div class="request-row-footer">
              <span>Request #{{ record.id }}</span
              ><a
                :href="`#requests/${record.id}`"
                :aria-label="`View details for ${record.book_name}, request ${record.id}`"
                >View details <span aria-hidden="true">→</span></a
              >
            </div>
          </article>
        </li>
      </ul>
      <div v-else-if="loaded && !error" class="state-panel">
        <h2>You haven’t requested any prayers yet</h2>
        <p>Explore the books offered for prayer recitations and rituals.</p>
        <a class="button primary" href="#readings">Browse books</a>
      </div>
      <div v-if="hasMore && !requestId" class="catalog-pagination">
        <button
          class="button secondary"
          :disabled="loading"
          @click="load(true)"
        >
          {{ loading ? 'Loading…' : 'Load more requests' }}
        </button>
      </div>
    </div>
  </section>
</template>
