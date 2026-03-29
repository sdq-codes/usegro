import { useApi } from "../useApi"
import type {ApiResult} from "@/composables/helpers/types/api";
import type { RawAxiosResponseHeaders, AxiosHeaders } from "axios";
import type {FormSubmission} from "@/composables/dto/customer/customer";

export interface PaginatedCustomers {
  data: FormSubmission[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export const useCustomerAPI = () => {
  const api = useApi()

  const FetchCustomers = async (page = 1, limit = 20): Promise<ApiResult<PaginatedCustomers, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.get(`/crm/customers`, { params: { page, limit } })
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Form not found"  }
    }
  }

  const DeleteCustomer = async (id:string, formId: string): Promise<ApiResult<null, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.delete<FormSubmission>(`/crm/customers/${formId}/${id}`)

      return { success: true, data: null, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Form not found"  }
    }
  }

  const FetchCustomer = async (formId: string, id:string): Promise<ApiResult<FormSubmission, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.get<FormSubmission>(`/crm/customers/${formId}/${id}`)

      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Form not found"  }
    }
  }

  const FetchCustomerActivity = async (submissionId: string): Promise<ApiResult<null, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.get(`/crm/customers/${submissionId}/activity`)
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Activity not found" }
    }
  }

  const PostComment = async (submissionId: string, comment: string): Promise<ApiResult<null, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post(`/crm/customers/${submissionId}/activity`, { comment })
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Failed to post comment" }
    }
  }

  return {
    FetchCustomers,
    FetchCustomer,
    DeleteCustomer,
    FetchCustomerActivity,
    PostComment,
  }
}
