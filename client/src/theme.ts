import { onMounted, onUnmounted, ref } from 'vue'

type Theme = 'light' | 'dark'
const THEME_KEY = 'open-monastery:theme'
function validTheme(value: unknown): value is Theme {
  return value === 'light' || value === 'dark'
}

export function useTheme() {
  const system = window.matchMedia('(prefers-color-scheme: dark)')
  let preference: Theme | null = null
  try {
    const saved = localStorage.getItem(THEME_KEY)
    if (validTheme(saved)) preference = saved
  } catch {
    /* Theme controls also work when storage is unavailable. */
  }
  const theme = ref<Theme>(preference ?? (system.matches ? 'dark' : 'light'))
  function apply(next: Theme) {
    theme.value = next
    document.documentElement.dataset.theme = next
    document
      .querySelector('meta[name="theme-color"]')
      ?.setAttribute('content', next === 'dark' ? '#211e19' : '#f6f1e7')
  }
  function toggle() {
    preference = theme.value === 'dark' ? 'light' : 'dark'
    apply(preference)
    try {
      localStorage.setItem(THEME_KEY, preference)
    } catch {
      /* Keep the choice for this visit. */
    }
  }
  function syncSystem() {
    if (!preference) apply(system.matches ? 'dark' : 'light')
  }
  function syncStorage(event: StorageEvent) {
    if (event.key !== THEME_KEY && event.key !== null) return
    preference = validTheme(event.newValue) ? event.newValue : null
    apply(preference ?? (system.matches ? 'dark' : 'light'))
  }
  apply(theme.value)
  onMounted(() => {
    system.addEventListener('change', syncSystem)
    window.addEventListener('storage', syncStorage)
  })
  onUnmounted(() => {
    system.removeEventListener('change', syncSystem)
    window.removeEventListener('storage', syncStorage)
  })
  return { theme, toggle }
}
