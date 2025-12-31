export type DashboardItemType = 'atraccion' | 'paquete' | 'salida'

export interface HeroSlide {
  id: number
  title: string
  location?: string
  departureDate?: string
  priceFrom?: string
  seats?: number
  description: string
  image: string
  route: string
  requiresLogin?: boolean
  meta?: { icon: string; text: string }[]
  ctaLabel?: string
  priceLabel?: string
  priceValue?: string
}

export interface DashboardModule {
  id: number
  title: string
  description: string
  image: string
  icon: string
  route: string
  requiresLogin?: boolean
}

export interface PopularAttraction {
  id: number
  name: string
  place: string
  image: string
  route: string
}

export interface PopularPackage {
  id: number
  name: string
  duration: string
  priceFrom: string
  image: string
  route: string
  requiresLogin?: boolean
}

export interface RecentItem {
  id: number
  type: DashboardItemType
  date: string
  title: string
  image: string
  route: string
  requiresLogin?: boolean
}
