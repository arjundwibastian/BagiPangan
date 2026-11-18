import { useEffect, useState } from "react";
import { listOwnFoodListings } from "../../api/food";
import { ApiError } from "../../api/client";
import type { FoodListing } from "../../types/food";
import type { User } from "../../types/auth";

interface MyListingsPageProps {
  user: User;
}

export function MyListingsPage({ user }: MyListingsPageProps) {
  const [listings, setListings] = useState<FoodListing[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let cancelled = false;
    listOwnFoodListings()
      .then((data) => {
        if (cancelled) return;
        setListings(data.filter((item) => item.user_id === user.user_id));
      })
      .catch((err) => {
        if (cancelled) return;
        setError(err instanceof ApiError ? err.message : "Cannot reach the server.");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [user.user_id, reloadKey]);

  return (
    <section className="page">
      <div className="page-header">
        <h2 className="page-title">My Listings</h2>
        <p className="page-subtitle">
          Food you have posted for others to claim.
        </p>
      </div>

      <div className="row-actions">
        <button
          type="button"
          className="btn-secondary"
          onClick={() => {
            setLoading(true);
            setError("");
            setListings([]);
            setReloadKey((k) => k + 1);
          }}
          disabled={loading}
        >
          Refresh
        </button>
      </div>

      {loading && <p className="muted">Loading...</p>}
      {error && <p className="auth-error">{error}</p>}

      {!loading && !error && listings.length === 0 && (
        <p className="muted">
          No listings yet. Create one from the "Create Listing" tab.
        </p>
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
              {new Date(listing.available_from).toLocaleString()} —{" "}
              {new Date(listing.available_until).toLocaleString()}
            </p>
            <p className="card-meta">
              At {listing.latitude}, {listing.longitude}
            </p>
          </article>
        ))}
      </div>
    </section>
  );
}
