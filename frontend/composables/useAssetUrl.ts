export const useAssetUrl = () => {
  const config = useRuntimeConfig()
  const assetsBase = config.public.apiBase.replace(/\/api\/v1\/?$/, '')

  const resolveAssetUrl = (path?: string | null) => {
    if (!path) return ''
    let normalized = path.replace(/\\/g, '/')
    if (/^https?:\/\//i.test(normalized)) return normalized
    const uploadsIndex = normalized.indexOf('uploads/')
    if (uploadsIndex > -1) normalized = normalized.slice(uploadsIndex)
    normalized = normalized.replace(/^\.\//, '')
    const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
    return `${assetsBase}/${clean}`
  }

  return { resolveAssetUrl }
}
