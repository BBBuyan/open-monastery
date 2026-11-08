<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import type { Category } from '../api'
import { categoryMenuPosition } from '../category-menu'

const props = defineProps<{
  modelValue?: number
  categories: Category[]
  loading?: boolean
}>()
const emit = defineEmits<{ 'update:modelValue': [value: number | undefined] }>()
const root = ref<HTMLElement>()
const trigger = ref<HTMLButtonElement>()
const menu = ref<HTMLElement>()
const open = ref(false)
const active = ref(0)
const pointerFocus = ref(false)
const menuStyle = ref<Record<string, string>>({})
const options = computed(() => [
  { id: undefined as number | undefined, name: 'All categories' },
  ...props.categories,
])
const selectedIndex = computed(() =>
  Math.max(
    0,
    options.value.findIndex((option) => option.id === props.modelValue),
  ),
)
const disabled = computed(() => props.loading || !props.categories.length)
const label = computed(() =>
  props.loading
    ? 'Loading categories…'
    : (options.value[selectedIndex.value]?.name ?? 'All categories'),
)
let typed = ''
let typeTimer: ReturnType<typeof setTimeout> | undefined
function close() {
  open.value = false
  typed = ''
  clearTimeout(typeTimer)
}
function position() {
  if (!trigger.value || !open.value) return
  const rect = trigger.value.getBoundingClientRect()
  if (rect.bottom <= 0 || rect.top >= window.innerHeight) {
    close()
    return
  }
  const placement = categoryMenuPosition(rect, {
    width: window.innerWidth,
    height: window.innerHeight,
  })
  menuStyle.value = {
    left: `${placement.left}px`,
    top: `${placement.top}px`,
    width: `${placement.width}px`,
    maxHeight: `${placement.maxHeight}px`,
    transform: placement.opensAbove ? 'translateY(-100%)' : 'none',
  }
}
async function revealActive() {
  await nextTick()
  const container = menu.value
  const option = document.getElementById(`category-option-${active.value}`)
  if (!container || !option) return
  const bounds = option.getBoundingClientRect()
  const visible = container.getBoundingClientRect()
  if (bounds.top < visible.top + 4)
    container.scrollTop -= visible.top + 4 - bounds.top
  else if (bounds.bottom > visible.bottom - 4)
    container.scrollTop += bounds.bottom - visible.bottom + 4
}
async function show(index = selectedIndex.value) {
  if (disabled.value) return
  active.value = index
  open.value = true
  position()
  await revealActive()
}
function commit(index = active.value) {
  const option = options.value[index]
  if (!option) return
  emit('update:modelValue', option.id)
  close()
  trigger.value?.focus({ preventScroll: true })
}
function onKey(event: KeyboardEvent) {
  pointerFocus.value = false
  if (event.key === 'Tab') {
    if (open.value) commit()
    return
  }
  if (event.key === 'Escape') {
    if (open.value) {
      event.preventDefault()
      event.stopPropagation()
      close()
    }
    return
  }
  if (['Enter', ' '].includes(event.key)) {
    event.preventDefault()
    if (open.value) commit()
    else void show()
    return
  }
  if (
    ['ArrowDown', 'ArrowUp', 'Home', 'End', 'PageDown', 'PageUp'].includes(
      event.key,
    )
  ) {
    event.preventDefault()
    const last = options.value.length - 1
    let index = open.value ? active.value : selectedIndex.value
    if (event.key === 'Home') index = 0
    else if (event.key === 'End') index = last
    else if (open.value)
      index = Math.max(
        0,
        Math.min(
          last,
          index +
            (event.key === 'ArrowDown'
              ? 1
              : event.key === 'ArrowUp'
                ? -1
                : event.key === 'PageDown'
                  ? 6
                  : -6),
        ),
      )
    void show(index)
    return
  }
  if (
    event.key.length === 1 &&
    !event.ctrlKey &&
    !event.metaKey &&
    !event.altKey
  ) {
    event.preventDefault()
    clearTimeout(typeTimer)
    typed += event.key.toLocaleLowerCase('mn-MN')
    const repeated = [...typed].every((character) => character === typed[0])
    const search = repeated ? typed[0]! : typed
    const start = repeated ? active.value + 1 : 0
    for (let offset = 0; offset < options.value.length; offset++) {
      const index = (start + offset) % options.value.length
      if (
        options.value[index]!.name.toLocaleLowerCase('mn-MN').startsWith(search)
      ) {
        void show(index)
        break
      }
    }
    typeTimer = setTimeout(() => {
      typed = ''
    }, 700)
  }
}
function outside(event: Event) {
  const target = event.target
  if (
    target instanceof Node &&
    !root.value?.contains(target) &&
    !menu.value?.contains(target)
  )
    close()
}
function viewportChanged(event: Event) {
  if (event.target instanceof Node && menu.value?.contains(event.target)) return
  position()
}
watch(disabled, (value) => {
  if (value) close()
})
watch(
  () => props.modelValue,
  () => {
    active.value = selectedIndex.value
  },
)
onMounted(() => {
  document.addEventListener('pointerdown', outside, true)
  document.addEventListener('focusin', outside)
  document.addEventListener('scroll', viewportChanged, true)
  window.addEventListener('resize', viewportChanged)
})
onUnmounted(() => {
  clearTimeout(typeTimer)
  document.removeEventListener('pointerdown', outside, true)
  document.removeEventListener('focusin', outside)
  document.removeEventListener('scroll', viewportChanged, true)
  window.removeEventListener('resize', viewportChanged)
})
</script>

<template>
  <div ref="root" class="category-picker">
    <button
      id="book-category"
      ref="trigger"
      class="category-trigger"
      :class="{ 'pointer-focus': pointerFocus }"
      type="button"
      role="combobox"
      aria-haspopup="listbox"
      aria-labelledby="book-category-label book-category-value"
      aria-controls="category-options"
      :aria-expanded="open"
      :aria-activedescendant="open ? `category-option-${active}` : undefined"
      :disabled="disabled"
      @pointerdown="pointerFocus = true"
      @keydown="onKey"
      @click="open ? close() : show()"
    >
      <span
        id="book-category-value"
        :lang="modelValue === undefined ? undefined : 'mn'"
        >{{ label }}</span
      >
      <svg
        viewBox="0 0 20 20"
        width="18"
        height="18"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        aria-hidden="true"
      >
        <path :d="open ? 'm5 12 5-5 5 5' : 'm5 8 5 5 5-5'" />
      </svg>
    </button>
    <Teleport to="body">
      <ul
        v-if="open"
        id="category-options"
        ref="menu"
        class="category-menu"
        :style="menuStyle"
        role="listbox"
        aria-labelledby="book-category-label"
      >
        <li
          v-for="(option, index) in options"
          :id="`category-option-${index}`"
          :key="option.id ?? 'all'"
          role="option"
          :aria-selected="option.id === modelValue"
          :class="{ 'is-active': active === index }"
          @mousedown.prevent
          @pointermove="active = index"
          @click="commit(index)"
        >
          <span :lang="option.id === undefined ? undefined : 'mn'">{{
            option.name
          }}</span>
          <span v-if="option.id === modelValue" aria-hidden="true">✓</span>
        </li>
      </ul>
    </Teleport>
  </div>
</template>
