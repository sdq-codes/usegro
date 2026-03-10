import { navigateTo } from 'nuxt/app';

export default defineNuxtRouteMiddleware((to, _) => {
  if (process.server) return;

  const hasSession = !!localStorage.getItem("refresh_token");

  const publicPaths = ["/authentication/login", "/authentication/register", "/authentication/forgot-password", "/reset-password"]
  if (!hasSession && !publicPaths.includes(to.path)) {
    return navigateTo("/authentication/login");
  }

  if (hasSession && (to.path === "/authentication/login" || to.path === "/authentication/register")) {
    return navigateTo("/dashboard");
  }
})
