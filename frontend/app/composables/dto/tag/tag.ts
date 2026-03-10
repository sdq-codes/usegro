export interface CreateCrmTagResponse {
  tag: string
}

export interface FetchCrmTagResponse {
  PK: number
  SK: string
  tag: string
  crmId: string
  status: string
}
