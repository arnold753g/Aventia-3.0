export const PHONE_PREFIXES = [
  { label: 'Estados Unidos / Canada (+1)', value: '+1' },
  { label: 'Mexico (+52)', value: '+52' },
  { label: 'Espana (+34)', value: '+34' },
  { label: 'Argentina (+54)', value: '+54' },
  { label: 'Bolivia (+591)', value: '+591' },
  { label: 'Brasil (+55)', value: '+55' },
  { label: 'Chile (+56)', value: '+56' },
  { label: 'Colombia (+57)', value: '+57' },
  { label: 'Costa Rica (+506)', value: '+506' },
  { label: 'Cuba (+53)', value: '+53' },
  { label: 'Ecuador (+593)', value: '+593' },
  { label: 'El Salvador (+503)', value: '+503' },
  { label: 'Guatemala (+502)', value: '+502' },
  { label: 'Honduras (+504)', value: '+504' },
  { label: 'Nicaragua (+505)', value: '+505' },
  { label: 'Panama (+507)', value: '+507' },
  { label: 'Paraguay (+595)', value: '+595' },
  { label: 'Peru (+51)', value: '+51' },
  { label: 'Uruguay (+598)', value: '+598' },
  { label: 'Venezuela (+58)', value: '+58' },
  { label: 'Portugal (+351)', value: '+351' },
  { label: 'Francia (+33)', value: '+33' },
  { label: 'Reino Unido (+44)', value: '+44' },
  { label: 'Alemania (+49)', value: '+49' },
  { label: 'Italia (+39)', value: '+39' },
  { label: 'China (+86)', value: '+86' },
  { label: 'Japon (+81)', value: '+81' },
  { label: 'India (+91)', value: '+91' }
] as const

export const DEFAULT_PHONE_PREFIX = '+591'

export const PHONE_PREFIX_VALUES = PHONE_PREFIXES.map((prefix) => prefix.value)

const PHONE_PREFIX_VALUES_BY_LENGTH = [...PHONE_PREFIX_VALUES].sort((a, b) => b.length - a.length)

const digitsOnly = (value: string) => value.replace(/\D/g, '')

export const buildPhone = (prefix: string, number: string) => {
  const cleanPrefix = PHONE_PREFIX_VALUES.includes(prefix) ? prefix : DEFAULT_PHONE_PREFIX
  const cleanNumber = digitsOnly(number || '')
  return `${cleanPrefix}${cleanNumber}`
}

export const splitPhone = (phone?: string | null) => {
  if (!phone) {
    return { prefix: DEFAULT_PHONE_PREFIX, number: '' }
  }

  const trimmed = phone.trim()
  if (!trimmed) {
    return { prefix: DEFAULT_PHONE_PREFIX, number: '' }
  }

  const compact = trimmed.replace(/\s+/g, '').replace(/-/g, '')
  if (compact.startsWith('+')) {
    const match = PHONE_PREFIX_VALUES_BY_LENGTH.find((prefix) => compact.startsWith(prefix))
    if (match) {
      return { prefix: match, number: digitsOnly(compact.slice(match.length)) }
    }
  }

  return { prefix: DEFAULT_PHONE_PREFIX, number: digitsOnly(compact) }
}
