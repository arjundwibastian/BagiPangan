export type RequestStatus =
  | "searching"
  | "claimed"
  | "completed"
  | "cancelled"
  | "expired";

export interface FoodRequest {
  request_id: string;
  user_id: string;
  latitude: number;
  longitude: number;
  radius_km: number;
  status: RequestStatus;
  claim_id?: string;
  expires_at: string;
  created_at: string;
  updated_at: string;
}

export interface CreateFoodRequestRequest {
  latitude: number;
  longitude: number;
  radius_km: number;
}
