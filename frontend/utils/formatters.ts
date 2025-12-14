// Formatear fecha
export const formatDate = (date: string | Date) => {
  return new Date(date).toLocaleDateString('es-BO', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// Formatear fecha corta
export const formatDateShort = (date: string | Date) => {
  return new Date(date).toLocaleDateString('es-BO', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

// Formatear fecha y hora
export const formatDateTime = (date: string | Date) => {
  return new Date(date).toLocaleString('es-BO', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Obtener iniciales de nombre completo
export const getInitials = (nombre: string, apellido: string) => {
  return `${nombre.charAt(0)}${apellido.charAt(0)}`.toUpperCase()
}

// Formatear nombre completo
export const getFullName = (nombre: string, apellidoPaterno: string, apellidoMaterno?: string) => {
  return apellidoMaterno
    ? `${nombre} ${apellidoPaterno} ${apellidoMaterno}`
    : `${nombre} ${apellidoPaterno}`
}

// Obtener color por rol
export const getRolColor = (rol: string) => {
  const colors: Record<string, string> = {
    admin: 'purple',
    turista: 'blue',
    encargado_agencia: 'green'
  }
  return colors[rol] || 'gray'
}

// Obtener color por status
export const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: 'green',
    inactive: 'orange',
    suspended: 'red'
  }
  return colors[status] || 'gray'
}

// Obtener label de rol
export const getRolLabel = (rol: string) => {
  const labels: Record<string, string> = {
    admin: 'Administrador',
    turista: 'Turista',
    encargado_agencia: 'Encargado de Agencia'
  }
  return labels[rol] || rol
}

// Obtener label de status
export const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    active: 'Activo',
    inactive: 'Inactivo',
    suspended: 'Suspendido'
  }
  return labels[status] || status
}

// Formatear teléfono boliviano
export const formatPhone = (phone: string) => {
  if (!phone || phone.length !== 8) return phone
  return `${phone.substring(0, 4)}-${phone.substring(4)}`
}

// Obtener icono por rol
export const getRolIcon = (rol: string) => {
  const icons: Record<string, string> = {
    admin: 'pi-shield',
    turista: 'pi-user',
    encargado_agencia: 'pi-building'
  }
  return icons[rol] || 'pi-user'
}
