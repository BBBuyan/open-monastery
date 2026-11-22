<script setup lang="ts">
import { computed } from 'vue'
import type { Category } from '../api'
import SelectDropdown from './SelectDropdown.vue'

const props = defineProps<{
  modelValue?: number
  categories: Category[]
  loading?: boolean
}>()
const emit = defineEmits<{ 'update:modelValue': [value: number | undefined] }>()
const options = computed(() => [
  { id: undefined, name: 'All categories' },
  ...props.categories,
])
</script>

<template>
  <SelectDropdown
    id="book-category"
    label-id="book-category-label"
    :model-value="modelValue"
    :options="options"
    :loading="loading"
    :disabled="!categories.length"
    placeholder="All categories"
    loading-label="Loading categories…"
    @update:model-value="emit('update:modelValue', $event)"
  />
</template>
