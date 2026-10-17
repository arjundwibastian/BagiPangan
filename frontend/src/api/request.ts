import { request } from "./client";
import type { CreateFoodRequestRequest, FoodRequest } from "../types/request";

export function createFoodRequest(
  payload: CreateFoodRequestRequest,
): Promise<FoodRequest> {
  return request<FoodRequest>("/food-requests", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function getFoodRequest(id: string): Promise<FoodRequest> {
  return request<FoodRequest>(`/food-requests/${id}`);
}

export function listFoodRequests(userId: string): Promise<FoodRequest[]> {
  return request<FoodRequest[]>(`/users/${userId}/food-requests`);
}

export function cancelFoodRequest(id: string): Promise<FoodRequest> {
  return request<FoodRequest>(`/food-requests/${id}/cancel`, {
    method: "PATCH",
  });
}
