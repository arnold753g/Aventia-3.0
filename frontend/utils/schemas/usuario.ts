import { z } from 'zod'
import { toFormValidator } from '@vee-validate/zod'

export const usuarioStep1Schema = z.object({
  nombre: z.string().min(1, 'El nombre es obligatorio'),
  apellido_paterno: z.string().min(1, 'El apellido paterno es obligatorio'),
  apellido_materno: z.string().optional(),
  telefono: z.string().min(8, 'El telefono debe tener al menos 8 digitos'),
  fecha_nacimiento: z.string().min(1, 'La fecha de nacimiento es obligatoria'),
  departamento: z.coerce.number({ required_error: 'Departamento requerido' }),
  provincia: z.coerce.number({ required_error: 'Provincia requerida' }),
  ci: z.string().min(5, 'El CI debe tener al menos 5 caracteres'),
  email: z.string().email('Email invalido')
})

export const usuarioStep1Resolver = () => toFormValidator(usuarioStep1Schema)
