import { computed, ref, shallowRef } from 'vue'
import type { Book } from './api'

export const MAX_QUANTITY = 99
export const BASKET_KEY = 'open-monastery:basket:v1'
export type BasketBook = Pick<Book, 'id' | 'name' | 'price'>
export type BasketItem = { book: BasketBook; quantity: number }
type BasketStorage = Pick<Storage, 'getItem' | 'setItem'>

function validBook(value: unknown): value is BasketBook {
  const book = value as Partial<BasketBook> | null
  return (
    !!book &&
    Number.isSafeInteger(book.id) &&
    Number(book.id) > 0 &&
    typeof book.name === 'string' &&
    Number.isSafeInteger(book.price) &&
    Number(book.price) >= 0
  )
}
function validQuantity(quantity: number) {
  return Number.isInteger(quantity) && quantity >= 1 && quantity <= MAX_QUANTITY
}

export function createBasket(storage?: BasketStorage) {
  const lines = shallowRef<BasketItem[]>([])
  const persistent = ref(!!storage)
  const announcement = ref('')
  let savedValue: string | null = null
  try {
    savedValue = storage?.getItem(BASKET_KEY) ?? null
  } catch {
    persistent.value = false
  }
  try {
    const saved = JSON.parse(savedValue ?? 'null')
    if (
      (saved?.version === 1 || saved?.version === 2) &&
      Array.isArray(saved.items)
    ) {
      const restored = new Map<number, BasketItem>()
      for (const item of saved.items) {
        const book = saved.version === 1 ? item?.book : item
        if (!validBook(book) || !validQuantity(item?.quantity)) continue
        const previous = restored.get(book.id)
        restored.set(book.id, {
          book: { id: book.id, name: book.name, price: book.price },
          quantity: Math.min(
            MAX_QUANTITY,
            (previous?.quantity ?? 0) + item.quantity,
          ),
        })
      }
      save([...restored.values()])
    }
  } catch {
    // Corrupt saved data must not stop browsing; a subsequent change replaces it.
  }
  function save(next: BasketItem[]) {
    lines.value = next
    try {
      const items = next.map(({ book, quantity }) => ({
        id: book.id,
        quantity,
        name: book.name,
        price: book.price,
      }))
      storage?.setItem(BASKET_KEY, JSON.stringify({ version: 2, items }))
      persistent.value = !!storage
    } catch {
      persistent.value = false
    }
  }
  function quantityOf(id: number) {
    return lines.value.find((item) => item.book.id === id)?.quantity ?? 0
  }
  function add(book: BasketBook) {
    if (!validBook(book) || quantityOf(book.id) >= MAX_QUANTITY) return false
    const quantity = quantityOf(book.id) + 1
    const item = {
      book: { id: book.id, name: book.name, price: book.price },
      quantity,
    }
    save(
      quantity === 1
        ? [...lines.value, item]
        : lines.value.map((line) => (line.book.id === book.id ? item : line)),
    )
    announcement.value = `${book.name}: ${quantity} in basket.`
    return true
  }
  function setQuantity(id: number, quantity: number) {
    if (quantity === 0 && quantityOf(id)) {
      remove(id)
      return true
    }
    if (!validQuantity(quantity) || !quantityOf(id)) return false
    save(
      lines.value.map((line) =>
        line.book.id === id ? { ...line, quantity } : line,
      ),
    )
    announcement.value = `Quantity updated to ${quantity}.`
    return true
  }
  function remove(id: number) {
    save(lines.value.filter((line) => line.book.id !== id))
    announcement.value = 'Book removed from basket.'
  }
  return {
    items: computed(() => lines.value),
    count: computed(() =>
      lines.value.reduce((sum, line) => sum + line.quantity, 0),
    ),
    total: computed(() =>
      lines.value.reduce(
        (sum, line) => sum + line.book.price * line.quantity,
        0,
      ),
    ),
    persistent,
    announcement,
    quantityOf,
    add,
    setQuantity,
    remove,
  }
}
function browserStorage() {
  try {
    return typeof window === 'undefined' ? undefined : window.localStorage
  } catch {
    return undefined
  }
}
export const basket = createBasket(browserStorage())
