import { request } from "./client";
import type { CreateFoodListingRequest, FoodListing } from "../types/food";

export function listFoodListings(): Promise<FoodListing[]> {
  return request<FoodListing[]>("/food-listings");
}

export function listOwnFoodListings(): Promise<FoodListing[]> {
  return request<FoodListing[]>("/food-listings/mine");
}

export function getFoodListing(id: string): Promise<FoodListing> {
  return request<FoodListing>(`/food-listings/${id}`);
}

export function listNearbyFoodListings(
  requestId: string,
): Promise<FoodListing[]> {
  return request<FoodListing[]>(
    `/food-listings/nearby?request_id=${encodeURIComponent(requestId)}`,
  );
}

export function createFoodListing(
  payload: CreateFoodListingRequest,
): Promise<FoodListing> {
  return request<FoodListing>("/food-listings", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}
