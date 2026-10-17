import { useState } from "react";
import type { FormEvent } from "react";
import { getClaim, verifyPickup } from "../../api/claims";
import { ApiError } from "../../api/client";
import type { FoodClaim } from "../../types/claim";

export function VerifyPickupPage() {
  const [claimId, setClaimId] = useState("");
  const [code, setCode] = useState("");
  const [claim, setClaim] = useState<FoodClaim | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleLookup() {
    setError("");
    setClaim(null);
    if (!claimId) return;
    setLoading(true);
    try {
      setClaim(await getClaim(claimId));
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Cannot reach the server.",
      );
    } finally {
      setLoading(false);
    }
  }

  async function handleVerify(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      const updated = await verifyPickup(claimId, { claim_code: code });
      setClaim(updated);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Cannot reach the server.",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="page">
      <div className="page-header">
        <h2 className="page-title">Verify Pickup</h2>
        <p className="page-subtitle">
          Enter the claim ID and the code the recipient gives you.
        </p>
      </div>

      <div className="form-grid">
        <label className="field field--full">
          <span className="field-label">Claim ID</span>
          <input
            className="field-input"
            value={claimId}
            onChange={(e) => setClaimId(e.target.value)}
            placeholder="00000000-0000-0000-0000-000000000000"
          />
        </label>
      </div>

      <div className="row-actions">
        <button
          type="button"
          className="btn-secondary"
          onClick={handleLookup}
          disabled={loading}
        >
          Look up
        </button>
      </div>

      {claim && (
        <div className="card">
          <span className={`badge badge--${claim.status}`}>{claim.status}</span>
          <h3 className="card-title">Claim {claim.claim_id.slice(0, 8)}</h3>
          <p className="card-meta">Quantity: {claim.claimed_quantity}</p>
          <p className="card-meta">
            Requested listing: {claim.food_listing_id}
          </p>
          <p className="card-meta">
            Code expires: {new Date(claim.code_expires_at).toLocaleString()}
          </p>
        </div>
      )}

      <form className="auth-form" onSubmit={handleVerify}>
        <label className="field">
          <span className="field-label">Pickup code</span>
          <input
            className="field-input"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            placeholder="123456"
            required
          />
        </label>

        {error && <p className="auth-error">{error}</p>}

        <button
          className="btn-primary"
          type="submit"
          disabled={loading || !claimId}
        >
          {loading ? "Verifying..." : "Verify pickup"}
        </button>
      </form>
    </section>
  );
}
