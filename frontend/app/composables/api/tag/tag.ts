import { useApi } from "../useApi"
import type {ApiResult} from "@/composables/helpers/types/api";
import type { RawAxiosResponseHeaders, AxiosHeaders } from "axios";
import type {CreateCrmTagResponse, FetchCrmTagResponse} from "@/composables/dto/tag/tag";

export const useCrmTagsAPI = () => {
  const api = useApi()

  const CreateTag = async (
    data: CreateCrmTagResponse
  ): Promise<ApiResult<Array<FetchCrmTagResponse>, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post("/crm/tags", {
        tag: data.tag
      })
      return { success: true, data: response.data, headers: response.headers}
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to verify email" }
    }
  }

  const FetchCRMTag = async (): Promise<ApiResult<Array<FetchCrmTagResponse>, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.get("/crm/tags")
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Failed to resend verification code" }
    }
  }
  return {
    CreateTag,
    FetchCRMTag
  }
}
