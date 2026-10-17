import { request } from "./client";
import type { UpdateProfileRequest, User } from "../types/auth";

export function getMe(): Promise<User> {
  return request<User>("/users/me");
}

export function updateMe(payload: UpdateProfileRequest): Promise<User> {
  return request<User>("/users/me", {
    method: "PATCH",
    body: JSON.stringify(payload),
  });
}
