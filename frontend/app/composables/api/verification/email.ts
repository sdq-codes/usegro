import { useApi } from "../useApi"
import type {LoggedInUserResponse} from "@/composables/dto/user";
import type {ApiResult} from "@/composables/helpers/types/api";
import type { RawAxiosResponseHeaders, AxiosHeaders } from "axios";
import type {EmailVerificationOtpPayload, EmailVerificationOtpResponse} from "@/composables/dto/verification/email";

export const useVerificationEmailAPI = () => {
  const api = useApi()

  const ResendVerificationEmail = async (): Promise<ApiResult<LoggedInUserResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.get<LoggedInUserResponse>("/verification/email/resend")

      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Failed to resend verification code" }
    }
  }

  const VerifyEmail = async (
    data: EmailVerificationOtpPayload
  ): Promise<ApiResult<EmailVerificationOtpResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post("/verification/email", {
        token_hash: data.code
      })

      return { success: true, data: response.data, headers: response.headers}
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to verify email" }
    }
  }
  return {
    ResendVerificationEmail,
    VerifyEmail
  }
}
