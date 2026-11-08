import assert from 'node:assert/strict'
import { test } from 'node:test'
import { createBasket, BASKET_KEY, MAX_QUANTITY } from '../src/basket.ts'

const book = (id, price = 10000) => ({
  id,
  name: `Book ${id}`,
  description: 'Description',
  price,
})
function memoryStorage(initial = null) {
  let value = initial
  return {
    getItem: () => value,
    setItem: (key, next) => {
      assert.equal(key, BASKET_KEY)
      value = next
    },
  }
}

test('adding books merges matching IDs and computes quantities and MNT totals', () => {
  const basket = createBasket()
  basket.add(book(1))
  basket.add(book(2, 2500))
  basket.add(book(1))
  assert.equal(basket.items.value.length, 2)
  assert.equal(basket.quantityOf(1), 2)
  assert.equal(basket.count.value, 3)
  assert.equal(basket.total.value, 22500)
  basket.setQuantity(2, 3)
  assert.equal(basket.total.value, 27500)
  basket.remove(1)
  assert.equal(basket.count.value, 3)
  assert.equal(basket.total.value, 7500)
})

test('basket survives recreation and removing the final item persists an empty basket', () => {
  const storage = memoryStorage()
  const first = createBasket(storage)
  first.add(book(1))
  first.setQuantity(1, 4)
  const second = createBasket(storage)
  assert.equal(second.quantityOf(1), 4)
  assert.equal(second.total.value, 40000)
  second.remove(1)
  assert.deepEqual(createBasket(storage).items.value, [])
})

test('invalid and excessive quantities never change a basket', () => {
  const basket = createBasket()
  basket.add(book(1))
  for (const quantity of [-1, 1.5, NaN, Infinity, MAX_QUANTITY + 1]) {
    assert.equal(basket.setQuantity(1, quantity), false)
  }
  assert.equal(basket.quantityOf(1), 1)
  assert.equal(basket.setQuantity(999, 2), false)
  basket.setQuantity(1, MAX_QUANTITY)
  assert.equal(basket.add(book(1)), false)
  assert.equal(basket.quantityOf(1), MAX_QUANTITY)
})

test('zero prices are valid, invalid prices and IDs are rejected', () => {
  const basket = createBasket()
  assert.equal(basket.add(book(1, 0)), true)
  for (const invalid of [
    book(0),
    book(2, -1),
    book(2, NaN),
    book(2, Infinity),
    book(2, 1.5),
  ]) {
    assert.equal(basket.add(invalid), false)
  }
  assert.equal(basket.count.value, 1)
  assert.equal(basket.total.value, 0)
})

test('adding a known book refreshes its name and price without duplicating it', () => {
  const basket = createBasket()
  basket.add(book(1))
  basket.add({ ...book(1, 12000), name: 'Updated' })
  assert.equal(basket.items.value.length, 1)
  assert.equal(basket.items.value[0].book.name, 'Updated')
  assert.equal(basket.total.value, 24000)
})

test('corrupt saved data is ignored; invalid lines are dropped and duplicates are merged', () => {
  for (const saved of [
    'broken json',
    'null',
    '{}',
    '{"version":2,"items":[]}',
  ]) {
    assert.deepEqual(createBasket(memoryStorage(saved)).items.value, [])
  }
  const storage = memoryStorage(
    JSON.stringify({
      version: 1,
      items: [
        { book: book(1), quantity: 2 },
        { book: book(1), quantity: 3 },
        { book: book(2), quantity: -1 },
        { book: book(0), quantity: 1 },
        null,
        { book: { ...book(3), price: '1000' }, quantity: 1 },
      ],
    }),
  )
  const basket = createBasket(storage)
  assert.equal(basket.items.value.length, 1)
  assert.equal(basket.quantityOf(1), 5)
})

test('unavailable storage does not prevent adding or editing books', () => {
  const storage = {
    getItem() {
      throw new Error('Blocked')
    },
    setItem() {
      throw new Error('Quota exceeded')
    },
  }
  const basket = createBasket(storage)
  basket.add(book(1))
  basket.setQuantity(1, 2)
  assert.equal(basket.quantityOf(1), 2)
  assert.equal(basket.total.value, 20000)
  assert.equal(basket.persistent.value, false)
  basket.remove(1)
  assert.equal(basket.count.value, 0)
})

test('a read failure is reported immediately and successful writes restore persistence', () => {
  let blocked = true
  const basket = createBasket({
    getItem() {
      throw new Error('Blocked')
    },
    setItem() {
      if (blocked) throw new Error('Full')
    },
  })
  assert.equal(basket.persistent.value, false)
  basket.add(book(1))
  assert.equal(basket.persistent.value, false)
  blocked = false
  basket.setQuantity(1, 2)
  assert.equal(basket.persistent.value, true)
})

test('reducing a quantity to zero removes the book and persists the updated total', () => {
  const storage = memoryStorage()
  const basket = createBasket(storage)
  basket.add(book(1))
  basket.add(book(2, 2500))
  assert.equal(basket.setQuantity(1, 0), true)
  assert.equal(basket.quantityOf(1), 0)
  assert.equal(basket.count.value, 1)
  assert.equal(basket.total.value, 2500)
  assert.equal(createBasket(storage).quantityOf(1), 0)
  assert.equal(basket.setQuantity(2, 0), true)
  assert.equal(basket.items.value.length, 0)
  assert.equal(basket.total.value, 0)
  assert.equal(basket.setQuantity(999, 0), false)
})

test('storage contains only ID, quantity, name, and price and restores without descriptions', () => {
  const storage = memoryStorage()
  const basket = createBasket(storage)
  basket.add(book(1))
  basket.setQuantity(1, 3)
  assert.deepEqual(JSON.parse(storage.getItem(BASKET_KEY)), {
    version: 2,
    items: [{ id: 1, quantity: 3, name: 'Book 1', price: 10000 }],
  })
  const restored = createBasket(storage)
  assert.equal(restored.quantityOf(1), 3)
  assert.equal(restored.total.value, 30000)
  assert.deepEqual(restored.items.value[0].book, {
    id: 1,
    name: 'Book 1',
    price: 10000,
  })
})

test('legacy baskets migrate immediately, dropping descriptions without losing quantities', () => {
  const storage = memoryStorage(
    JSON.stringify({
      version: 1,
      items: [
        { book: book(1), quantity: 2 },
        { book: book(2, 2500), quantity: 3 },
      ],
    }),
  )
  const basket = createBasket(storage)
  assert.equal(basket.count.value, 5)
  assert.equal(basket.total.value, 27500)
  assert.deepEqual(JSON.parse(storage.getItem(BASKET_KEY)), {
    version: 2,
    items: [
      { id: 1, quantity: 2, name: 'Book 1', price: 10000 },
      { id: 2, quantity: 3, name: 'Book 2', price: 2500 },
    ],
  })
})
