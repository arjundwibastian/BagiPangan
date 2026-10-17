import { useState } from "react";
import type { FormEvent } from "react";
import { createFoodListing } from "../../api/food";
import { ApiError } from "../../api/client";
import type { CreateFoodListingRequest, FoodCategory } from "../../types/food";

interface CreateListingPageProps {
  onCreated: () => void;
}

const FOOD_TYPES: FoodCategory[] = ["cooked_food", "raw_food", "bakery"];

export function CreateListingPage({ onCreated }: CreateListingPageProps) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [latitude, setLatitude] = useState("-6.2");
  const [longitude, setLongitude] = useState("106.816666");
  const [foodType, setFoodType] = useState<FoodCategory>("cooked_food");
  const [quantity, setQuantity] = useState("1");
  const [availableFrom, setAvailableFrom] = useState("");
  const [availableUntil, setAvailableUntil] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setLoading(true);

    const payload: CreateFoodListingRequest = {
      title,
      description: description.trim() ? description : undefined,
      latitude: Number(latitude),
      longitude: Number(longitude),
      food_type: foodType,
      quantity: Number(quantity),
      available_from: new Date(availableFrom).toISOString(),
      available_until: new Date(availableUntil).toISOString(),
    };

    try {
      await createFoodListing(payload);
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
        <h2 className="page-title">Create Listing</h2>
        <p className="page-subtitle">Post food you want to share.</p>
      </div>

      <form className="auth-form" onSubmit={handleSubmit}>
        <div className="form-grid">
          <label className="field field--full">
            <span className="field-label">Title</span>
            <input
              className="field-input"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Nasi Kotak"
              minLength={3}
              required
            />
          </label>

          <label className="field field--full">
            <span className="field-label">Description (optional)</span>
            <textarea
              className="field-input"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Nasi ayam dan sayur"
            />
          </label>

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
            <span className="field-label">Food type</span>
            <select
              className="field-input"
              value={foodType}
              onChange={(e) => setFoodType(e.target.value as FoodCategory)}
            >
              {FOOD_TYPES.map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
          </label>

          <label className="field">
            <span className="field-label">Quantity</span>
            <input
              className="field-input"
              type="number"
              min={0}
              value={quantity}
              onChange={(e) => setQuantity(e.target.value)}
              required
            />
          </label>

          <label className="field">
            <span className="field-label">Available from</span>
            <input
              className="field-input"
              type="datetime-local"
              value={availableFrom}
              onChange={(e) => setAvailableFrom(e.target.value)}
              required
            />
          </label>

          <label className="field">
            <span className="field-label">Available until</span>
            <input
              className="field-input"
              type="datetime-local"
              value={availableUntil}
              onChange={(e) => setAvailableUntil(e.target.value)}
              required
            />
          </label>
        </div>

        {error && <p className="auth-error">{error}</p>}

        <button className="btn-primary" type="submit" disabled={loading}>
          {loading ? "Publishing..." : "Publish listing"}
        </button>
      </form>
    </section>
  );
}