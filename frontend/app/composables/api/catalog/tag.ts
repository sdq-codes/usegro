import { useApi } from '../useApi'

interface ApiResult<T> {
  success: boolean
  data?: T
  error?: string
}

export interface CatalogTag {
  id: string
  crm_id: string
  name: string
  created_at: string
  updated_at: string
}

export const useCatalogTagAPI = () => {
  const api = useApi()

  const CreateTag = async (name: string): Promise<ApiResult<{ data: CatalogTag }>> => {
    try {
      const response = await api.post('/catalog/tags', { name })
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to create tag' }
    }
  }

  const ListTags = async (): Promise<ApiResult<{ data: CatalogTag[] }>> => {
    try {
      const response = await api.get('/catalog/tags')
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch tags' }
    }
  }

  const DeleteTag = async (tagId: string): Promise<ApiResult<null>> => {
    try {
      await api.delete(`/catalog/tags/${tagId}`)
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to delete tag' }
    }
  }

  return { CreateTag, ListTags, DeleteTag }
}
