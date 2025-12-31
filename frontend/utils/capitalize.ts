/**
 * Capitaliza la primera letra de cada palabra
 */
export const capitalizeWords = (text: string): string => {
  if (!text) return ''

  return text
    .toLowerCase()
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

/**
 * Capitaliza solo la primera letra del texto completo
 */
export const capitalizeFirst = (text: string): string => {
  if (!text) return ''
  return text.charAt(0).toUpperCase() + text.slice(1).toLowerCase()
}

/**
 * Capitaliza para nombres propios (mantiene mayúsculas intermedias)
 */
export const capitalizeProperName = (text: string): string => {
  if (!text) return ''

  const exceptions = ['de', 'del', 'la', 'las', 'los', 'el', 'y']

  return text
    .toLowerCase()
    .split(' ')
    .map((word, index) => {
      // Primera palabra siempre con mayúscula
      if (index === 0) {
        return word.charAt(0).toUpperCase() + word.slice(1)
      }
      // Excepciones en minúscula
      if (exceptions.includes(word)) {
        return word
      }
      // Resto con mayúscula inicial
      return word.charAt(0).toUpperCase() + word.slice(1)
    })
    .join(' ')
}
