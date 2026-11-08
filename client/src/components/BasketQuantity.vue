<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { basket, MAX_QUANTITY, type BasketBook } from '../basket'
defineProps<{ book: BasketBook; disabled?: boolean }>()
const quantity = (id: number) => basket.quantityOf(id)
const increaseButton = ref<HTMLButtonElement>()
const addButton = ref<HTMLButtonElement>()
async function decrease(book: BasketBook, event: MouseEvent) {
  const restoreFocus = document.activeElement === event.currentTarget
  const remaining = quantity(book.id) - 1
  basket.setQuantity(book.id, remaining)
  if (remaining === 0 && restoreFocus) {
    await nextTick()
    addButton.value?.focus({ preventScroll: true })
  }
}
async function addBook(book: BasketBook, event: MouseEvent) {
  const restoreFocus = document.activeElement === event.currentTarget
  if (!basket.add(book)) return
  await nextTick()
  if (restoreFocus) increaseButton.value?.focus({ preventScroll: true })
}
</script>

<template>
  <button
    v-if="!quantity(book.id)"
    ref="addButton"
    class="button secondary add-book"
    type="button"
    :disabled="disabled"
    :aria-label="`Add ${book.name} to basket`"
    title="Add to basket"
    @click="addBook(book, $event)"
  >
    <span aria-hidden="true">+</span>
  </button>
  <div
    v-else
    class="quantity-control"
    role="group"
    :aria-label="`Quantity for ${book.name}`"
  >
    <button
      type="button"
      :disabled="disabled"
      :aria-label="
        quantity(book.id) === 1
          ? `Remove ${book.name} from basket`
          : `Decrease quantity of ${book.name}`
      "
      @click="decrease(book, $event)"
    >
      −
    </button>
    <output :aria-label="`Quantity: ${quantity(book.id)}`">{{
      quantity(book.id)
    }}</output>
    <button
      type="button"
      ref="increaseButton"
      :disabled="disabled || quantity(book.id) >= MAX_QUANTITY"
      :aria-label="`Increase quantity of ${book.name}`"
      @click="basket.add(book)"
    >
      +
    </button>
  </div>
</template>
