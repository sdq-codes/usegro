import { useApi } from "../useApi"
import type {ApiResult} from "@/composables/helpers/types/api";
import type { RawAxiosResponseHeaders, AxiosHeaders } from "axios";
import type {FormSubmission} from "@/composables/dto/customer/customer";

export const useCustomerAPI = () => {
  const api = useApi()

  const FetchCustomers = async (): Promise<ApiResult<FormSubmission, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.get<FormSubmission>(`/crm/customers`)

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

  return {
    FetchCustomers,
    FetchCustomer,
    DeleteCustomer,
  }
}
