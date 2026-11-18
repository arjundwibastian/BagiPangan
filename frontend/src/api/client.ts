import type { ApiResponse } from "../types/auth";

const API_ORIGIN = (import.meta.env?.VITE_API_BASE_URL ?? "").replace(/\/$/, "");
const BASE_URL = `${API_ORIGIN}/api/v1`;

export class ApiError extends Error {
  code: string;
  status: number;

  constructor(code: string, message: string, status = 0) {
    super(message);
    this.code = code;
    this.status = status;
  }
}

export async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const token = localStorage.getItem("access_token");

  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  });

  const body: (Partial<ApiResponse<T>> & { error?: string }) | null =
    await res.json().catch(() => null);

  if (!res.ok || body?.responseCode !== "00") {
    if (res.status === 401 && token) {
      for (const key of ["access_token", "refresh_token", "user"]) {
        localStorage.removeItem(key);
      }
      window.location.reload();
    }
    throw new ApiError(
      body?.responseCode ?? String(res.status),
      body?.responseMessage ?? body?.error ?? `Server returned HTTP ${res.status}. Please try again.`,
      res.status,
    );
  }

  return body.responseData as T;
}
