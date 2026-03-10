import { useApi } from "../useApi"
import type { CreateCRMPayload, UpdateCRMPayload, CRMApiResponse } from "@/composables/dto/crm"
import type { ApiResult } from "@/composables/helpers/types/api"
import type { RawAxiosResponseHeaders, AxiosHeaders } from "axios"

type Result<T = CRMApiResponse> = Promise<ApiResult<T, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>>

export const useCRMAPI = () => {
  const api = useApi()

  const CheckBusinessNameExists = async (business_name: string): Promise<{ exists: boolean; error?: string }> => {
    try {
      const response = await api.post("/crm/business-name/exist", { business_name })
      return { exists: response.data?.data?.exists ?? false }
    } catch (error: unknown) {
      return { exists: false, error: error.response?.data?.response_message || "Failed to check business name" }
    }
  }

  const ListCRMs = async (): Result<CRMApiResponse[]> => {
    try {
      const response = await api.get("/crm")
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to fetch CRMs" }
    }
  }

  const GetCRM = async (id: string): Result => {
    try {
      const response = await api.get(`/crm/${id}`)
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to fetch CRM" }
    }
  }

  const CreateCRM = async (data: CreateCRMPayload): Result => {
    try {
      const response = await api.post("/crm", data)
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to create CRM" }
    }
  }

  const UpdateCRM = async (id: string, data: UpdateCRMPayload): Result => {
    try {
      const response = await api.patch(`/crm/${id}`, data)
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to update CRM" }
    }
  }

  const ToggleCRMStatus = async (): Result => {
    try {
      const response = await api.patch("/crm/status")
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to toggle CRM status" }
    }
  }

  const CreateSalesChannels = async (id: string, sales_channel_type: string[]): Result => {
    try {
      const response = await api.post(`/crm/${id}/sales-channels`, { sales_channel_type })
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to save sales channels" }
    }
  }

  const UpdateSalesChannels = async (id: string, sales_channel_type: string[]): Result => {
    try {
      const response = await api.patch(`/crm/${id}/sales-channels`, { sales_channel_type })
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to update sales channels" }
    }
  }

  const CreateStockProducts = async (id: string, product_type: string[]): Result => {
    try {
      const response = await api.post(`/crm/${id}/stock-products`, { product_type })
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to save product types" }
    }
  }

  const UpdateStockProducts = async (id: string, product_type: string[]): Result => {
    try {
      const response = await api.patch(`/crm/${id}/stock-products`, { product_type })
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to update product types" }
    }
  }

  return {
    CheckBusinessNameExists,
    ListCRMs,
    GetCRM,
    CreateCRM,
    UpdateCRM,
    ToggleCRMStatus,
    CreateSalesChannels,
    UpdateSalesChannels,
    CreateStockProducts,
    UpdateStockProducts,
  }
}
