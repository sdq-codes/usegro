export type InvoiceTemplate = 'classic' | 'modern' | 'minimal'

export const INVOICE_TEMPLATES = [
  {
    id: 'classic' as InvoiceTemplate,
    label: 'Classic',
    description: 'Warm orange accents with a structured card layout',
  },
  {
    id: 'modern' as InvoiceTemplate,
    label: 'Modern',
    description: 'Bold dark header with a clean, contemporary design',
  },
  {
    id: 'minimal' as InvoiceTemplate,
    label: 'Minimal',
    description: 'Typography-first design with subtle blue accents',
  },
]
