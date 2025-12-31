import { computed, watch } from 'vue'

type ThemeMode = 'light' | 'dark' | 'system'

const STORAGE_KEY = 'andaria-theme-mode'

export const useThemeMode = () => {
  const mode = useState<ThemeMode>('theme-mode', () => 'light')
  const prefersDark = useState<boolean>('theme-prefers-dark', () => false)
  const initialized = useState<boolean>('theme-initialized', () => false)

  const isDark = computed(() => {
    if (mode.value === 'dark') return true
    if (mode.value === 'light') return false
    return prefersDark.value
  })

  const applyThemeClass = (value: boolean) => {
    if (!process.client) return
    const root = document.documentElement
    root.classList.toggle('dark', value)
  }

  if (process.client && !initialized.value) {
    const stored = localStorage.getItem(STORAGE_KEY) as ThemeMode | null
    if (stored === 'light' || stored === 'dark' || stored === 'system') {
      mode.value = stored
    }

    const media = window.matchMedia('(prefers-color-scheme: dark)')
    prefersDark.value = media.matches
    media.addEventListener('change', (event) => {
      prefersDark.value = event.matches
    })

    watch(isDark, (value) => {
      applyThemeClass(value)
      localStorage.setItem(STORAGE_KEY, mode.value)
    }, { immediate: true })

    initialized.value = true
  }

  const toggleMode = () => {
    mode.value = isDark.value ? 'light' : 'dark'
  }

  const setMode = (value: ThemeMode) => {
    mode.value = value
  }

  return {
    mode,
    isDark,
    toggleMode,
    setMode
  }
}
