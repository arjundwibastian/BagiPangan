import { useState } from "react";
import type { FormEvent } from "react";
import { login } from "../api/auth";
import { ApiError } from "../api/client";
import type { User } from "../types/auth";

interface LoginPageProps {
  onLoggedIn: (user: User) => void;
  onSwitchToRegister: () => void;
  notice?: string;
}

export function LoginPage({
  onLoggedIn,
  onSwitchToRegister,
  notice,
}: LoginPageProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setLoading(true);

    try {
      const result = await login({ email, password });
      localStorage.setItem("access_token", result.access_token);
      localStorage.setItem("refresh_token", result.refresh_token);
      localStorage.setItem("user", JSON.stringify(result.user));
      onLoggedIn(result.user);
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
        <h1 className="auth-title">Welcome back</h1>
        <p className="auth-subtitle">Share extra food with people nearby.</p>

        {notice && <p className="success-note">{notice}</p>}

        <form className="auth-form" onSubmit={handleSubmit}>
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
            <span className="field-label">Password</span>
            <input
              className="field-input"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="********"
              autoComplete="current-password"
              required
            />
          </label>

          {error && <p className="auth-error">{error}</p>}

          <button className="btn-primary" type="submit" disabled={loading}>
            {loading ? "Signing in..." : "Log in"}
          </button>
        </form>

        <p className="auth-switch">
          New here?{" "}
          <button
            type="button"
            className="auth-link"
            onClick={onSwitchToRegister}
          >
            Create an account
          </button>
        </p>
      </main>
    </div>
  );
}
