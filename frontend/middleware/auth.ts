import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  // Si no hay sesión, redirige a login
  if (!authStore.isAuthenticated) {
    const redirect = encodeURIComponent(to.fullPath || '/')
    return navigateTo(`/login?redirect=${redirect}`)
  }

  // Rutas bajo /admin requieren rol admin
  if (to.path.startsWith('/admin') && !authStore.isAdmin) {
    return navigateTo('/dashboard')
  }
})
