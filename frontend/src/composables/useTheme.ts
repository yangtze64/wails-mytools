import { ref, watch } from 'vue'

type ThemeMode = 'dark' | 'light'

const STORAGE_KEY = 'mytools.theme'

const currentTheme = ref<ThemeMode>(readStoredTheme())

function readStoredTheme(): ThemeMode {
  if (typeof window === 'undefined') {
    return 'dark'
  }
  try {
    const stored = window.localStorage.getItem(STORAGE_KEY)
    if (stored === 'light' || stored === 'dark') {
      return stored
    }
  } catch { /* ignore */ }
  return 'dark'
}

function applyTheme(mode: ThemeMode) {
  const root = document.documentElement
  root.setAttribute('data-theme', mode)
  root.style.colorScheme = mode
}

export function useTheme() {
  function toggleTheme() {
    currentTheme.value = currentTheme.value === 'dark' ? 'light' : 'dark'
  }

  function setTheme(mode: ThemeMode) {
    currentTheme.value = mode
  }

  watch(currentTheme, (mode) => {
    applyTheme(mode)
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(STORAGE_KEY, mode)
    }
  }, { immediate: true })

  return {
    currentTheme,
    isDark: () => currentTheme.value === 'dark',
    toggleTheme,
    setTheme,
  }
}
