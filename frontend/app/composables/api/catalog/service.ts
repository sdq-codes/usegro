import { useApi } from '../useApi'
import type { CreateServiceDTO, UpdateServiceDTO } from '@/composables/dto/catalog/service'
import type { CatalogItem } from '@/composables/dto/catalog/product'

interface ApiResult<T> {
  success: boolean
  data?: T
  error?: string
}

export const useServiceAPI = () => {
  const api = useApi()

  const CreateService = async (body: CreateServiceDTO): Promise<ApiResult<{ data: CatalogItem }>> => {
    try {
      const response = await api.post('/catalog/services', body)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to create service' }
    }
  }

  const ListServices = async (params?: { search?: string; status?: string }): Promise<ApiResult<{ data: CatalogItem[] }>> => {
    try {
      const response = await api.get('/catalog/services', { params })
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch services' }
    }
  }

  const GetService = async (serviceId: string): Promise<ApiResult<{ data: CatalogItem }>> => {
    try {
      const response = await api.get(`/catalog/services/${serviceId}`)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Service not found' }
    }
  }

  const UpdateService = async (serviceId: string, body: UpdateServiceDTO): Promise<ApiResult<{ data: CatalogItem }>> => {
    try {
      const response = await api.patch(`/catalog/services/${serviceId}`, body)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to update service' }
    }
  }

  const DeleteService = async (serviceId: string): Promise<ApiResult<null>> => {
    try {
      await api.delete(`/catalog/services/${serviceId}`)
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to delete service' }
    }
  }

  return { CreateService, ListServices, GetService, UpdateService, DeleteService }
}
