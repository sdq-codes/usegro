export interface RegisterUserPayload {
  email: string
  password: string
}

export interface ForgotPasswordPayload {
  email: string
}

export interface ResetPasswordPayload {
  password: string
  confirm_password: string
  token: string | string[] | undefined
}

export interface RegisterUserResponse {
  id: number
  email: string
  token: string
}

export type ForgotPasswordResponse = object
