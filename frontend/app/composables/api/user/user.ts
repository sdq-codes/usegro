import { useApi } from "../useApi"
import type {LoggedInUserResponse} from "@/composables/dto/user";
import type {ApiResult} from "@/composables/helpers/types/api";
import type { RawAxiosResponseHeaders, AxiosHeaders } from "axios";
import { navigateTo } from "nuxt/app";
import { verifications } from "@/composables/helpers/format/verifications";

const ONBOARDING_PATH = "/onboarding";
const VERIFICATION_PATH = "/verification/email";

export const useUserAPI = () => {
  const api = useApi()

  const FetchLoggedInUser = async (): Promise<ApiResult<LoggedInUserResponse, RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)>> => {
    let response;

    try {
      response = await api.get<LoggedInUserResponse>("/base/user")
    } catch (error: unknown) {
      return { success: false, error: error.response?.data?.message || "Failed to fetch user" }
    }

    const currentPath = process.client ? window.location.pathname : "";

    // 1. Email not verified → verification page
    const userVerifications = verifications(response.data?.data?.verifications ?? []);
    if (userVerifications?.email !== "VERIFIED") {
      if (currentPath !== VERIFICATION_PATH) {
        await navigateTo(VERIFICATION_PATH);
      }
      return { success: true, data: response.data, headers: response.headers }
    }

    // 2. Email verified → ensure crm-id is set in localStorage
    if (process.client && !localStorage.getItem("crm-id")) {
      const organizations = response.data?.data?.organizations;
      if (Array.isArray(organizations) && organizations.length > 0) {
        localStorage.setItem("crm-id", organizations[0].id);
      } else if (currentPath !== ONBOARDING_PATH) {
        await navigateTo(ONBOARDING_PATH);
        return { success: true, data: response.data, headers: response.headers }
      }
    }

    return { success: true, data: response.data, headers: response.headers }
  }

  return {
    FetchLoggedInUser,
  }
}
