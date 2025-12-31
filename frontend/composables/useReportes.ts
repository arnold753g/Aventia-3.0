import { useAuthStore } from '~/stores/auth'

type ReporteFormato = 'json' | 'pdf' | 'csv' | 'excel'

type ReporteParamsBase = {
  mes?: number
  anio?: number
  fechaInicio?: string
  fechaFin?: string
  formato?: ReporteFormato
}

type ReporteVentasParams = ReporteParamsBase & {
  paqueteId?: number
  tipoCompra?: 'compartido' | 'privado'
}

type ReportePaqueteParams = ReporteParamsBase & {
  paqueteId?: number
}

const buildQuery = (params: Record<string, string | number | undefined>) => {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      query.set(key, String(value))
    }
  })
  return query.toString()
}

const parseFilename = (contentDisposition?: string | null) => {
  if (!contentDisposition) return ''
  const match = contentDisposition.match(/filename\*?=(?:UTF-8''|")?([^\";]+)"?/i)
  if (!match) return ''
  return decodeURIComponent(match[1])
}

export const useReportes = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const authStore = useAuthStore()

  const authHeader = () => (authStore.token ? { Authorization: `Bearer ${authStore.token}` } : undefined)

  const requestReporte = async (agenciaId: number, endpoint: string, params: Record<string, string | number | undefined>) => {
    const qs = buildQuery(params)
    const url = qs ? `${apiBase}/agencias/${agenciaId}/reportes/${endpoint}?${qs}` : `${apiBase}/agencias/${agenciaId}/reportes/${endpoint}`
    const formato = String(params.formato || 'json')
    if (formato !== 'json') {
      let filename = ''
      const blob = await $fetch(url, {
        headers: authHeader(),
        responseType: 'blob',
        onResponse({ response }) {
          filename = parseFilename(response.headers.get('content-disposition'))
        }
      })
      return { blob, filename }
    }
    return $fetch(url, { headers: authHeader() })
  }

  const generarReporteVentas = async (agenciaId: number, params: ReporteVentasParams) => {
    return requestReporte(agenciaId, 'ventas', {
      mes: params.mes,
      anio: params.anio,
      fecha_inicio: params.fechaInicio,
      fecha_fin: params.fechaFin,
      formato: params.formato,
      paquete_id: params.paqueteId,
      tipo_compra: params.tipoCompra
    })
  }

  const generarReporteOcupacion = async (agenciaId: number, params: ReportePaqueteParams) => {
    return requestReporte(agenciaId, 'ocupacion', {
      mes: params.mes,
      anio: params.anio,
      fecha_inicio: params.fechaInicio,
      fecha_fin: params.fechaFin,
      formato: params.formato,
      paquete_id: params.paqueteId
    })
  }

  const generarReporteFinanciero = async (agenciaId: number, params: ReportePaqueteParams) => {
    return requestReporte(agenciaId, 'financiero', {
      mes: params.mes,
      anio: params.anio,
      fecha_inicio: params.fechaInicio,
      fecha_fin: params.fechaFin,
      formato: params.formato,
      paquete_id: params.paqueteId
    })
  }

  const generarReporteTuristas = async (agenciaId: number, params: ReporteParamsBase) => {
    return requestReporte(agenciaId, 'turistas', {
      mes: params.mes,
      anio: params.anio,
      fecha_inicio: params.fechaInicio,
      fecha_fin: params.fechaFin,
      formato: params.formato
    })
  }

  return {
    generarReporteVentas,
    generarReporteOcupacion,
    generarReporteFinanciero,
    generarReporteTuristas
  }
}
