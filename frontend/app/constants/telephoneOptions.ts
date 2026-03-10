interface CountryOption {
  value: string
  label: string
  flag: string
}

export const countryTelephoneOptions: CountryOption[] = [
  { value: 'us', label: '+1', flag: '🇺🇸' },
  { value: 'ng', label: '+234', flag: '🇳🇬' },
  { value: 'gb', label: '+44', flag: '🇬🇧' }
]
