import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  if (!authStore.isAuthenticated) {
    const redirect = encodeURIComponent(to.fullPath || '/')
    return navigateTo(`/login?redirect=${redirect}`)
  }

  if (!authStore.isEncargado) {
    return navigateTo(authStore.isAdmin ? '/admin/dashboard' : '/dashboard')
  }
})

