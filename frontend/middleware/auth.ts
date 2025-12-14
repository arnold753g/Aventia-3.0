import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  // Si no hay sesión, redirige a login
  if (!authStore.isAuthenticated) {
    return navigateTo('/login')
  }

  // Rutas bajo /admin requieren rol admin
  if (to.path.startsWith('/admin') && !authStore.isAdmin) {
    return navigateTo('/dashboard')
  }
})
