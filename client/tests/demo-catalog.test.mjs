import assert from 'node:assert/strict'
import { test } from 'node:test'
import { createBasket } from '../src/basket.ts'
import {
  demoMonasteries,
  getDemoBooks,
  monasteryCatalogLink,
  monasteryIdFromHash,
} from '../src/demo-catalog.ts'

test('monastery catalog links round-trip through the URL fragment', () => {
  for (const { id } of demoMonasteries) {
    assert.equal(monasteryIdFromHash(monasteryCatalogLink(id)), id)
  }
  assert.equal(monasteryIdFromHash('#readings'), undefined)
  assert.equal(monasteryIdFromHash('#basket?monastery=amber-valley'), undefined)
})

test('demo catalogs differ by monastery and search stays within the chosen catalog', async () => {
  const [amber, cedar] = demoMonasteries
  const amberPage = await getDemoBooks(amber, { query: '' })
  const cedarPage = await getDemoBooks(cedar, { query: '' })
  assert.notDeepEqual(
    amberPage.data.map(({ id }) => id),
    cedarPage.data.map(({ id }) => id),
  )
  const [sharedPrayer] = amberPage.data
  const matching = await getDemoBooks(cedar, {
    categoryId: 1,
    query: sharedPrayer.name,
  })
  assert.equal(matching.data.length, 1)
  assert.equal(matching.data[0].id, sharedPrayer.id)
  assert.notEqual(matching.data[0].price, sharedPrayer.price)
})

test('the same prayer from two monasteries keeps separate quantities and prices after reload', async () => {
  const [amber, cedar] = demoMonasteries
  const {
    data: [amberOffer],
  } = await getDemoBooks(amber, { query: 'Авидын' })
  const {
    data: [cedarOffer],
  } = await getDemoBooks(cedar, { query: 'Авидын' })
  let stored = null
  const storage = {
    getItem: () => stored,
    setItem: (_, value) => {
      stored = value
    },
  }
  const basket = createBasket(storage)
  basket.add(amberOffer)
  basket.add(cedarOffer)
  basket.setQuantity(amberOffer.id, 2, amber.id)

  const restored = createBasket(storage)
  assert.equal(restored.items.value.length, 2)
  assert.equal(restored.quantityOf(amberOffer.id, amber.id), 2)
  assert.equal(restored.quantityOf(cedarOffer.id, cedar.id), 1)
  assert.equal(restored.total.value, amberOffer.price * 2 + cedarOffer.price)
  restored.remove(cedarOffer.id, cedar.id)
  assert.equal(restored.total.value, amberOffer.price * 2)
})
