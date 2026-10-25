import { request } from "./client";
import type {
  CreateClaimRequest,
  FoodClaim,
  VerifyPickupRequest,
} from "../types/claim";

export function createClaim(payload: CreateClaimRequest): Promise<FoodClaim> {
  return request<FoodClaim>("/claims", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function getClaim(id: string): Promise<FoodClaim> {
  return request<FoodClaim>(`/claims/${id}`);
}

export function verifyListingPickup(listingId: string, code: string): Promise<FoodClaim> {
  return request<FoodClaim>("/claims/verify-pickup", {
    method: "POST",
    body: JSON.stringify({ food_listing_id: listingId, claim_code: code }),
  });
}

export function verifyPickup(
  id: string,
  payload: VerifyPickupRequest,
): Promise<FoodClaim> {
  return request<FoodClaim>(`/claims/${id}/verify-pickup`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function cancelClaim(id: string): Promise<FoodClaim> {
  return request<FoodClaim>(`/claims/${id}/cancel`, {
    method: "POST",
  });
}
