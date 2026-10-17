export type FoodCategory = "cooked_food" | "raw_food" | "bakery";
export type ListingStatus = "available" | "claimed" | "expired";

export interface FoodListing {
  food_listing_id: string;
  user_id: string;
  title: string;
  description?: string;
  latitude: number;
  longitude: number;
  food_type: FoodCategory;
  quantity: number;
  available_from: string;
  available_until: string;
  status: ListingStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateFoodListingRequest {
  title: string;
  description?: string;
  latitude: number;
  longitude: number;
  food_type: FoodCategory;
  quantity: number;
  available_from: string;
  available_until: string;
}
