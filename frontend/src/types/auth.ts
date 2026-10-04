export type Role = "donor" | "recipient";

export interface User {
  user_id: string;
  name: string;
  email: string;
  phone: string;
  role: Role;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResult {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

export interface ApiResponse<T> {
  responseCode: string;
  responseMessage: string;
  responseData: T;
}
