import { useAuthStore } from '~/stores/auth'

export const useAuthAPI = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => {
    if (!authStore.token) return {}
    return { Authorization: `Bearer ${authStore.token}` }
  }

  const register = async (data: any) => {
    return $fetch(`${baseURL}/auth/register`, {
      method: 'POST',
      body: data
    })
  }

  const verifyEmail = async (email: string, code: string) => {
    return $fetch(`${baseURL}/auth/verify-email`, {
      method: 'POST',
      body: { email, code }
    })
  }

  const resendCode = async (email: string) => {
    return $fetch(`${baseURL}/auth/resend-email-code`, {
      method: 'POST',
      body: { email }
    })
  }

  const forgotPassword = async (email: string) => {
    return $fetch(`${baseURL}/auth/forgot-password`, {
      method: 'POST',
      body: { email }
    })
  }

  const resetPassword = async (email: string, code: string, newPassword: string) => {
    return $fetch(`${baseURL}/auth/reset-password`, {
      method: 'POST',
      body: { email, code, new_password: newPassword }
    })
  }

  const changePassword = async (currentPassword: string, newPassword: string) => {
    return $fetch(`${baseURL}/auth/change-password`, {
      method: 'POST',
      headers: authHeader(),
      body: { current_password: currentPassword, new_password: newPassword }
    })
  }

  const setInitialPassword = async (email: string, newPassword: string) => {
    return $fetch(`${baseURL}/auth/set-initial-password`, {
      method: 'POST',
      body: { email, new_password: newPassword }
    })
  }

  return {
    register,
    verifyEmail,
    resendCode,
    forgotPassword,
    resetPassword,
    changePassword,
    setInitialPassword
  }
}
