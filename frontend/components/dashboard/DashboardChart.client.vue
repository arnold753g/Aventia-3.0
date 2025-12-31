<template>
  <VueChart :type="type" :data="data" :options="options" />
</template>

<script setup lang="ts">
import { Chart, registerables } from 'chart.js'
import type { ChartData, ChartOptions, ChartType } from 'chart.js'
import { Chart as VueChart } from 'vue-chartjs'

Chart.register(...registerables)

defineProps<{
  type: ChartType
  data: ChartData
  options?: ChartOptions
}>()
</script>
