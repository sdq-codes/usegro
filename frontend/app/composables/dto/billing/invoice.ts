export interface InvoiceLineItem {
  id: string
  name: string
  type: 'product' | 'service'
  catalog_id: string
  qty: number
  rate: number
  amount: number
  billing_type: 'one-time' | 'recurring'
  billing_cycle: string
  image_url: string
}

export type InvoiceStatus = 'draft' | 'sent' | 'paid' | 'overdue' | 'cancelled'

export interface Invoice {
  id: string
  crm_id: string
  invoice_number: string
  status: InvoiceStatus
  customer_ids: string[]
  customer_names: string[]
  customer_emails: string[]
  line_items: InvoiceLineItem[]
  tax_rate: number
  subtotal: number
  tax_amount: number
  total: number
  subject: string
  terms_and_conditions: string
  reference_number: string
  memo: string
  due_date: string | null
  recurring_start_date: string | null
  sent_at: string | null
  created_at: string
  updated_at: string
}

export interface CreateInvoicePayload {
  customer_ids: string[]
  customer_names: string[]
  customer_emails: string[]
  line_items: CreateLineItemPayload[]
  tax_rate: number
  subject?: string
  terms_and_conditions?: string
  reference_number?: string
  memo?: string
  due_date?: string | null
  recurring_start_date?: string | null
}

export interface CreateLineItemPayload {
  name: string
  type: 'product' | 'service'
  catalog_id?: string
  qty: number
  rate: number
  billing_type: 'one-time' | 'recurring'
  billing_cycle?: string
  image_url?: string
}

export interface UpdateInvoicePayload {
  customer_ids?: string[]
  customer_names?: string[]
  customer_emails?: string[]
  line_items?: CreateLineItemPayload[]
  tax_rate?: number
  subject?: string
  terms_and_conditions?: string
  reference_number?: string
  memo?: string
  due_date?: string | null
  recurring_start_date?: string | null
}

export interface ListInvoicesResponse {
  data: Invoice[]
  total: number
  page: number
  limit: number
}
