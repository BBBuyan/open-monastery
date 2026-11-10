<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

defineProps<{ theme: 'light' | 'dark'; active: string }>()
defineEmits<{ toggleTheme: [] }>()
const disclosure = ref<HTMLDetailsElement>()
const trigger = ref<HTMLElement>()
function close() {
  if (disclosure.value) disclosure.value.open = false
}
function escape(event: KeyboardEvent) {
  if (!disclosure.value?.open) return
  event.preventDefault()
  close()
  trigger.value?.focus()
}
function focusout(event: FocusEvent) {
  if (!disclosure.value?.contains(event.relatedTarget as Node | null)) close()
}
function outside(event: PointerEvent) {
  if (!disclosure.value?.contains(event.target as Node)) close()
}
onMounted(() => document.addEventListener('pointerdown', outside))
onUnmounted(() => document.removeEventListener('pointerdown', outside))
</script>

<template>
  <details
    ref="disclosure"
    class="profile-menu"
    @keydown.esc="escape"
    @focusout="focusout"
  >
    <summary
      ref="trigger"
      class="profile-link"
      aria-label="Account options"
      title="Account options"
    >
      <svg
        viewBox="0 0 32 32"
        width="32"
        height="32"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        aria-hidden="true"
      >
        <circle cx="16" cy="11" r="5" />
        <path d="M6 28v-3a10 10 0 0 1 20 0v3" />
      </svg>
    </summary>
    <div class="profile-menu-panel">
      <nav aria-label="Account">
        <a
          href="#profile"
          :aria-current="active === 'profile' ? 'page' : undefined"
          @click="close"
          >Profile Settings</a
        >
        <a
          href="#requests"
          :aria-current="active === 'requests' ? 'page' : undefined"
          @click="close"
          >Readings</a
        >
      </nav>
      <button
        type="button"
        :aria-pressed="theme === 'dark'"
        @click="$emit('toggleTheme')"
      >
        <span>Dark mode</span
        ><span
          class="theme-switch"
          :class="{ 'is-on': theme === 'dark' }"
          aria-hidden="true"
          ><span></span
        ></span>
      </button>
    </div>
  </details>
</template>
