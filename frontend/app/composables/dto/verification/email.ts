export interface EmailVerificationOtpPayload {
  code: string
}

export interface EmailVerificationOtpResponse {
  id: number
  email: string
}
