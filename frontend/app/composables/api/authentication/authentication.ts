import { useApi, setAccessToken } from "../useApi"
import type {
  ForgotPasswordPayload,
  ForgotPasswordResponse,
  RegisterUserPayload,
  RegisterUserResponse, ResetPasswordPayload
} from "@/composables/dto/auth";
import type {ApiResult} from "@/composables/helpers/types/api";
import type { RawAxiosResponseHeaders, AxiosHeaders } from "axios";

export const useAuthentication = () => {
  const api = useApi()

  const CheckEmailExists = async (email: string): Promise<{ exists: boolean; errors?: string[] }> => {
    try {
      const response = await api.post("/base/authentication/email/exist", { email })
      const exists = response.data?.responseMessage === "Email already exists"
      return { exists }
    } catch (error: unknown) {
      const body = error.response?.data
      const errors: string[] = body?.errors?.map((e: { message: string }) => e.message) ?? [
        body?.response_message ?? "Could not verify email"
      ]
      return { exists: false, errors }
    }
  }

  const RegisterUser = async (
    data: RegisterUserPayload
  ): Promise<ApiResult<RegisterUserResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post<RegisterUserResponse>("/base/authentication/register", {
        email: data.email,
        password: data.password,
      })

      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Registration failed" }
    }
  }

  const LoginUser = async (
    data: RegisterUserPayload
  ): Promise<ApiResult<RegisterUserResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post("/base/authentication/login", {
        email: data.email,
        password: data.password,
      })

      return { success: true, data: response.data, headers: response.headers}
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Login failed" }
    }
  }

  const Logout = async (): Promise<ApiResult<unknown, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      // No body — the HttpOnly cookie is sent automatically and cleared by the server
      const response = await api.post("/base/authentication/logout")
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Logout failed" }
    } finally {
      setAccessToken(null)
      localStorage.removeItem("session")
    }
  }

  const RefreshToken = async (): Promise<boolean> => {
    // No body — the HttpOnly cookie is sent automatically
    if (!localStorage.getItem('session')) return false

    try {
      const { data } = await api.post("/base/authentication/refresh")
      setAccessToken(data.data.access_token)
      return true
    } catch {
      localStorage.removeItem("session")
      return false
    }
  }

  const ForgotPassword = async (
    data: ForgotPasswordPayload
  ): Promise<ApiResult<ForgotPasswordResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post("/base/authentication/forgot-password", {
        email: data.email,
      })

      return { success: true, data: response.data, headers: response.headers}
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Forgot password failed" }
    }
  }

  const ResetPassword = async (
    data: ResetPasswordPayload
  ): Promise<ApiResult<ForgotPasswordResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post("/base/authentication/reset-password", {
        password: data.password,
        confirm_password: data.confirm_password,
        token: data.token,
      })
      return { success: true, data: response.data, headers: response.headers}
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Reset password failed" }
    }
  }

  const RequestEmailCode = async (email: string): Promise<ApiResult<unknown, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post("/base/authentication/email-code/request", { email })
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Failed to send code" }
    }
  }

  const VerifyEmailCode = async (
    email: string,
    code: string
  ): Promise<ApiResult<RegisterUserResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    try {
      const response = await api.post("/base/authentication/email-code/verify", { email, code })
      return { success: true, data: response.data, headers: response.headers }
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.response_message || "Invalid or expired code" }
    }
  }

  return {
    CheckEmailExists,
    RegisterUser,
    LoginUser,
    Logout,
    RefreshToken,
    ForgotPassword,
    ResetPassword,
    RequestEmailCode,
    VerifyEmailCode,
  }
}
