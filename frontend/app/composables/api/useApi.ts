// composables/useApi.ts
import axios, { type AxiosInstance } from "axios"

// Access token lives in memory — cleared on page refresh (re-hydrated via refresh token)
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
    api = axios.create({
      baseURL: process.env.NUXT_PUBLIC_API_BASE || "http://usegro-production-alb-973426588.eu-west-1.elb.amazonaws.com/api/v1",
      timeout: 10000,
      headers: { "Content-Type": "application/json" },
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

    // 401 → refresh → retry, with queue to handle parallel failures
    api.interceptors.response.use(
      (response) => response,
      async (error) => {
        const original = error.config

        if (error.response?.status === 401 && !original._retry) {
          if (!localStorage.getItem('refresh_token')) return Promise.reject(error)
          original._retry = true

          if (isRefreshing) {
            // Queue this request until the ongoing refresh completes
            return new Promise<string>((resolve, reject) => {
              queue.push({ resolve, reject })
            }).then((token) => {
              original.headers.Authorization = `Bearer ${token}`
              return api(original)
            })
          }

          isRefreshing = true

          try {
            const { data } = await api.post("/authentication/refresh", {
              refresh_token: localStorage.getItem("refresh_token"),
            })
            accessToken = data.data.access_token
            localStorage.setItem("refresh_token", data.data.refresh_token)
            processQueue(null, accessToken)
            original.headers.Authorization = `Bearer ${accessToken}`
            return api(original)
          } catch (err) {
            processQueue(err, null)
            accessToken = null
            localStorage.removeItem("refresh_token")
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
