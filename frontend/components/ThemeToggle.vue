<template>
  <button
    type="button"
    class="panel-theme-toggle"
    :aria-label="label"
    :title="label"
    @click="toggleMode"
  >
    <i :class="icon"></i>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const { isDark, toggleMode } = useThemeMode()

const icon = computed(() => (isDark.value ? 'pi pi-sun' : 'pi pi-moon'))
const label = computed(() => (isDark.value ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'))
</script>
