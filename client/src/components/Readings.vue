<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { getBooks, getCategories, type Category } from '../api'
import { createCatalog } from '../catalog'

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
const sentinel = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | undefined
let searchTimer: ReturnType<typeof setTimeout> | undefined

watch(
  [query, selectedCategory],
  ([value, category], [, previousCategory]) => {
    clearTimeout(searchTimer)
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
  <section aria-labelledby="book-catalog-title">
    <div class="page-heading">
      <p class="eyebrow">Book catalog</p>
      <h1 id="book-catalog-title">Books</h1>
      <p class="intro">
        Find a Buddhist scripture for monks to recite on behalf of you or your
        loved ones. Browse the original Mongolian names and the traditional
        purposes associated with each book.
      </p>
    </div>
    <div class="context-note">
      <strong>About the descriptions</strong>
      <p>
        The catalog describes the intentions and spiritual benefits associated
        with these recitations and rituals in Buddhist tradition.
      </p>
      <p>
        Online prayer requests are not available yet. You can browse the books
        here.
      </p>
    </div>
    <div class="library-toolbar">
      <label for="book-search"
        >Search books
        <input
          id="book-search"
          v-model="query"
          type="search"
          placeholder="Enter a Mongolian name or purpose"
          aria-describedby="book-search-hint"
        />
        <span id="book-search-hint" class="search-hint"
          >For example: <span lang="mn">Манал, Ногоон дарь эх</span></span
        >
      </label>
      <label for="book-category" class="category-filter"
        >Category
        <select
          id="book-category"
          v-model="selectedCategory"
          :disabled="categoriesLoading || !categories.length"
          aria-describedby="category-help"
        >
          <option :value="undefined">All categories</option>
          <option
            v-for="category in categories"
            :key="category.id"
            :value="category.id"
            lang="mn"
          >
            {{ category.name }}
          </option>
        </select>
        <span id="category-help" class="search-hint">{{
          categoriesLoading
            ? 'Loading categories…'
            : categoriesError
              ? 'Categories are unavailable.'
              : !categories.length
                ? 'No categories available.'
                : 'Search within the selected category.'
        }}</span>
      </label>
      <p role="status">
        {{
          busy
            ? 'Loading books…'
            : error && !books.length
              ? 'Book catalog unavailable'
              : books.length +
                (books.length === 1 ? ' book loaded' : ' books loaded')
        }}
      </p>
    </div>
    <div v-if="categoriesError" class="category-error" role="status">
      <p>Categories couldn’t be loaded. You can still browse all books.</p>
      <button class="button secondary" type="button" @click="loadCategories">
        Retry categories
      </button>
    </div>
    <div v-if="busy && !books.length" class="state-panel" aria-busy="true">
      <h2>Loading books</h2>
      <p>Book names and descriptions will appear here when they are ready.</p>
    </div>
    <div v-else-if="error && !books.length" class="state-panel">
      <h2>The books couldn’t be loaded</h2>
      <p>Check your connection and try again.</p>
      <button class="button primary" type="button" @click="loadMore">
        Try again
      </button>
    </div>
    <div v-else-if="!books.length" class="state-panel">
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
    <ul v-if="books.length" class="reading-list" :aria-busy="busy">
      <li v-for="book in books" :key="book.id">
        <article>
          <h2 lang="mn">{{ book.name }}</h2>
          <p class="purpose-label">Traditional purpose</p>
          <p v-if="book.description" class="description" lang="mn">
            {{ book.description }}
          </p>
          <p v-else class="description">
            A description of this book is not yet available.
          </p>
        </article>
      </li>
    </ul>
    <div ref="sentinel" class="catalog-pagination">
      <template v-if="books.length">
        <p v-if="error" role="status">
          More books couldn’t be loaded. Your loaded books are still here.
        </p>
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
  </section>
</template>
