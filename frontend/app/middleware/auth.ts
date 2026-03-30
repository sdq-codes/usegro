import { navigateTo } from 'nuxt/app';

export default defineNuxtRouteMiddleware((to, _) => {
  if (process.server) return;

  // 'session' is a non-sensitive flag indicating an HttpOnly refresh cookie exists.
  // The actual security is enforced server-side; this is for UX routing only.
  const hasSession = !!localStorage.getItem("session");

  const publicPaths = ["/authentication/login", "/authentication/register", "/authentication/forgot-password", "/reset-password"]
  if (!hasSession && !publicPaths.includes(to.path)) {
    return navigateTo("/authentication/login");
  }

  if (hasSession && (to.path === "/authentication/login" || to.path === "/authentication/register")) {
    return navigateTo("/dashboard");
  }
})
