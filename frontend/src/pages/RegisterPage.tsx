import { useState } from "react";
import type { FormEvent } from "react";
import { register } from "../api/auth";
import { ApiError } from "../api/client";
import type { Role } from "../types/auth";

interface RegisterPageProps {
  onRegistered: () => void;
  onSwitchToLogin: () => void;
}

export function RegisterPage({
  onRegistered,
  onSwitchToLogin,
}: RegisterPageProps) {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<Role>("recipient");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setLoading(true);

    try {
      await register({ name, email, phone, role, password });
      onRegistered();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Cannot reach the server. Is it running?");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-screen">
      <main className="auth-card">
        <p className="auth-eyebrow">BagiPangan</p>
        <h1 className="auth-title">Create your account</h1>
        <p className="auth-subtitle">Join as a donor or a recipient.</p>

        <form className="auth-form" onSubmit={handleSubmit}>
          <label className="field">
            <span className="field-label">Name</span>
            <input
              className="field-input"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Budi Santoso"
              autoComplete="name"
              required
            />
          </label>

          <label className="field">
            <span className="field-label">Email</span>
            <input
              className="field-input"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="budi@example.com"
              autoComplete="email"
              required
            />
          </label>

          <label className="field">
            <span className="field-label">Phone</span>
            <input
              className="field-input"
              type="tel"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              placeholder="08123456789"
              autoComplete="tel"
              required
            />
          </label>

          <label className="field">
            <span className="field-label">Password</span>
            <input
              className="field-input"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="********"
              autoComplete="new-password"
              required
            />
          </label>

          <fieldset className="field role-field">
            <legend className="field-label">I want to join as</legend>
            <div className="role-group">
              <label className="role-option">
                <input
                  type="radio"
                  name="role"
                  value="donor"
                  checked={role === "donor"}
                  onChange={() => setRole("donor")}
                />
                <span>Donor</span>
              </label>
              <label className="role-option">
                <input
                  type="radio"
                  name="role"
                  value="recipient"
                  checked={role === "recipient"}
                  onChange={() => setRole("recipient")}
                />
                <span>Recipient</span>
              </label>
            </div>
          </fieldset>

          {error && <p className="auth-error">{error}</p>}

          <button className="btn-primary" type="submit" disabled={loading}>
            {loading ? "Creating account..." : "Create account"}
          </button>
        </form>

        <p className="auth-switch">
          Already have an account?{" "}
          <button type="button" className="auth-link" onClick={onSwitchToLogin}>
            Log in
          </button>
        </p>
      </main>
    </div>
  );
}
