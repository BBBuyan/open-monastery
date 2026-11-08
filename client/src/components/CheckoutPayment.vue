<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { createMockPayment } from '../mock-payment'
import { formatPrice } from '../money'

const props = defineProps<{ total: number; selectionKey: string }>()
const { method, stage, choose, preview, complete, reset } = createMockPayment()
const methods = [
  { id: 'qpay' as const, name: 'QPay', logo: '/QPay_logo_blue.svg' },
  { id: 'monpay' as const, name: 'MonPay', logo: '/monpay_white_logo.svg' },
]
const selected = computed(() =>
  methods.find((option) => option.id === method.value),
)
const previewPanel = ref<HTMLElement>()
const completePanel = ref<HTMLElement>()
const selectPanel = ref<HTMLElement>()
watch(() => props.selectionKey, reset)
async function showPreview() {
  if (!preview()) return
  await nextTick()
  previewPanel.value?.focus({ preventScroll: true })
}
async function simulatePayment() {
  if (!complete()) return
  await nextTick()
  completePanel.value?.focus({ preventScroll: true })
}
async function backToSelection() {
  reset()
  await nextTick()
  selectPanel.value
    ?.querySelector<HTMLInputElement>('input:checked')
    ?.focus({ preventScroll: true })
}
</script>

<template>
  <div class="checkout-payment">
    <div class="payment-demo-label">Demo checkout</div>
    <p id="payment-demo-notice" class="checkout-notice">
      Try the payment flow. No money will be charged and no prayer request will
      be submitted.
    </p>
    <form
      v-if="stage === 'selection'"
      ref="selectPanel"
      @submit.prevent="showPreview"
    >
      <fieldset class="payment-methods" aria-describedby="payment-demo-notice">
        <legend>Payment method</legend>
        <label
          v-for="option in methods"
          :key="option.id"
          class="payment-option"
          :class="{ 'is-selected': method === option.id }"
        >
          <input
            type="radio"
            name="payment-method"
            :value="option.id"
            :checked="method === option.id"
            required
            @change="choose(option.id)"
          />
          <span class="payment-logo" :class="`payment-logo-${option.id}`"
            ><img :src="option.logo" alt=""
          /></span>
          <span>{{ option.name }}</span>
        </label>
      </fieldset>
      <button class="button primary" type="submit" :disabled="!method">
        {{
          selected
            ? `Preview ${selected.name} payment`
            : 'Select a payment method'
        }}
      </button>
    </form>
    <section
      v-else-if="stage === 'preview'"
      ref="previewPanel"
      class="payment-preview"
      tabindex="-1"
      aria-labelledby="payment-preview-title"
    >
      <h3 id="payment-preview-title">{{ selected?.name }} payment preview</h3>
      <div class="payment-preview-provider">
        <span class="payment-logo" :class="`payment-logo-${method}`"
          ><img :src="selected?.logo" alt=""
        /></span>
        <strong>{{ formatPrice(total) }}</strong>
      </div>
      <p class="checkout-notice">
        This is a demo screen. No payment invoice has been created.
      </p>
      <button class="button primary" type="button" @click="simulatePayment">
        Simulate successful payment
      </button>
      <button class="payment-back" type="button" @click="backToSelection">
        Change payment method
      </button>
    </section>
    <section
      v-else
      ref="completePanel"
      class="payment-complete"
      tabindex="-1"
      aria-labelledby="payment-complete-title"
    >
      <svg
        viewBox="0 0 24 24"
        width="32"
        height="32"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        aria-hidden="true"
      >
        <circle cx="12" cy="12" r="10" />
        <path d="m7 12 3 3 7-7" />
      </svg>
      <h3 id="payment-complete-title">Demo complete</h3>
      <p class="checkout-notice">
        The {{ selected?.name }} payment was simulated. No money was charged or
        prayer request submitted. Your basket is unchanged.
      </p>
      <button class="button secondary" type="button" @click="backToSelection">
        Try another payment method
      </button>
    </section>
  </div>
</template>
