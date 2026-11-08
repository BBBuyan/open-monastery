<script setup lang="ts">
import { computed } from 'vue'
import { basket } from '../basket'
import { formatPrice } from '../money'
import BasketQuantity from './BasketQuantity.vue'
import CheckoutPayment from './CheckoutPayment.vue'
defineProps<{ checkout?: boolean }>()
const { items, count, total, persistent } = basket
const selectionKey = computed(() =>
  JSON.stringify(
    items.value.map((item) => [item.book.id, item.quantity, item.book.price]),
  ),
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
      <a class="button primary" href="#readings">Browse books</a>
    </div>
    <template v-else>
      <p v-if="!persistent" class="basket-notice" role="status">
        Your basket is saved for this visit only.
      </p>
      <div class="basket-layout">
        <ul class="basket-list">
          <li v-for="item in items" :key="item.book.id">
            <div class="book-heading">
              <h2 lang="mn">{{ item.book.name }}</h2>
              <p class="book-price">
                {{ formatPrice(item.book.price * item.quantity) }}
              </p>
            </div>
            <p class="unit-price">{{ formatPrice(item.book.price) }} each</p>
            <div class="basket-line-actions">
              <BasketQuantity :book="item.book" />
              <button
                type="button"
                class="remove-book"
                :aria-label="`Remove ${item.book.name} from basket`"
                @click="basket.remove(item.book.id)"
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
            <CheckoutPayment :total="total" :selection-key="selectionKey" />
            <a class="basket-secondary-link" href="#readings"
              >Continue browsing</a
            >
          </template>
          <template v-else>
            <a class="button primary" href="#checkout"
              >Go to checkout <span aria-hidden="true">→</span></a
            >
            <p class="checkout-notice">Review your books before payment.</p>
          </template>
        </aside>
      </div>
    </template>
  </section>
</template>
