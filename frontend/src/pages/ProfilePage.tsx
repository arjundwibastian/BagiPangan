import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { getMe, updateMe } from "../api/users";
import { ApiError } from "../api/client";
import type { User } from "../types/auth";

interface ProfilePageProps {
  user: User;
  onUpdated: (user: User) => void;
}

export function ProfilePage({ user, onUpdated }: ProfilePageProps) {
  const [name, setName] = useState(user.name);
  const [phone, setPhone] = useState(user.phone);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    let cancelled = false;

    getMe()
      .then((me) => {
        if (cancelled) return;
        setName(me.name);
        setPhone(me.phone);
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
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSaved(false);
    setSaving(true);

    try {
      const updated = await updateMe({ name, phone });
      localStorage.setItem("user", JSON.stringify(updated));
      onUpdated(updated);
      setSaved(true);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Cannot reach the server.",
      );
    } finally {
      setSaving(false);
    }
  }

  return (
    <section className="page">
      <div className="page-header">
        <h2 className="page-title">Profile</h2>
        <p className="page-subtitle">Your account details.</p>
      </div>

      {loading && <p className="muted">Loading...</p>}

      <form className="auth-form" onSubmit={handleSubmit}>
        <div className="form-grid">
          <label className="field field--full">
            <span className="field-label">Name</span>
            <input
              className="field-input"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </label>

          <label className="field field--full">
            <span className="field-label">Phone</span>
            <input
              className="field-input"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              required
            />
          </label>

          <div className="field">
            <span className="field-label">Email</span>
            <span className="field-static">{user.email}</span>
          </div>

          <div className="field">
            <span className="field-label">Role</span>
            <span className="field-static">{user.role}</span>
          </div>
        </div>

        {error && <p className="auth-error">{error}</p>}
        {saved && <p className="success-note">Profile updated.</p>}

        <button className="btn-primary" type="submit" disabled={saving}>
          {saving ? "Saving..." : "Save changes"}
        </button>
      </form>
    </section>
  );
}
