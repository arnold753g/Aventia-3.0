import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  if (!authStore.isAuthenticated) {
    const redirect = encodeURIComponent(to.fullPath || '/')
    return navigateTo(`/login?redirect=${redirect}`)
  }

  if (!authStore.isEncargado) {
    if (authStore.isAdmin) return navigateTo('/admin/dashboard')
    if (authStore.isTurista) return navigateTo('/turista/dashboard')
    return navigateTo('/dashboard')
  }
})
