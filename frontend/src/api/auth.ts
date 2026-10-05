import { request } from "./client";
import type { AuthResult, LoginRequest } from "../types/auth";

export function login(payload: LoginRequest): Promise<AuthResult> {
  return request<AuthResult>("/auth/login", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}
