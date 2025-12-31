import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  if (!authStore.isAuthenticated) {
    const redirect = encodeURIComponent(to.fullPath || '/')
    return navigateTo(`/login?redirect=${redirect}`)
  }

  if (!authStore.isTurista) {
    if (authStore.isAdmin) return navigateTo('/admin/dashboard')
    if (authStore.isEncargado) return navigateTo('/agencia/dashboard')
    return navigateTo('/dashboard')
  }
})

