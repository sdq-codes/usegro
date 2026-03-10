import { useApi } from "../../useApi"
import type {LoggedInUserResponse} from "@/composables/dto/user";
import type {ApiResult} from "@/composables/helpers/types/api";
import type { RawAxiosResponseHeaders, AxiosHeaders } from "axios";
import type {FormVersionResponse, SubmissionPayload} from "@/composables/dto/customer/form/form";

export const useFormAPI = () => {
  const api = useApi()

  const FetchForm = async (): Promise<ApiResult<FormVersionResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.get<LoggedInUserResponse>("/crm/forms/customers/create")

      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Form not found"  }
    }
  }

  const CreateCustomer = async (formId: string, versionId: string, formData: SubmissionPayload): Promise<ApiResult<null, RawAxiosResponseHeaders>> => {
    try {
      const response = await api.post(
        `/forms/${formId}/version/${versionId}/submission`,
        formData
      )

      return { success: true, data: response.data, headers: response.headers}
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Customer creation failed"  }
    }
  }

  return {
    FetchForm,
    CreateCustomer
  }
}
