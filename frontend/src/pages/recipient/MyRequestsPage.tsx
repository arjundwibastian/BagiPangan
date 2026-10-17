import { useEffect, useState } from "react";
import { cancelFoodRequest, listFoodRequests } from "../../api/requests";
import { ApiError } from "../../api/client";
import type { FoodRequest } from "../../types/request";
import type { User } from "../../types/auth";

interface MyRequestsPageProps {
  user: User;
}

export function MyRequestsPage({ user }: MyRequestsPageProps) {
  const [requests, setRequests] = useState<FoodRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError("");

    listFoodRequests(user.user_id)
      .then((data) => {
        if (!cancelled) setRequests(data ?? []);
      })
      .catch((err) => {
        if (!cancelled) {
          setError(
            err instanceof ApiError ? err.message : "Cannot reach the server.",
          );
        }
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
      await cancelFoodRequest(id);
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
        <h2 className="page-title">My Requests</h2>
        <p className="page-subtitle">What you are looking for nearby.</p>
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

      {!loading && !error && requests.length === 0 && (
        <p className="muted">
          No requests yet. Create one from the "Create Request" tab.
        </p>
      )}

      <div className="grid">
        {requests.map((request) => (
          <article className="card" key={request.request_id}>
            <span className={`badge badge--${request.status}`}>
              {request.status}
            </span>
            <h3 className="card-title">
              Request {request.request_id.slice(0, 8)}
            </h3>
            <p className="card-meta">Radius: {request.radius_km} km</p>
            <p className="card-meta">
              At {request.latitude}, {request.longitude}
            </p>
            <p className="card-meta">
              Expires: {new Date(request.expires_at).toLocaleString()}
            </p>
            {request.status === "searching" && (
              <button
                type="button"
                className="btn-secondary"
                onClick={() => handleCancel(request.request_id)}
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
