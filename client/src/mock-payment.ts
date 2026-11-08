import { ref } from 'vue'
export type PaymentMethod = 'qpay' | 'monpay'
export function createMockPayment() {
  const method = ref<PaymentMethod | null>(null)
  const stage = ref<'selection' | 'preview' | 'complete'>('selection')
  function choose(value: PaymentMethod) {
    if (value !== 'qpay' && value !== 'monpay') return
    method.value = value
    stage.value = 'selection'
  }
  function preview() {
    if (!method.value) return false
    stage.value = 'preview'
    return true
  }
  function complete() {
    if (stage.value !== 'preview') return false
    stage.value = 'complete'
    return true
  }
  function reset() {
    stage.value = 'selection'
  }
  return { method, stage, choose, preview, complete, reset }
}
