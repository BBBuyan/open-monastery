<script setup lang="ts">
import {
  computed,
  nextTick,
  onMounted,
  onUnmounted,
  ref,
  shallowRef,
  watch,
} from 'vue'
import { getBooks, getCategories, type Category, type Book } from '../api'
import { createCatalog } from '../catalog'
import { formatPrice } from '../money'
import BasketQuantity from './BasketQuantity.vue'
import BookDescription from './BookDescription.vue'
import CategoryFilter from './CategoryFilter.vue'

const { books, loading, error, hasMore, loadMore, reset, dispose } =
  createCatalog(({ categoryId, ...options }) =>
    getBooks(categoryId, { ...options, limit: 20 }),
  )
const query = ref('')
const selectedCategory = ref<number | undefined>()
const categories = ref<Category[]>([])
const categoriesLoading = ref(true)
const categoriesError = ref(false)
const categoryController = new AbortController()
const hasFilters = computed(
  () => !!query.value.trim() || selectedCategory.value !== undefined,
)

async function loadCategories() {
  categoriesLoading.value = true
  categoriesError.value = false
  try {
    const result = await getCategories(categoryController.signal)
    if (!categoryController.signal.aborted) categories.value = result
  } catch {
    if (!categoryController.signal.aborted) categoriesError.value = true
  } finally {
    if (!categoryController.signal.aborted) categoriesLoading.value = false
  }
}

function clearFilters() {
  query.value = ''
  selectedCategory.value = undefined
}
const searchPending = ref(false)
const busy = computed(() => loading.value || searchPending.value)
const retainedBooks = shallowRef<Book[]>([])
const displayedBooks = computed(() =>
  busy.value && !books.value.length ? retainedBooks.value : books.value,
)
const sentinel = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | undefined
let searchTimer: ReturnType<typeof setTimeout> | undefined

watch(
  [query, selectedCategory],
  ([value, category], [, previousCategory]) => {
    clearTimeout(searchTimer)
    retainedBooks.value = books.value.length
      ? books.value
      : busy.value
        ? retainedBooks.value
        : []
    if (category !== previousCategory) {
      searchPending.value = false
      void reset(value, true, category)
      return
    }
    searchPending.value = true
    void reset(value, false, category)
    searchTimer = setTimeout(() => {
      searchPending.value = false
      void loadMore()
    }, 300)
  },
  { flush: 'sync' },
)

// Re-observe after a page renders, including when the viewport is still unfilled.
watch([loading, hasMore, error, searchPending], async () => {
  await nextTick()
  if (!observer || !sentinel.value) return
  observer.unobserve(sentinel.value)
  if (hasMore.value && !busy.value && !error.value)
    observer.observe(sentinel.value)
})

onMounted(() => {
  void loadCategories()
  if ('IntersectionObserver' in window) {
    observer = new IntersectionObserver(
      (entries) => {
        if (
          entries.some((entry) => entry.isIntersecting) &&
          !busy.value &&
          !error.value
        )
          void loadMore()
      },
      { rootMargin: '300px 0px' },
    )
  }
  void reset()
})
onUnmounted(() => {
  categoryController.abort()
  clearTimeout(searchTimer)
  observer?.disconnect()
  observer = undefined
  dispose()
})
</script>

<template>
  <section aria-label="Search books">
    <div class="library-toolbar">
      <label for="book-search" class="search-label">
        Search books
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
            id="book-search"
            v-model="query"
            type="search"
            placeholder="Search by name or description"
            autocomplete="off"
          />
        </div>
      </label>
      <div class="category-filter">
        <label id="book-category-label" for="book-category">Category</label>
        <CategoryFilter
          v-model="selectedCategory"
          :categories="categories"
          :loading="categoriesLoading"
        />
      </div>
    </div>
    <p
      class="visually-hidden"
      role="status"
      aria-live="polite"
      aria-atomic="true"
    >
      {{
        busy ? 'Searching…' : error && !books.length ? 'Books unavailable' : ''
      }}
    </p>
    <div v-if="categoriesError" class="category-error" role="status">
      <p>Categories unavailable.</p>
      <button class="button secondary" type="button" @click="loadCategories">
        Retry categories
      </button>
    </div>
    <div class="catalog-results" :aria-busy="busy">
      <ul
        v-if="busy && !displayedBooks.length"
        class="reading-list catalog-skeleton"
        aria-hidden="true"
      >
        <li v-for="index in 4" :key="index">
          <article>
            <span class="skeleton-line skeleton-title"></span>
            <span class="skeleton-line"></span
            ><span class="skeleton-line skeleton-short"></span>
          </article>
        </li>
      </ul>
      <div v-else-if="!busy && error && !books.length" class="state-panel">
        <h2>The books couldn’t be loaded</h2>
        <p>Check your connection and try again.</p>
        <button class="button primary" type="button" @click="loadMore">
          Try again
        </button>
      </div>
      <div v-else-if="!busy && !books.length" class="state-panel">
        <h2>
          {{ hasFilters ? 'No matching books' : 'No books are listed yet' }}
        </h2>
        <p>
          {{
            hasFilters
              ? 'Try another category or a different Mongolian name or purpose.'
              : 'Please check back later for the book catalog.'
          }}
        </p>
        <button
          v-if="hasFilters"
          class="button secondary"
          type="button"
          @click="clearFilters"
        >
          Clear filters
        </button>
      </div>
      <ul
        v-if="displayedBooks.length"
        class="reading-list"
        :class="{ 'is-updating': busy && !books.length }"
      >
        <li v-for="book in displayedBooks" :key="book.id">
          <article class="catalog-book">
            <div class="book-copy">
              <h2 lang="mn">{{ book.name }}</h2>
              <p class="visually-hidden">Description</p>
              <BookDescription
                v-if="book.description"
                :book-id="book.id"
                :name="book.name"
                :text="book.description"
              />
              <p v-else class="description">
                A description of this book is not yet available.
              </p>
            </div>
            <div class="book-purchase">
              <p class="book-price">{{ formatPrice(book.price) }}</p>
              <BasketQuantity :book="book" :disabled="busy && !books.length" />
            </div>
          </article>
        </li>
      </ul>
      <div ref="sentinel" class="catalog-pagination">
        <template v-if="books.length">
          <p v-if="error" role="status">Couldn’t load more books. Try again.</p>
          <button
            v-if="hasMore"
            class="button secondary"
            type="button"
            :disabled="busy"
            @click="loadMore"
          >
            {{
              busy
                ? 'Loading more books…'
                : error
                  ? 'Try again'
                  : 'Load more books'
            }}
          </button>
          <p v-else role="status">
            {{
              hasFilters
                ? 'All matching books are shown.'
                : 'All books are shown.'
            }}
          </p>
        </template>
      </div>
    </div>
  </section>
</template>
