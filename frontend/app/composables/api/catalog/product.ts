import { useApi } from '../useApi'
import type { CreateProductDTO, UpdateProductDTO, CatalogItem } from '@/composables/dto/catalog/product'

interface ApiResult<T> {
  success: boolean
  data?: T
  error?: string
}

export const useProductAPI = () => {
  const api = useApi()

  const CreateProduct = async (body: CreateProductDTO): Promise<ApiResult<{ data: CatalogItem }>> => {
    try {
      const response = await api.post('/catalog/products', body)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to create product' }
    }
  }

  const ListProducts = async (params?: { search?: string; status?: string; page?: number; limit?: number }): Promise<ApiResult<{ data: { data: CatalogItem[]; total: number; page: number; limit: number; total_pages: number } }>> => {
    try {
      const response = await api.get('/catalog/products', { params })
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to fetch products' }
    }
  }

  const GetProduct = async (productId: string): Promise<ApiResult<{ data: CatalogItem }>> => {
    try {
      const response = await api.get(`/catalog/products/${productId}`)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Product not found' }
    }
  }

  const UpdateProduct = async (productId: string, body: UpdateProductDTO): Promise<ApiResult<{ data: CatalogItem }>> => {
    try {
      const response = await api.patch(`/catalog/products/${productId}`, body)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to update product' }
    }
  }

  const DeleteProduct = async (productId: string): Promise<ApiResult<null>> => {
    try {
      await api.delete(`/catalog/products/${productId}`)
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Failed to delete product' }
    }
  }

  return { CreateProduct, ListProducts, GetProduct, UpdateProduct, DeleteProduct }
}
