export type ApiResult<T, H> = {
  success: true; data: T; headers: H
} | { success: false; error: string }
