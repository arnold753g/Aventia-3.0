// Obtener label de nivel de dificultad
export const getNivelDificultadLabel = (nivel: string) => {
  const labels: Record<string, string> = {
    facil: 'Fácil',
    medio: 'Medio',
    dificil: 'Difícil',
    extremo: 'Extremo'
  }
  return labels[nivel] || nivel
}

// Obtener color por nivel de dificultad
export const getNivelDificultadColor = (nivel: string) => {
  const colors: Record<string, string> = {
    facil: 'green',
    medio: 'blue',
    dificil: 'orange',
    extremo: 'red'
  }
  return colors[nivel] || 'gray'
}

// Obtener icono por nivel de dificultad
export const getNivelDificultadIcon = (nivel: string) => {
  const icons: Record<string, string> = {
    facil: 'pi-check-circle',
    medio: 'pi-info-circle',
    dificil: 'pi-exclamation-triangle',
    extremo: 'pi-times-circle'
  }
  return icons[nivel] || 'pi-circle'
}

// Obtener label de status
export const getStatusAtraccionLabel = (status: string) => {
  const labels: Record<string, string> = {
    activa: 'Activa',
    inactiva: 'Inactiva',
    mantenimiento: 'En Mantenimiento',
    fuera_temporada: 'Fuera de Temporada'
  }
  return labels[status] || status
}

// Obtener color por status
export const getStatusAtraccionColor = (status: string) => {
  const colors: Record<string, string> = {
    activa: 'green',
    inactiva: 'gray',
    mantenimiento: 'orange',
    fuera_temporada: 'blue'
  }
  return colors[status] || 'gray'
}

// Formatear precio boliviano
export const formatPrecioBoliviano = (precio: number) => {
  if (precio === 0) return 'Gratis'
  return `Bs. ${precio.toFixed(2)}`
}

export const extractHora = (value?: string) => {
  if (!value) return ''
  const match = value.match(/(\d{1,2}):(\d{2})/)
  if (!match) return value
  const hh = match[1].padStart(2, '0')
  const mm = match[2]
  return `${hh}:${mm}`
}

// Formatear horario
export const formatHorario = (apertura?: string, cierre?: string) => {
  const aperturaHora = extractHora(apertura)
  const cierreHora = extractHora(cierre)
  if (!aperturaHora || !cierreHora) return 'Horario no especificado'
  return `${aperturaHora} - ${cierreHora}`
}

// Formatear mejor época
export const formatMejorEpoca = (mesInicio?: any, mesFin?: any) => {
  if (!mesInicio || !mesFin) return 'Todo el año'
  return `${mesInicio.nombre} - ${mesFin.nombre}`
}

// Obtener iniciales para avatar de categoría
export const getCategoriaInitials = (nombre: string) => {
  const words = nombre.split(' ')
  if (words.length >= 2) {
    return `${words[0].charAt(0)}${words[1].charAt(0)}`.toUpperCase()
  }
  return nombre.substring(0, 2).toUpperCase()
}
