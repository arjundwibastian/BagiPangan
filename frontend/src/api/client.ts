import type { ApiResponse } from "../types/auth";

const BASE_URL = "/api/v1";

export class ApiError extends Error {
  code: string;

  constructor(code: string, message: string) {
    super(message);
    this.code = code;
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

  const body: ApiResponse<T> = await res.json();

  if (body.responseCode !== "00") {
    if (body.responseCode === "41" && token) {
      localStorage.clear();
      window.location.reload();
    }
    throw new ApiError(body.responseCode, body.responseMessage);
  }

  return body.responseData;
}
