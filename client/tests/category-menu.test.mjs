import assert from 'node:assert/strict'
import { test } from 'node:test'
import { categoryMenuPosition } from '../src/category-menu.ts'

const viewport = { width: 1000, height: 800 }
const trigger = (top) => ({ left: 700, top, bottom: top + 58, width: 240 })
test('a dropdown near the top opens below with six rows of space', () => {
  const result = categoryMenuPosition(trigger(80), viewport)
  assert.equal(result.opensAbove, false)
  assert.equal(result.maxHeight, 274)
  assert.ok(result.top + result.maxHeight <= viewport.height - 8)
})
test('a dropdown in the lower third opens above instead of spilling below', () => {
  const result = categoryMenuPosition(trigger(540), viewport)
  assert.equal(result.opensAbove, true)
  assert.equal(result.maxHeight, 274)
  assert.ok(result.top - result.maxHeight >= 8)
})
test('a dropdown near the bottom opens above within the viewport', () => {
  const result = categoryMenuPosition(trigger(730), viewport)
  assert.equal(result.opensAbove, true)
  assert.ok(result.top < 730)
  assert.ok(result.top - result.maxHeight >= 8)
})
test('small windows reduce the menu height to the available room', () => {
  const result = categoryMenuPosition(trigger(120), { width: 320, height: 300 })
  assert.ok(result.maxHeight < 274)
  assert.ok(result.top + result.maxHeight <= 292)
  assert.ok(result.left + result.width <= 312)
})
test('a menu above the trigger also shrinks when the viewport is short', () => {
  const result = categoryMenuPosition(trigger(180), { width: 320, height: 300 })
  assert.equal(result.opensAbove, true)
  assert.ok(result.maxHeight < 274)
  assert.ok(result.top - result.maxHeight >= 8)
})
