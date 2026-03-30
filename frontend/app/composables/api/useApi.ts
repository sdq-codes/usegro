// composables/useApi.ts
import axios, { type AxiosInstance } from "axios"

// Access token lives in memory — cleared on page refresh (re-hydrated via refresh cookie)
let accessToken: string | null = null

let isRefreshing = false
let queue: Array<{ resolve: (token: string) => void; reject: (err: unknown) => void }> = []

function processQueue(error: unknown, token: string | null) {
  queue.forEach(p => (error ? p.reject(error) : p.resolve(token!)))
  queue = []
}

export function setAccessToken(token: string | null) {
  accessToken = token
}

export function getAccessToken(): string | null {
  return accessToken
}

let api: AxiosInstance

export const useApi = () => {
  if (!api) {
    const config = useRuntimeConfig()
    const baseURL = (config.public.apiBase as string) || "http://localhost/api/v1"
    api = axios.create({
      baseURL,
      timeout: 10000,
      headers: { "Content-Type": "application/json" },
      withCredentials: true, // send HttpOnly refresh_token cookie automatically
    })

    // Attach access token + CRM-ID on every request
    api.interceptors.request.use((config) => {
      if (accessToken) {
        config.headers.Authorization = `Bearer ${accessToken}`
      }
      const crmId = localStorage.getItem("crm-id")
      if (crmId) {
        config.headers["X-CRM-ID"] = crmId
      }
      return config
    })

    // 401 → refresh via cookie → retry, with queue to handle parallel failures
    api.interceptors.response.use(
      (response) => response,
      async (error) => {
        const original = error.config

        if (error.response?.status === 401 && !original._retry) {
          // Only attempt refresh if we believe a session exists
          if (!localStorage.getItem('session')) return Promise.reject(error)
          original._retry = true

          if (isRefreshing) {
            return new Promise<string>((resolve, reject) => {
              queue.push({ resolve, reject })
            }).then((token) => {
              original.headers.Authorization = `Bearer ${token}`
              return api(original)
            })
          }

          isRefreshing = true

          try {
            // No body needed — the HttpOnly cookie is sent automatically
            const { data } = await api.post("/base/authentication/refresh")
            accessToken = data.data.access_token
            processQueue(null, accessToken)
            original.headers.Authorization = `Bearer ${accessToken}`
            return api(original)
          } catch (err) {
            processQueue(err, null)
            accessToken = null
            localStorage.removeItem("session")
            window.location.href = "/authentication/login"
            return Promise.reject(err)
          } finally {
            isRefreshing = false
          }
        }

        return Promise.reject(error)
      }
    )
  }

  return api
}
