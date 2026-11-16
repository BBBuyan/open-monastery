<script setup lang="ts">
import { computed } from 'vue'
import { basket, basketBookKey } from '../basket'
import { demoMonasteries, monasteryCatalogLink } from '../demo-catalog'
import { formatPrice } from '../money'
import BasketQuantity from './BasketQuantity.vue'
defineProps<{ checkout?: boolean }>()
const { items, count, total, persistent } = basket
const lines = computed(() =>
  items.value.map((item) => ({
    ...item,
    monastery: demoMonasteries.find(({ id }) => id === item.book.monasteryId),
  })),
)
</script>

<template>
  <section :aria-labelledby="checkout ? 'checkout-title' : 'basket-title'">
    <div class="basket-heading">
      <h1 :id="checkout ? 'checkout-title' : 'basket-title'">
        {{ checkout ? 'Checkout' : 'Your basket' }}
      </h1>
      <p>{{ count }} {{ count === 1 ? 'book' : 'books' }}</p>
    </div>
    <div v-if="!items.length" class="state-panel">
      <h2>Your basket is empty</h2>
      <p>Add books to review them here.</p>
      <a class="button primary" href="#readings">Browse monasteries</a>
    </div>
    <template v-else>
      <p v-if="!persistent" class="basket-notice" role="status">
        Your basket is saved for this visit only.
      </p>
      <div class="basket-layout">
        <ul class="basket-list">
          <li
            v-for="item in lines"
            :key="basketBookKey(item.book.id, item.book.monasteryId)"
          >
            <div class="book-heading">
              <h2 lang="mn">{{ item.book.name }}</h2>
              <p class="book-price">
                {{ formatPrice(item.book.price * item.quantity) }}
              </p>
            </div>
            <p class="unit-price">
              {{ formatPrice(item.book.price) }} each
              <span aria-hidden="true"> · </span>
              <a
                v-if="item.monastery"
                :href="monasteryCatalogLink(item.monastery.id)"
              >
                {{ item.monastery.name }}
              </a>
              <span v-else>Monastery unspecified</span>
            </p>
            <div class="basket-line-actions">
              <BasketQuantity :book="item.book" />
              <button
                type="button"
                class="remove-book"
                :aria-label="`Remove ${item.book.name}${item.monastery ? ` from ${item.monastery.name}` : ''} from basket`"
                @click="basket.remove(item.book.id, item.book.monasteryId)"
              >
                Remove
              </button>
            </div>
          </li>
        </ul>
        <aside class="basket-summary" aria-label="Basket summary">
          <h2>{{ checkout ? 'Review your selection' : 'Summary' }}</h2>
          <div class="basket-total">
            <span>Total</span><strong>{{ formatPrice(total) }}</strong>
          </div>
          <template v-if="checkout">
            <p class="checkout-notice">
              Demo checkout. You can review your selection here.
            </p>
            <a class="basket-secondary-link" href="#readings"
              >Continue browsing</a
            >
          </template>
          <template v-else>
            <a class="button primary" href="#checkout"
              >Go to checkout <span aria-hidden="true">→</span></a
            >
          </template>
        </aside>
      </div>
    </template>
  </section>
</template>
