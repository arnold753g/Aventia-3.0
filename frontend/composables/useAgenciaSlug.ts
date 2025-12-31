export const useAgenciaSlug = () => {
  const config = useRuntimeConfig()

  /**
   * Genera un slug URL-friendly desde un texto
   */
  const generateSlug = (text: string): string => {
    return text
      .toLowerCase()
      .normalize('NFD') // Descomponer caracteres con acentos
      .replace(/[\u0300-\u036f]/g, '') // Eliminar diacríticos
      .replace(/[^a-z0-9]+/g, '-') // Reemplazar no-alfanuméricos con guiones
      .replace(/^-+|-+$/g, '') // Eliminar guiones al inicio y final
      .substring(0, 200) // Limitar longitud
  }

  /**
   * Obtiene la URL pública completa de una agencia
   */
  const getPublicUrl = (slug: string): string => {
    if (process.client) {
      const origin = window.location.origin
      return `${origin}/agencias/${slug}`
    }
    // En el servidor, usar la URL base del config
    return `/agencias/${slug}`
  }

  /**
   * Obtiene la URL completa absoluta de una agencia (para compartir)
   */
  const getShareUrl = (slug: string): string => {
    if (process.client) {
      const origin = window.location.origin
      return `${origin}/agencias/${slug}`
    }
    // Fallback
    return `https://andaria.com/agencias/${slug}`
  }

  /**
   * Copia la URL pública al portapapeles
   */
  const copyToClipboard = async (slug: string): Promise<boolean> => {
    if (!process.client) return false

    try {
      const url = getShareUrl(slug)
      await navigator.clipboard.writeText(url)
      return true
    } catch (error) {
      // Fallback para navegadores antiguos
      try {
        const textArea = document.createElement('textarea')
        textArea.value = getShareUrl(slug)
        textArea.style.position = 'fixed'
        textArea.style.left = '-999999px'
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        const successful = document.execCommand('copy')
        document.body.removeChild(textArea)
        return successful
      } catch (err) {
        console.error('Error al copiar:', err)
        return false
      }
    }
  }

  return {
    generateSlug,
    getPublicUrl,
    getShareUrl,
    copyToClipboard
  }
}
