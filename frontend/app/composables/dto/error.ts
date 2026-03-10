export interface ApiErrorResponse {
  error: {
    code: string
    message: string
    status: number
    timestamp: string
    requestId: string
    details?: []
  }
}
