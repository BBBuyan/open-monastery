import assert from 'node:assert/strict'
import { test } from 'node:test'
import { createMockPayment } from '../src/mock-payment.ts'

test('a demo requires a payment method and preview before completion', () => {
  const payment = createMockPayment()
  assert.equal(payment.method.value, null)
  assert.equal(payment.preview(), false)
  assert.equal(payment.complete(), false)
  assert.equal(payment.stage.value, 'selection')
  payment.choose('qpay')
  assert.equal(payment.preview(), true)
  assert.equal(payment.stage.value, 'preview')
  assert.equal(payment.complete(), true)
  assert.equal(payment.stage.value, 'complete')
})

test('switching providers resets the demo and supports both payment methods', () => {
  const payment = createMockPayment()
  payment.choose('qpay')
  payment.preview()
  payment.complete()
  payment.choose('monpay')
  assert.equal(payment.method.value, 'monpay')
  assert.equal(payment.stage.value, 'selection')
  assert.equal(payment.complete(), false)
  assert.equal(payment.preview(), true)
  assert.equal(payment.complete(), true)
})

test('basket changes can invalidate a demo without forgetting the selected provider', () => {
  const payment = createMockPayment()
  payment.choose('monpay')
  payment.preview()
  payment.complete()
  payment.reset()
  assert.equal(payment.stage.value, 'selection')
  assert.equal(payment.method.value, 'monpay')
  assert.equal(payment.complete(), false)
})

test('unsupported provider IDs are ignored', () => {
  const payment = createMockPayment()
  payment.choose('unsupported')
  assert.equal(payment.method.value, null)
  assert.equal(payment.preview(), false)
})
