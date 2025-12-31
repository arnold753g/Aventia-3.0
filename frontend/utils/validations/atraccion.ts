import { z } from 'zod'

// Validación de coordenadas GPS de Bolivia
const latitudBolivia = z.number()
  .min(-22.90, 'La latitud debe estar entre -22.90 y -9.67 (territorio boliviano)')
  .max(-9.67, 'La latitud debe estar entre -22.90 y -9.67 (territorio boliviano)')

const longitudBolivia = z.number()
  .min(-69.65, 'La longitud debe estar entre -69.65 y -57.45 (territorio boliviano)')
  .max(-57.45, 'La longitud debe estar entre -69.65 y -57.45 (territorio boliviano)')

// Validación de horarios
const horario = z.string()
  .regex(/^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/, 'Formato de hora inválido (HH:MM)')
  .optional()

// Esquema completo para crear atracción
export const createAtraccionSchema = z.object({
  nombre: z.string()
    .min(3, 'El nombre debe tener al menos 3 caracteres')
    .max(255, 'El nombre no puede exceder 255 caracteres')
    .transform(val => {
      return val
        .toLowerCase()
        .split(' ')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ')
    }),

  descripcion: z.string()
    .min(10, 'La descripción debe tener al menos 10 caracteres')
    .max(5000, 'La descripción no puede exceder 5000 caracteres'),

  provincia_id: z.number({
    required_error: 'Debe seleccionar una provincia'
  }).min(1, 'Debe seleccionar una provincia válida'),

  direccion: z.string().optional(),

  latitud: latitudBolivia.optional().nullable(),

  longitud: longitudBolivia.optional().nullable(),

  horario_apertura: horario,

  horario_cierre: horario,

  precio_entrada: z.number()
    .min(0, 'El precio no puede ser negativo')
    .max(10000, 'El precio parece demasiado alto')
    .optional()
    .default(0),

  nivel_dificultad: z.enum(['facil', 'medio', 'dificil', 'extremo'], {
    errorMap: () => ({ message: 'Seleccione un nivel de dificultad válido' })
  }).optional(),

  requiere_agencia: z.boolean().default(false),

  acceso_particular: z.boolean().default(true),

  mes_inicio_id: z.number()
    .min(1, 'Mes inválido')
    .max(12, 'Mes inválido')
    .optional()
    .nullable(),

  mes_fin_id: z.number()
    .min(1, 'Mes inválido')
    .max(12, 'Mes inválido')
    .optional()
    .nullable(),

  status: z.enum(['activa', 'inactiva', 'mantenimiento', 'fuera_temporada'])
    .default('activa'),

  visible_publico: z.boolean().default(true),

  telefono: z.string()
    .regex(/^[67]\d{7}$/, 'El teléfono debe tener 8 dígitos y comenzar con 6 o 7')
    .optional()
    .or(z.literal('')),

  email: z.string()
    .email('Ingrese un email válido')
    .optional()
    .or(z.literal('')),

  sitio_web: z.string()
    .url('Ingrese una URL válida')
    .optional()
    .or(z.literal('')),

  facebook: z.string().optional(),

  instagram: z.string().optional(),

  subcategorias_ids: z.array(z.number())
    .min(1, 'Debe seleccionar al menos una subcategoría')
    .max(4, 'Máximo 4 subcategorías permitidas'),

  dias_ids: z.array(z.number())
    .min(1, 'Debe seleccionar al menos un día')
    .max(7, 'No puede seleccionar más de 7 días'),

  fotos: z.array(z.string().url('URL de foto inválida'))
    .max(10, 'Máximo 10 fotos permitidas')
    .optional()
    .default([])
}).refine(data => {
  if (data.latitud && !data.longitud) return false
  if (data.longitud && !data.latitud) return false
  return true
}, {
  message: 'Debe proporcionar latitud y longitud juntas',
  path: ['latitud']
}).refine(data => {
  if (data.mes_inicio_id && !data.mes_fin_id) return false
  if (data.mes_fin_id && !data.mes_inicio_id) return false
  return true
}, {
  message: 'Debe proporcionar mes de inicio y fin juntos',
  path: ['mes_inicio_id']
})

// Esquema para actualizar atracción (todos los campos opcionales)
export const updateAtraccionSchema = z.object({
  nombre: z.string()
    .min(3, 'El nombre debe tener al menos 3 caracteres')
    .max(255, 'El nombre no puede exceder 255 caracteres')
    .optional(),

  descripcion: z.string()
    .min(10, 'La descripción debe tener al menos 10 caracteres')
    .max(5000, 'La descripción no puede exceder 5000 caracteres')
    .optional(),

  provincia_id: z.number().min(1, 'Seleccione una provincia válida').optional(),

  direccion: z.string().optional(),

  latitud: latitudBolivia.optional().nullable(),

  longitud: longitudBolivia.optional().nullable(),

  horario_apertura: horario,

  horario_cierre: horario,

  precio_entrada: z.number()
    .min(0, 'El precio no puede ser negativo')
    .max(10000, 'El precio parece demasiado alto')
    .optional(),

  nivel_dificultad: z.enum(['facil', 'medio', 'dificil', 'extremo']).optional(),

  requiere_agencia: z.boolean().optional(),

  acceso_particular: z.boolean().optional(),

  mes_inicio_id: z.number().min(1).max(12).optional().nullable(),

  mes_fin_id: z.number().min(1).max(12).optional().nullable(),

  status: z.enum(['activa', 'inactiva', 'mantenimiento', 'fuera_temporada']).optional(),

  visible_publico: z.boolean().optional(),

  telefono: z.string()
    .regex(/^[67]\d{7}$/, 'El teléfono debe tener 8 dígitos y comenzar con 6 o 7')
    .optional()
    .or(z.literal('')),

  email: z.string().email('Ingrese un email válido').optional().or(z.literal('')),

  sitio_web: z.string().url('Ingrese una URL válida').optional().or(z.literal('')),

  facebook: z.string().optional(),

  instagram: z.string().optional(),

  subcategorias_ids: z.array(z.number())
    .min(1, 'Debe seleccionar al menos una subcategoría')
    .max(4, 'Máximo 4 subcategorías permitidas')
    .optional(),

  dias_ids: z.array(z.number())
    .min(1, 'Debe seleccionar al menos un día')
    .max(7, 'No puede seleccionar más de 7 días')
    .optional(),

  fotos: z.array(z.string().url('URL de foto inválida'))
    .max(10, 'Máximo 10 fotos permitidas')
    .optional()
}).refine(data => {
  if (data.latitud && !data.longitud) return false
  if (data.longitud && !data.latitud) return false
  return true
}, {
  message: 'Debe proporcionar latitud y longitud juntas',
  path: ['latitud']
}).refine(data => {
  if (data.mes_inicio_id && !data.mes_fin_id) return false
  if (data.mes_fin_id && !data.mes_inicio_id) return false
  return true
}, {
  message: 'Debe proporcionar mes de inicio y fin juntos',
  path: ['mes_inicio_id']
})

// Tipos TypeScript
export type CreateAtraccionInput = z.infer<typeof createAtraccionSchema>
export type UpdateAtraccionInput = z.infer<typeof updateAtraccionSchema>
