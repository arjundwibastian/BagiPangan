package dto

type CreateClaimRequest struct {
	RequestID       string `json:"request_id"`
	FoodListingID   string `json:"food_listing_id"`
	ClaimedQuantity int32  `json:"claimed_quantity"`
}

type VerifyListingPickupRequest struct {
	FoodListingID string `json:"food_listing_id"`
	ClaimCode     string `json:"claim_code"`
}

type VerifyPickupRequest struct {
	ClaimCode string `json:"claim_code"`
}
