export type ClaimStatus = "waiting_for_pickup" | "picked_up" | "cancelled";

export interface FoodClaim {
  claim_id: string;
  request_id: string;
  food_listing_id: string;
  claim_code: string;
  code_expires_at: string;
  claimed_at?: string;
  claimed_quantity: number;
  status: ClaimStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateClaimRequest {
  request_id: string;
  food_listing_id: string;
  claimed_quantity: number;
}

export interface VerifyPickupRequest {
  claim_code: string;
}
