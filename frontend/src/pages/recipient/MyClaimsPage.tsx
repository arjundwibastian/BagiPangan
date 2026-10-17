import { useEffect, useState } from "react";
import { cancelClaim, getClaim } from "../../api/claims";
import { ApiError } from "../../api/client";
import { getClaimIds, removeClaimId } from "../../lib/claims";
import type { FoodClaim } from "../../types/claim";
import type { User } from "../../types/auth";

interface MyClaimsPageProps {
  user: User;
}

export function MyClaimsPage({ user }: MyClaimsPageProps) {
  const [claims, setClaims] = useState<FoodClaim[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError("");

    const ids = getClaimIds(user.user_id);
    if (ids.length === 0) {
      setClaims([]);
      setLoading(false);
      return;
    }

    Promise.all(
      ids.map((id) =>
        getClaim(id).catch((err) => {
          if (
            err instanceof ApiError &&
            (err.code === "01" || err.code === "42")
          ) {
            removeClaimId(user.user_id, id);
            return null;
          }
          throw err;
        }),
      ),
    )
      .then((results) => {
        if (cancelled) return;
        setClaims(results.filter((c): c is FoodClaim => c !== null));
      })
      .catch((err) => {
        if (cancelled) return;
        setError(
          err instanceof ApiError ? err.message : "Cannot reach the server.",
        );
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [user.user_id, reloadKey]);

  async function handleCancel(id: string) {
    setError("");
    try {
      await cancelClaim(id);
      setReloadKey((k) => k + 1);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Cannot reach the server.",
      );
    }
  }

  return (
    <section className="page">
      <div className="page-header">
        <h2 className="page-title">My Claims</h2>
        <p className="page-subtitle">Pickups you have claimed.</p>
      </div>

      <div className="row-actions">
        <button
          type="button"
          className="btn-secondary"
          onClick={() => setReloadKey((k) => k + 1)}
        >
          Refresh
        </button>
      </div>

      {loading && <p className="muted">Loading...</p>}
      {error && <p className="auth-error">{error}</p>}

      {!loading && !error && claims.length === 0 && (
        <p className="muted">
          No claims yet. Claim a listing from the "Browse Nearby" tab.
        </p>
      )}

      <div className="grid">
        {claims.map((claim) => (
          <article className="card" key={claim.claim_id}>
            <span className={`badge badge--${claim.status}`}>
              {claim.status}
            </span>
            <h3 className="card-title">Claim {claim.claim_id.slice(0, 8)}</h3>
            <p className="card-meta">Quantity: {claim.claimed_quantity}</p>
            <p className="card-meta">
              Pickup code: <strong>{claim.claim_code}</strong>
            </p>
            <p className="card-meta">
              Code expires: {new Date(claim.code_expires_at).toLocaleString()}
            </p>
            {claim.status === "waiting_for_pickup" && (
              <button
                type="button"
                className="btn-secondary"
                onClick={() => handleCancel(claim.claim_id)}
              >
                Cancel
              </button>
            )}
          </article>
        ))}
      </div>
    </section>
  );
}
