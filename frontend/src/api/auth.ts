import { request } from "./client";
import type {
  AuthResult,
  LoginRequest,
  RegisterRequest,
  User,
} from "../types/auth";

export function login(payload: LoginRequest): Promise<AuthResult> {
  return request<AuthResult>("/auth/login", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function register(payload: RegisterRequest): Promise<User> {
  return request<User>("/auth/register", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}