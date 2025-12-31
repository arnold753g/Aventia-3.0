import { z } from 'zod'
import { PHONE_PREFIX_VALUES } from '~/utils/phone'

const ciBoliviano = z.string()
  .min(5, 'El CI debe tener al menos 5 caracteres')
  .max(15, 'El CI no puede tener más de 15 caracteres')
  .regex(/^[0-9A-Za-z]+$/, 'El CI solo puede contener letras y números')
  .transform(val => val.toUpperCase())

const phonePrefixSchema = z.string()
  .min(1, 'Seleccione un prefijo valido')
  .refine((value) => PHONE_PREFIX_VALUES.includes(value), 'Seleccione un prefijo valido')

const phoneNumberSchema = z.string()
  .min(1, 'Ingrese el numero de telefono')
  .regex(/^\d+$/, 'El telefono solo debe contener numeros')
  .min(4, 'El telefono debe tener al menos 4 digitos')
  .max(14, 'El telefono no puede tener mas de 14 digitos')

const validatePhoneByPrefix = (
  data: { phone_prefix?: string; phone_number?: string },
  ctx: z.RefinementCtx
) => {
  const prefix = data.phone_prefix
  const number = data.phone_number

  if (!prefix && !number) return

  if (!prefix) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ['phone_prefix'],
      message: 'Seleccione un prefijo valido'
    })
  }

  if (!number) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ['phone_number'],
      message: 'Ingrese el numero de telefono'
    })
  }

  if (!prefix || !number) return

}

const email = z.string()
  .email('Ingrese un email válido')
  .toLowerCase()

// Validacion de fecha de nacimiento (mayor de 18 anos)
const fechaNacimiento = z.union([z.string(), z.date()])
  .transform((value) => {
    if (value instanceof Date) return value.toISOString().slice(0, 10)
    return value
  })
  .refine((date) => {
    if (!date) return false
    const birthDate = new Date(`${date}T00:00:00`)
    if (Number.isNaN(birthDate.getTime())) return false

    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const cutoff = new Date(today)
    cutoff.setFullYear(cutoff.getFullYear() - 18)

    return birthDate <= cutoff
  }, 'Debe ser mayor de 18 anos')

const password = z.string()
  .min(8, 'La contraseña debe tener al menos 8 caracteres')
  .regex(/[A-Z]/, 'Debe contener al menos una letra mayúscula')
  .regex(/[a-z]/, 'Debe contener al menos una letra minúscula')
  .regex(/[0-9]/, 'Debe contener al menos un número')
  .regex(/[^A-Za-z0-9]/, 'Debe contener al menos un carácter especial')

// Esquema completo para crear usuario
export const createUsuarioSchema = z.object({
  nombre: z.string()
    .min(2, 'El nombre debe tener al menos 2 caracteres')
    .max(100, 'El nombre no puede exceder 100 caracteres'),

  apellido_paterno: z.string()
    .min(2, 'El apellido paterno debe tener al menos 2 caracteres')
    .max(100, 'El apellido paterno no puede exceder 100 caracteres'),

  apellido_materno: z.string()
    .min(2, 'El apellido materno debe tener al menos 2 caracteres')
    .max(100, 'El apellido materno no puede exceder 100 caracteres'),

  ci: ciBoliviano,

  expedido: z.enum(['LP', 'CB', 'SC', 'PT', 'OR', 'TJ', 'CH', 'BE', 'BN', 'PD'], {
    errorMap: () => ({ message: 'Seleccione un departamento válido' })
  }),

  phone_prefix: phonePrefixSchema,
  phone_number: phoneNumberSchema,

  email: email,

  fecha_nacimiento: fechaNacimiento,

  ciudad: z.string().optional(),

  nationality: z.string().default('Bolivia'),

  password: password,

  rol: z.enum(['admin', 'turista', 'encargado_agencia'], {
    errorMap: () => ({ message: 'Seleccione un rol válido' })
  })
}).superRefine(validatePhoneByPrefix)

// Esquema para actualizar usuario (sin password y campos no editables)
export const updateUsuarioSchema = z.object({
  nombre: z.string()
    .min(2, 'El nombre debe tener al menos 2 caracteres')
    .max(100, 'El nombre no puede exceder 100 caracteres')
    .optional(),

  apellido_paterno: z.string()
    .min(2, 'El apellido paterno debe tener al menos 2 caracteres')
    .max(100, 'El apellido paterno no puede exceder 100 caracteres')
    .optional(),

  apellido_materno: z.string()
    .min(2, 'El apellido materno debe tener al menos 2 caracteres')
    .max(100, 'El apellido materno no puede exceder 100 caracteres')
    .optional(),

  phone_prefix: phonePrefixSchema.optional(),
  phone_number: phoneNumberSchema.optional(),

  ciudad: z.string().optional(),

  profile_photo: z.string().url('Ingrese una URL válida').optional()
}).superRefine(validatePhoneByPrefix)

// Esquema para registro público (con confirmación de contraseña, sin rol)
export const registroPublicoSchema = z.object({
  nombre: z.string()
    .min(2, 'El nombre debe tener al menos 2 caracteres')
    .max(100, 'El nombre no puede exceder 100 caracteres'),

  apellido_paterno: z.string()
    .min(2, 'El apellido paterno debe tener al menos 2 caracteres')
    .max(100, 'El apellido paterno no puede exceder 100 caracteres'),

  apellido_materno: z.string()
    .min(2, 'El apellido materno debe tener al menos 2 caracteres')
    .max(100, 'El apellido materno no puede exceder 100 caracteres')
    .optional()
    .or(z.literal('')),

  ci: ciBoliviano,

  expedido: z.enum(['LP', 'CB', 'SC', 'PT', 'OR', 'TJ', 'CH', 'BE', 'BN', 'PD'], {
    errorMap: () => ({ message: 'Seleccione un departamento válido' })
  }),

  phone_prefix: phonePrefixSchema,
  phone_number: phoneNumberSchema,

  email: email,

  fecha_nacimiento: fechaNacimiento,

  ciudad: z.string().optional().or(z.literal('')),

  password: password,

  confirmPassword: z.string().min(1, 'Confirme su contraseña')
}).superRefine(validatePhoneByPrefix).refine((data) => data.password === data.confirmPassword, {
  message: 'Las contraseñas no coinciden',
  path: ['confirmPassword']
})

// Tipos TypeScript generados desde los esquemas
export type CreateUsuarioInput = z.infer<typeof createUsuarioSchema>
export type UpdateUsuarioInput = z.infer<typeof updateUsuarioSchema>
export type RegistroPublicoInput = z.infer<typeof registroPublicoSchema>

