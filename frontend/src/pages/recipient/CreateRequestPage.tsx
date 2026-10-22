import { useState } from "react";
import type { FormEvent } from "react";
import { createFoodRequest } from "../../api/requests";
import { ApiError } from "../../api/client";
import { MapPicker } from "../../components/MapPicker";
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
          <div className="field field--full">
            <span className="field-label">Search location</span>
            <MapPicker
              latitude={Number(latitude)}
              longitude={Number(longitude)}
              radiusKm={Number(radiusKm)}
              onChange={(lat, lng) => {
                setLatitude(lat.toFixed(6));
                setLongitude(lng.toFixed(6));
              }}
            />
          </div>

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
