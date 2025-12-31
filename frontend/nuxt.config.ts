// https://nuxt.com/docs/api/configuration/nuxt-config
const apiBase = process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:5750/api/v1'
const wsBaseUrl = process.env.NUXT_PUBLIC_WS_BASE || 'ws://localhost:5750'
const isProd = process.env.NODE_ENV === 'production'

export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: !isProd },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],

  css: [
    'primeicons/primeicons.css',
    '~/assets/theme.css'
  ],

  build: {
    transpile: ['primevue']
  },

  runtimeConfig: {
    public: {
      apiBase,
      wsBaseUrl
    }
  },

  app: {
    head: {
      title: 'Sistema Andaria - Gestión Turística',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content: 'Plataforma de gestión turística Andaria'
        }
      ]
    }
  }
})
