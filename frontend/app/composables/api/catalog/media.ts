import { useApi } from '../useApi'
import type { CatalogItemMedia } from '@/composables/dto/catalog/product'

interface ApiResult<T> {
  success: boolean
  data?: T
  error?: string
}

interface PresignResponse {
  key: string
  upload_url: string
  public_url: string
}

export const useCatalogMediaAPI = () => {
  const api = useApi()

  /**
   * Get presigned PUT URLs from the server, upload files directly to R2, return keys.
   */
  const presignAndUpload = async (files: File[]): Promise<{ keys: string[]; error?: string }> => {
    if (files.length === 0) return { keys: [] }

    try {
      // Step 1: get presigned URLs
      const presignRes = await api.post('/catalog/media/presign', {
        files: files.map(f => ({ name: f.name, content_type: f.type })),
      })

      const presigned = presignRes.data.data as PresignResponse[]

      // Step 2: upload each file directly to R2
      await Promise.all(
        presigned.map((p, i) =>
          fetch(p.upload_url, {
            method: 'PUT',
            headers: { 'Content-Type': files[i].type },
            body: files[i],
          }).then(r => {
            if (!r.ok) throw new Error(`Failed to upload ${files[i].name}`)
          })
        )
      )

      return { keys: presigned.map(p => p.key) }
    } catch (err: any) {
      return { keys: [], error: err?.response?.data?.response_message || err?.message || 'Upload failed' }
    }
  }

  const DeleteMedia = async (mediaId: string): Promise<ApiResult<null>> => {
    try {
      await api.delete(`/catalog/media/${mediaId}`)
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error?.response?.data?.response_message || 'Delete failed' }
    }
  }

  return { presignAndUpload, DeleteMedia }
}
