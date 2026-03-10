// Re-hydrate the in-memory access token on every page load using the stored refresh token
import { useAuthentication } from "@/composables/api/authentication/authentication"

export default defineNuxtPlugin(async () => {
  await useAuthentication().RefreshToken()
})
