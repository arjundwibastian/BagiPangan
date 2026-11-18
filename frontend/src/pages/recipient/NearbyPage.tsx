import { useEffect, useState } from "react";
import { listFoodRequests } from "../../api/requests";
import { listNearbyFoodListings } from "../../api/food";
import { createClaim } from "../../api/claims";
import { ApiError } from "../../api/client";
import type { FoodRequest } from "../../types/request";
import type { FoodListing } from "../../types/food";
import type { FoodClaim } from "../../types/claim";
import type { User } from "../../types/auth";

interface NearbyPageProps {
  user: User;
}

export function NearbyPage({ user }: NearbyPageProps) {
  const [requests, setRequests] = useState<FoodRequest[]>([]);
  const [requestId, setRequestId] = useState("");
  const [listings, setListings] = useState<FoodListing[]>([]);
  const [quantities, setQuantities] = useState<Record<string, string>>({});
  const [claimed, setClaimed] = useState<FoodClaim | null>(null);
  const [claimingId, setClaimingId] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [searched, setSearched] = useState(false);

  useEffect(() => {
    let cancelled = false;

    listFoodRequests(user.user_id)
      .then((data) => {
        if (cancelled) return;
        const searchable = (data ?? []).filter((r) => r.status === "searching");
        setRequests(searchable);
        if (searchable.length > 0) {
          setRequestId((current) => current || searchable[0].request_id);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(
            err instanceof ApiError ? err.message : "Cannot reach the server.",
          );
        }
      });

    return () => {
      cancelled = true;
    };
  }, [user.user_id]);

  async function handleSearch() {
    if (!requestId) return;
    setLoading(true);
    setError("");
    setSearched(true);
    setClaimed(null);

    try {
      const data = await listNearbyFoodListings(requestId);
      setListings(data ?? []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Cannot reach the server.");
    } finally {
      setLoading(false);
    }
  }

  async function handleClaim(listing: FoodListing) {
    if (claimingId || loading || !requestId) return;
    setError("");
    setClaimed(null);
    const quantity = Number(quantities[listing.food_listing_id] ?? "1");
    if (!Number.isInteger(quantity) || quantity < 1 || quantity > listing.quantity) {
      setError("Enter a whole quantity between 1 and the available stock.");
      return;
    }
    setClaimingId(listing.food_listing_id);

    try {
      const claim = await createClaim({
        request_id: requestId,
        food_listing_id: listing.food_listing_id,
        claimed_quantity: quantity,
      });
      setClaimed(claim);
      // One request can back only one claim. Its old search results are stale.
      const remaining = requests.filter((request) => request.request_id !== requestId);
      setRequests(remaining);
      setRequestId(remaining[0]?.request_id ?? "");
      setListings([]);
      setSearched(false);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Cannot reach the server.",
      );
    } finally {
      setClaimingId("");
    }
  }

  return (
    <section className="page">
      <div className="page-header">
        <h2 className="page-title">Browse Nearby</h2>
        <p className="page-subtitle">
          Find active listings around one of your requests.
        </p>
      </div>

      {requests.length === 0 && (
        <p className="muted">
          Create a request first — nearby search needs one.
        </p>
      )}

      {requests.length > 0 && (
        <div className="row-actions">
          <select
            className="field-input"
            value={requestId}
            onChange={(e) => {
              setRequestId(e.target.value);
              setListings([]);
              setSearched(false);
              setClaimed(null);
            }}
            disabled={loading || claimingId !== ""}
          >
            {requests.map((request) => (
              <option key={request.request_id} value={request.request_id}>
                {request.request_id.slice(0, 8)} · radius {request.radius_km} km
              </option>
            ))}
          </select>
          <button
            type="button"
            className="btn-primary"
            onClick={handleSearch}
            disabled={loading || claimingId !== ""}
          >
            {loading ? "Searching..." : "Search"}
          </button>
        </div>
      )}

      {error && <p className="auth-error">{error}</p>}

      {claimed && (
        <p className="success-note">
          Claimed! Pickup code: <strong>{claimed.claim_code}</strong>
        </p>
      )}

      {searched && !loading && !error && listings.length === 0 && (
        <p className="muted">No listings found within that request's radius.</p>
      )}

      <div className="grid">
        {listings.map((listing) => (
          <article className="card" key={listing.food_listing_id}>
            <span className={`badge badge--${listing.status}`}>
              {listing.status}
            </span>
            <h3 className="card-title">{listing.title}</h3>
            {listing.description && (
              <p className="card-meta">{listing.description}</p>
            )}
            <p className="card-meta">Type: {listing.food_type}</p>
            <p className="card-meta">Quantity: {listing.quantity}</p>
            <p className="card-meta">
              At {listing.latitude}, {listing.longitude}
            </p>

            <div className="claim-row">
              <input
                className="field-input"
                type="number"
                min={1}
                max={listing.quantity}
                value={quantities[listing.food_listing_id] ?? "1"}
                onChange={(e) =>
                  setQuantities((current) => ({
                    ...current,
                    [listing.food_listing_id]: e.target.value,
                  }))
                }
              />
              <button
                type="button"
                className="btn-primary"
                onClick={() => handleClaim(listing)}
                disabled={claimingId !== ""}
              >
                {claimingId === listing.food_listing_id
                  ? "Claiming..."
                  : "Claim"}
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
