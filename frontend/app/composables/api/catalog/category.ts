import { useApi } from '../useApi'
import type { Category, CreateCategoryDTO, UpdateCategoryDTO, StandardCategory, StandardCategorySearchResult, StandardAttribute } from '@/composables/dto/catalog/category'

interface ApiResult<T> {
  success: boolean
  data?: T
  error?: string
}

export const useCategoryAPI = () => {
  const api = useApi()

  const CreateCategory = async (body: CreateCategoryDTO): Promise<ApiResult<{ data: Category }>> => {
    try {
      const response = await api.post('/catalog/categories', body)
      return { success: true, data: response.data }
    } catch (error: unknown) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to create category' }
    }
  }

  const ListCategories = async (): Promise<ApiResult<{ data: Category[] }>> => {
    try {
      const response = await api.get('/catalog/categories')
      return { success: true, data: response.data }
    } catch (error: unknown) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch categories' }
    }
  }

  const UpdateCategory = async (categoryId: string, body: UpdateCategoryDTO): Promise<ApiResult<{ data: Category }>> => {
    try {
      const response = await api.patch(`/catalog/categories/${categoryId}`, body)
      return { success: true, data: response.data }
    } catch (error: unknown) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to update category' }
    }
  }

  const DeleteCategory = async (categoryId: string): Promise<ApiResult<null>> => {
    try {
      await api.delete(`/catalog/categories/${categoryId}`)
      return { success: true }
    } catch (error: unknown) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to delete category' }
    }
  }

  return {
    CreateCategory,
    ListCategories,
    UpdateCategory,
    DeleteCategory,
  }
}

export const useStandardCategoryAPI = () => {
  const api = useApi()

  const ListRootCategories = async (): Promise<ApiResult<{ data: StandardCategory[] }>> => {
    try {
      const response = await api.get('/catalog/standard-categories')
      return { success: true, data: response.data }
    } catch (error: unknown) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch categories' }
    }
  }

  const ListChildren = async (parentId: string): Promise<ApiResult<{ data: StandardCategory[] }>> => {
    try {
      const response = await api.get(`/catalog/standard-categories/${parentId}/children`)
      return { success: true, data: response.data }
    } catch (error: unknown) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch subcategories' }
    }
  }

  const GetCategory = async (id: string): Promise<ApiResult<{ data: StandardCategory }>> => {
    try {
      const response = await api.get(`/catalog/standard-categories/${id}`)
      return { success: true, data: response.data }
    } catch (error: unknown) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch category' }
    }
  }

  const SearchCategories = async (q: string): Promise<ApiResult<{ data: StandardCategorySearchResult[] }>> => {
    try {
      const response = await api.get(`/catalog/standard-categories/search?q=${encodeURIComponent(q)}`)
      return { success: true, data: response.data }
    } catch (error: unknown) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to search categories' }
    }
  }

  return { ListRootCategories, ListChildren, GetCategory, SearchCategories }
}
