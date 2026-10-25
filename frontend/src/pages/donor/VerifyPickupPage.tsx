import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { verifyListingPickup } from "../../api/claims";
import { listOwnFoodListings } from "../../api/food";
import { ApiError } from "../../api/client";
import type { User } from "../../types/auth";
import type { FoodListing } from "../../types/food";

function PickupForm({ listing }: { listing: FoodListing }) {
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [verifying, setVerifying] = useState(false);

  async function handleVerify(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (verifying) return;
    setError("");
    setSuccess("");
    if (!/^\d{6}$/.test(code)) {
      setError("Enter the recipient's 6-digit pickup code.");
      return;
    }
    setVerifying(true);
    try {
      const claim = await verifyListingPickup(listing.food_listing_id, code);
      setSuccess(`Pickup confirmed for ${claim.claimed_quantity} item(s).`);
      setCode("");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Cannot reach the server.");
    } finally {
      setVerifying(false);
    }
  }

  return (
    <article className="card">
      <span className={`badge badge--${listing.status}`}>
        {listing.status === "claimed" ? "Fully reserved" : "Available"}
      </span>
      <h3 className="card-title">{listing.title}</h3>
      {listing.description && <p className="card-meta">{listing.description}</p>}
      <p className="card-meta">Unreserved quantity: {listing.quantity}</p>
      <p className="card-meta">
        Pickup: {new Date(listing.available_from).toLocaleString()} —{" "}
        {new Date(listing.available_until).toLocaleString()}
      </p>
      <form className="pickup-form" onSubmit={handleVerify}>
        <label className="field">
          <span className="field-label">Recipient's pickup code</span>
          <input
            className="field-input"
            value={code}
            onChange={(event) => {
              setCode(event.target.value.replace(/\D/g, "").slice(0, 6));
              setError("");
              setSuccess("");
            }}
            inputMode="numeric"
            pattern="[0-9]{6}"
            maxLength={6}
            autoComplete="off"
            placeholder="123456"
            disabled={verifying}
            required
          />
        </label>
        {error && <p className="auth-error" role="alert">{error}</p>}
        {success && <p className="success-note" role="status">{success}</p>}
        <button className="btn-primary" type="submit" disabled={verifying || code.length !== 6}>
          {verifying ? "Verifying..." : "Verify pickup"}
        </button>
      </form>
    </article>
  );
}

export function VerifyPickupPage({ user }: { user: User }) {
  const [listings, setListings] = useState<FoodListing[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let cancelled = false;
    listOwnFoodListings()
      .then((data) => {
        if (cancelled) return;
        // Fully reserved listings still need pickup verification.
        setListings(data.filter((listing) =>
          listing.user_id === user.user_id &&
          (listing.status === "available" || listing.status === "claimed")
        ));
      })
      .catch((err) => {
        if (cancelled) return;
        setError(err instanceof ApiError ? err.message : "Cannot reach the server.");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [user.user_id, reloadKey]);

  function refresh() {
    setLoading(true);
    setError("");
    setListings([]);
    setReloadKey((key) => key + 1);
  }

  return (
    <section className="page">
      <div className="page-header">
        <h2 className="page-title">Verify Pickup</h2>
        <p className="page-subtitle">
          Find the recipient's listing below and enter their 6-digit pickup code.
          No claim ID is needed.
        </p>
      </div>
      <div className="row-actions">
        <button type="button" className="btn-secondary" onClick={refresh} disabled={loading}>
          Refresh listings
        </button>
      </div>
      {loading && <p className="muted" role="status">Loading active listings...</p>}
      {error && <p className="auth-error" role="alert">{error}</p>}
      {!loading && !error && listings.length === 0 && (
        <p className="muted">No active listings. Create one from the "Create Listing" tab.</p>
      )}
      {!loading && !error && (
        <div className="grid">
          {listings.map((listing) => <PickupForm key={listing.food_listing_id} listing={listing} />)}
        </div>
      )}
    </section>
  );
}
