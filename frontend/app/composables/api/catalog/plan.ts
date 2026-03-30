import { useApi } from '../useApi'
import type { CreatePlanDTO, Plan } from '@/composables/dto/catalog/service'

interface ApiResult<T> {
  success: boolean
  data?: T
  error?: string
}

export const usePlanAPI = () => {
  const api = useApi()

  const CreatePlan = async (serviceId: string, body: CreatePlanDTO): Promise<ApiResult<{ data: Plan }>> => {
    try {
      const response = await api.post(`/catalog/services/${serviceId}/plans`, body)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to create plan' }
    }
  }

  const ListPlans = async (serviceId: string): Promise<ApiResult<{ data: Plan[] }>> => {
    try {
      const response = await api.get(`/catalog/services/${serviceId}/plans`)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch plans' }
    }
  }

  const UpdatePlan = async (serviceId: string, planId: string, body: CreatePlanDTO): Promise<ApiResult<{ data: Plan }>> => {
    try {
      const response = await api.patch(`/catalog/services/${serviceId}/plans/${planId}`, body)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to update plan' }
    }
  }

  const DeletePlan = async (serviceId: string, planId: string): Promise<ApiResult<null>> => {
    try {
      await api.delete(`/catalog/services/${serviceId}/plans/${planId}`)
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to delete plan' }
    }
  }

  return { CreatePlan, ListPlans, UpdatePlan, DeletePlan }
}
