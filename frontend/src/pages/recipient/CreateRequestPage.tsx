import { useState } from "react";
import type { FormEvent } from "react";
import { createFoodRequest } from "../../api/requests";
import { ApiError } from "../../api/client";
import type { CreateFoodRequestRequest } from "../../types/request";

interface CreateRequestPageProps {
  onCreated: () => void;
}

export function CreateRequestPage({ onCreated }: CreateRequestPageProps) {
  const [latitude, setLatitude] = useState("-6.2");
  const [longitude, setLongitude] = useState("106.816666");
  const [radiusKm, setRadiusKm] = useState("5");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setLoading(true);

    const payload: CreateFoodRequestRequest = {
      latitude: Number(latitude),
      longitude: Number(longitude),
      radius_km: Number(radiusKm),
    };

    try {
      await createFoodRequest(payload);
      onCreated();
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
        <h2 className="page-title">Create Request</h2>
        <p className="page-subtitle">Tell us where and how far to look.</p>
      </div>

      <form className="auth-form" onSubmit={handleSubmit}>
        <div className="form-grid">
          <label className="field">
            <span className="field-label">Latitude</span>
            <input
              className="field-input"
              type="number"
              step="any"
              value={latitude}
              onChange={(e) => setLatitude(e.target.value)}
              required
            />
          </label>

          <label className="field">
            <span className="field-label">Longitude</span>
            <input
              className="field-input"
              type="number"
              step="any"
              value={longitude}
              onChange={(e) => setLongitude(e.target.value)}
              required
            />
          </label>

          <label className="field">
            <span className="field-label">Radius (km)</span>
            <input
              className="field-input"
              type="number"
              step="any"
              min={1}
              value={radiusKm}
              onChange={(e) => setRadiusKm(e.target.value)}
              required
            />
          </label>
        </div>

        {error && <p className="auth-error">{error}</p>}

        <button className="btn-primary" type="submit" disabled={loading}>
          {loading ? "Creating..." : "Create request"}
        </button>
      </form>
    </section>
  );
}
