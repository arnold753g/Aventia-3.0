import { z } from 'zod'
import { DEFAULT_PHONE_PREFIX } from '~/utils/phone'

// Validador de telefono boliviano
const phonePrefixBolivia = z.string()
  .refine((value) => value === DEFAULT_PHONE_PREFIX, 'El prefijo debe ser +591')

const phoneNumberBolivia = z.string()
  .length(8, 'El telefono debe tener 8 digitos')
  .regex(/^[2-7]\d{7}$/, 'El telefono debe comenzar con 2-7')

// Validador de email
const email = z.string()
  .email('Email invalido')
  .toLowerCase()

// Validador de URL opcional
const urlOpcional = z.string()
  .url('URL invalida')
  .or(z.literal(''))
  .optional()

// Validador de coordenadas GPS (Bolivia)
const latitud = z.number()
  .min(-22.9, 'Latitud fuera del rango de Bolivia')
  .max(-9.67, 'Latitud fuera del rango de Bolivia')
  .optional()
  .nullable()

const longitud = z.number()
  .min(-69.65, 'Longitud fuera del rango de Bolivia')
  .max(-57.45, 'Longitud fuera del rango de Bolivia')
  .optional()
  .nullable()

// Validador de horario (HH:MM)
const horario = z.string()
  .regex(/^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/, 'Formato de hora invalido (HH:MM)')
  .optional()
  .or(z.literal(''))

// Schema para creacion RAPIDA (solo campos esenciales)
export const createAgenciaRapidaSchema = z.object({
  nombre_comercial: z.string()
    .min(2, 'El nombre debe tener al menos 2 caracteres')
    .max(255, 'El nombre no puede exceder 255 caracteres'),

  departamento_id: z.number({
    required_error: 'Seleccione un departamento'
  }).min(1, 'Seleccione un departamento'),

  phone_prefix: phonePrefixBolivia,

  phone_number: phoneNumberBolivia,

  encargado_principal_id: z.number({
    required_error: 'Seleccione un encargado'
  }).min(1, 'Seleccione un encargado')
})

// Schema para creacion COMPLETA
export const createAgenciaCompletaSchema = z.object({
  nombre_comercial: z.string()
    .min(2, 'El nombre debe tener al menos 2 caracteres')
    .max(255, 'El nombre no puede exceder 255 caracteres'),

  descripcion: z.string()
    .optional()
    .or(z.literal('')),

  direccion: z.string()
    .min(5, 'La direccion debe tener al menos 5 caracteres'),

  departamento_id: z.number({
    required_error: 'Seleccione un departamento'
  }).min(1, 'Seleccione un departamento'),

  latitud,
  longitud,

  phone_prefix: phonePrefixBolivia,

  phone_number: phoneNumberBolivia,

  email: email,

  sitio_web: urlOpcional,
  facebook: urlOpcional,
  instagram: z.string().optional().or(z.literal('')),

  licencia_turistica: z.boolean().default(false),

  horario_apertura: horario,
  horario_cierre: horario,

  acepta_qr: z.boolean().default(true),
  acepta_transferencia: z.boolean().default(true),
  acepta_efectivo: z.boolean().default(true),

  encargado_principal_id: z.number({
    required_error: 'Seleccione un encargado'
  }).min(1, 'Seleccione un encargado'),

  status: z.enum(['activa', 'inactiva', 'suspendida', 'en_revision'])
    .default('activa'),

  visible_publico: z.boolean().default(true),

  dias_ids: z.array(z.number()).optional()
})

// Schema para actualizacion (todos los campos opcionales)
export const updateAgenciaSchema = createAgenciaCompletaSchema.partial()
