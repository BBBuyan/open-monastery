<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

const props = defineProps<{ bookId: number; name: string; text: string }>()
const expanded = ref(false)
const overflows = ref(false)
const description = ref<HTMLElement>()
let observer: ResizeObserver | undefined
function measure() {
  const element = description.value
  if (!element) return
  const lineHeight = Number.parseFloat(getComputedStyle(element).lineHeight)
  overflows.value = element.scrollHeight > lineHeight * 3 + 1
}
watch(
  () => props.text,
  async () => {
    expanded.value = false
    await nextTick()
    measure()
  },
)
onMounted(() => {
  measure()
  if (typeof ResizeObserver !== 'undefined') {
    observer = new ResizeObserver(measure)
    if (description.value) observer.observe(description.value)
  } else window.addEventListener('resize', measure)
})
onUnmounted(() => {
  observer?.disconnect()
  window.removeEventListener('resize', measure)
})
</script>

<template>
  <div>
    <p
      :id="`book-description-${bookId}`"
      ref="description"
      class="description"
      :class="{ 'description-collapsed': !expanded }"
      lang="mn"
    >
      {{ text }}
    </p>
    <button
      v-if="overflows"
      class="description-toggle"
      type="button"
      :aria-expanded="expanded"
      :aria-controls="`book-description-${bookId}`"
      :aria-label="`${expanded ? 'See less' : 'See more'} of the description for ${name}`"
      @click="expanded = !expanded"
    >
      {{ expanded ? 'See less' : 'See more' }}
    </button>
  </div>
</template>
