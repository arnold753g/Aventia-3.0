import { defineStore } from 'pinia'
import type {
  DashboardModule,
  HeroSlide,
  PopularAttraction,
  PopularPackage,
  RecentItem
} from '~/types/dashboard'

const getPrimaryFoto = (atraccion: any) => {
  const principal = atraccion?.fotos?.find((foto: any) => foto?.es_principal)
  return principal?.foto || atraccion?.fotos?.[0]?.foto || ''
}

const formatDuration = (paquete: any) => {
  const dias = paquete?.duracion_dias ?? paquete?.duracionDias
  const noches = paquete?.duracion_noches ?? paquete?.duracionNoches
  if (!dias && !noches) return 'Duracion flexible'
  if (dias && noches) return `${dias} dias / ${noches} noches`
  if (dias) return `${dias} dias`
  return `${noches} noches`
}

const formatHorario = (atraccion: any) => {
  const apertura = atraccion?.horario_apertura ?? atraccion?.horarioApertura
  const cierre = atraccion?.horario_cierre ?? atraccion?.horarioCierre
  const aperturaShort = apertura ? String(apertura).slice(0, 5) : ''
  const cierreShort = cierre ? String(cierre).slice(0, 5) : ''
  if (aperturaShort && cierreShort) return `${aperturaShort} - ${cierreShort}`
  return aperturaShort || cierreShort || ''
}

const formatMoney = (value: number) => {
  return value.toLocaleString('es-BO', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

const formatEntryPrice = (atraccion: any) => {
  const precio = Number(atraccion?.precio_entrada ?? 0)
  if (!Number.isFinite(precio) || precio <= 0) return 'Libre'
  return `Bs. ${formatMoney(precio)}`
}

const formatDescription = (value?: string) => {
  const clean = String(value || '').replace(/\s+/g, ' ').trim()
  if (!clean) return 'Explora una atraccion turistica destacada.'
  if (clean.length <= 140) return clean
  return `${clean.slice(0, 137)}...`
}

const formatLocation = (atraccion: any) => {
  const provincia = atraccion?.provincia?.nombre || ''
  const departamento = atraccion?.provincia?.departamento?.nombre || ''
  if (provincia && departamento) return `${provincia}, ${departamento}`
  return provincia || departamento || 'Bolivia'
}

const formatShortDate = (value?: string) => {
  if (!value) return 'Fecha reciente'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Fecha reciente'
  const months = ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic']
  const day = String(date.getDate()).padStart(2, '0')
  const month = months[date.getMonth()] || ''
  const year = date.getFullYear()
  return `${day} ${month} ${year}`
}

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    heroSlides: [] as HeroSlide[],
    modules: [
      {
        id: 1,
        title: 'Atracciones turisticas',
        description: 'Explora destinos publicos y su informacion.',
        image: 'uploads/fotografias/atracciones/atraccion_1766433805791432200.jpg',
        icon: 'pi pi-map-marker',
        route: '/atracciones'
      },
      {
        id: 2,
        title: 'Paquetes turisticos',
        description: 'Colecciones exclusivas para tu viaje.',
        image: 'uploads/fotografias/paquetes/paquete_1766409161229754000.jpg',
        icon: 'pi pi-briefcase',
        route: '/paquetes',
        requiresLogin: true
      },
      {
        id: 3,
        title: 'Salidas habilitadas',
        description: 'Fechas confirmadas y cupos disponibles.',
        image: 'uploads/fotografias/paquetes/paquete_1766069742054348400.jpg',
        icon: 'pi pi-calendar',
        route: '/salidas',
        requiresLogin: true
      }
    ] as DashboardModule[],
    popularAttractions: [] as PopularAttraction[],
    popularPackages: [] as PopularPackage[],
    recentItems: [] as RecentItem[],
    loadingAttractions: false,
    attractionsLoaded: false,
    attractionsError: '',
    loadingPackages: false,
    packagesLoaded: false,
    packagesError: ''
  }),

  actions: {
    async loadAttractions() {
      if (this.loadingAttractions || this.attractionsLoaded) return
      this.loadingAttractions = true
      this.attractionsError = ''
      try {
        const { getAtracciones } = useAtracciones()
        const response: any = await getAtracciones({
          page: 1,
          limit: 12,
          visible_publico: 'true',
          status: 'activa',
          sort_by: 'created_at',
          sort_order: 'desc'
        })

        if (!response?.success) {
          this.attractionsError = response?.error?.message || 'No se pudieron cargar las atracciones'
          return
        }

        const atracciones = response.data?.atracciones || []
        this.heroSlides = atracciones.slice(0, 4).map((atraccion: any) => {
          const location = formatLocation(atraccion)
          const horario = formatHorario(atraccion)
          return {
            id: atraccion.id,
            title: atraccion.nombre || 'Atraccion turistica',
            location,
            description: formatDescription(atraccion.descripcion),
            image: getPrimaryFoto(atraccion),
            route: `/atracciones/${atraccion.id}`,
            requiresLogin: false,
            meta: [
              ...(location ? [{ icon: 'pi pi-map-marker', text: location }] : []),
              ...(horario ? [{ icon: 'pi pi-clock', text: `Horario ${horario}` }] : [])
            ],
            ctaLabel: 'Ver atraccion',
            priceLabel: 'Entrada',
            priceValue: formatEntryPrice(atraccion)
          }
        })
        this.popularAttractions = atracciones.slice(0, 6).map((atraccion: any) => ({
          id: atraccion.id,
          name: atraccion.nombre || 'Atraccion',
          place: formatLocation(atraccion),
          image: getPrimaryFoto(atraccion),
          route: `/atracciones/${atraccion.id}`
        }))

        this.recentItems = atracciones.slice(0, 6).map((atraccion: any) => ({
          id: atraccion.id,
          type: 'atraccion',
          date: formatShortDate(atraccion.created_at),
          title: atraccion.nombre || 'Atraccion turistica',
          image: getPrimaryFoto(atraccion),
          route: `/atracciones/${atraccion.id}`
        }))

        this.attractionsLoaded = true
      } catch (error: any) {
        this.attractionsError = error?.data?.error?.message || error?.message || 'No se pudieron cargar las atracciones'
      } finally {
        this.loadingAttractions = false
      }
    },
    async loadPackages() {
      if (this.loadingPackages || this.packagesLoaded) return
      this.loadingPackages = true
      this.packagesError = ''
      try {
        const { getPaquetes } = usePaquetesTuristicos()
        const response: any = await getPaquetes({
          page: 1,
          limit: 8,
          sort_by: 'created_at',
          sort_order: 'desc'
        })

        if (!response?.success) {
          this.packagesError = response?.error?.message || 'No se pudieron cargar los paquetes'
          return
        }

        const paquetes = response.data?.paquetes || []
        this.popularPackages = paquetes.slice(0, 6).map((paquete: any) => ({
          id: paquete.id,
          name: paquete.nombre || 'Paquete turistico',
          duration: formatDuration(paquete),
          priceFrom: String(paquete.precio_base_nacionales ?? paquete.precioBaseNacionales ?? '0'),
          image: getPrimaryFoto(paquete),
          route: '/paquetes',
          requiresLogin: true
        }))

        this.packagesLoaded = true
      } catch (error: any) {
        this.packagesError = error?.data?.error?.message || error?.message || 'No se pudieron cargar los paquetes'
      } finally {
        this.loadingPackages = false
      }
    }
  }
})
