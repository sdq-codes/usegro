export interface CreateCRMPayload {
  business_name: string
  full_name: string
  business_info?: string
}

export interface UpdateCRMPayload {
  business_name?: string
  full_name?: string
  business_info?: string
}

export interface CRMData {
  id: string
  business_name: string
  full_name: string
  business_info: string
}

export interface CRMApiResponse {
  data: CRMData
  response_message: string
  response_code: number
  request_id: string
}
