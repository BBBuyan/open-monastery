import { computed, ref, shallowRef } from 'vue'
import type { Book } from './api'

export const MAX_QUANTITY = 99
export const BASKET_KEY = 'open-monastery:basket:v1'
export type BasketBook = Pick<Book, 'id' | 'name' | 'price'> & {
  monasteryId?: string
}
export type BasketItem = { book: BasketBook; quantity: number }
type BasketStorage = Pick<Storage, 'getItem' | 'setItem'>

export function basketBookKey(id: number, monasteryId?: string) {
  return monasteryId === undefined ? String(id) : `${id}:${monasteryId}`
}

function snapshotBook({
  id,
  name,
  price,
  monasteryId,
}: BasketBook): BasketBook {
  return monasteryId === undefined
    ? { id, name, price }
    : { id, name, price, monasteryId }
}

function validBook(value: unknown): value is BasketBook {
  const book = value as Partial<BasketBook> | null
  return (
    !!book &&
    Number.isSafeInteger(book.id) &&
    Number(book.id) > 0 &&
    typeof book.name === 'string' &&
    (book.monasteryId === undefined ||
      (typeof book.monasteryId === 'string' &&
        book.monasteryId.trim() !== '')) &&
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
      const restored = new Map<string, BasketItem>()
      for (const item of saved.items) {
        const book = saved.version === 1 ? item?.book : item
        if (!validBook(book) || !validQuantity(item?.quantity)) continue
        const key = basketBookKey(book.id, book.monasteryId)
        const previous = restored.get(key)
        restored.set(key, {
          book: snapshotBook(book),
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
        ...book,
        quantity,
      }))
      storage?.setItem(BASKET_KEY, JSON.stringify({ version: 2, items }))
      persistent.value = !!storage
    } catch {
      persistent.value = false
    }
  }
  function quantityOf(id: number, monasteryId?: string) {
    return (
      lines.value.find(
        ({ book }) => book.id === id && book.monasteryId === monasteryId,
      )?.quantity ?? 0
    )
  }
  function add(book: BasketBook) {
    if (!validBook(book)) return false
    const { id, monasteryId } = book
    const quantity = quantityOf(id, monasteryId) + 1
    if (quantity > MAX_QUANTITY) return false
    const item = {
      book: snapshotBook(book),
      quantity,
    }
    save(
      quantity === 1
        ? [...lines.value, item]
        : lines.value.map((line) =>
            line.book.id === id && line.book.monasteryId === monasteryId
              ? item
              : line,
          ),
    )
    announcement.value = `${book.name}: ${quantity} in basket.`
    return true
  }
  function setQuantity(id: number, quantity: number, monasteryId?: string) {
    const current = quantityOf(id, monasteryId)
    if (quantity === 0 && current) {
      remove(id, monasteryId)
      return true
    }
    if (!validQuantity(quantity) || !current) return false
    save(
      lines.value.map((line) =>
        line.book.id === id && line.book.monasteryId === monasteryId
          ? { ...line, quantity }
          : line,
      ),
    )
    announcement.value = `Quantity updated to ${quantity}.`
    return true
  }
  function remove(id: number, monasteryId?: string) {
    save(
      lines.value.filter(
        ({ book }) => book.id !== id || book.monasteryId !== monasteryId,
      ),
    )
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
