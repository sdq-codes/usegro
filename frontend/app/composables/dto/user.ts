export interface UserOrganization {
  id: string
}

export interface LoggedInUserResponse {
  id: number
  email: string
  organizations?: UserOrganization[]
}
