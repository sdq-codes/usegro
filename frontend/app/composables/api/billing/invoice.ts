import { useApi } from '../useApi'
import type {
  Invoice,
  CreateInvoicePayload,
  UpdateInvoicePayload,
  ListInvoicesResponse,
} from '@/composables/dto/billing/invoice'

interface ApiResult<T> {
  success: boolean
  data?: T
  error?: string
}

export const useInvoiceAPI = () => {
  const api = useApi()

  const CreateInvoice = async (body: CreateInvoicePayload): Promise<ApiResult<{ data: Invoice }>> => {
    try {
      const response = await api.post('/billing/invoices', body)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to create invoice' }
    }
  }

  const ListInvoices = async (params?: {
    page?: number
    limit?: number
    status?: string
    customer_name?: string
    invoice_number?: string
    due_date_from?: string
    due_date_to?: string
    created_from?: string
    created_to?: string
    amount_min?: number
    amount_max?: number
    billing_type?: string
  }): Promise<ApiResult<{ data: ListInvoicesResponse }>> => {
    try {
      // Strip empty string / undefined values so they don't pollute the query string
      const cleaned = params
        ? Object.fromEntries(Object.entries(params).filter(([, v]) => v !== '' && v !== undefined && v !== null))
        : {}
      const response = await api.get('/billing/invoices', { params: cleaned })
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch invoices' }
    }
  }

  const GetInvoice = async (id: string): Promise<ApiResult<{ data: Invoice }>> => {
    try {
      const response = await api.get(`/billing/invoices/${id}`)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Invoice not found' }
    }
  }

  const UpdateInvoice = async (id: string, body: UpdateInvoicePayload): Promise<ApiResult<{ data: Invoice }>> => {
    try {
      const response = await api.patch(`/billing/invoices/${id}`, body)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to update invoice' }
    }
  }

  const DeleteInvoice = async (id: string): Promise<ApiResult<null>> => {
    try {
      await api.delete(`/billing/invoices/${id}`)
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to delete invoice' }
    }
  }

  const SendInvoice = async (id: string): Promise<ApiResult<{ data: Invoice }>> => {
    try {
      const response = await api.post(`/billing/invoices/${id}/send`)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to send invoice' }
    }
  }

  return { CreateInvoice, ListInvoices, GetInvoice, UpdateInvoice, DeleteInvoice, SendInvoice }
}
