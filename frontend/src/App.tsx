import { useState } from "react";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import type { User } from "./types/auth";

export default function App() {
  const [user, setUser] = useState<User | null>(() => {
    const stored = localStorage.getItem("user");
    return stored ? (JSON.parse(stored) as User) : null;
  });
  const [authView, setAuthView] = useState<"login" | "register">("login");

  if (!user) {
    return authView === "login" ? (
      <LoginPage
        onLoggedIn={setUser}
        onSwitchToRegister={() => setAuthView("register")}
      />
    ) : (
      <RegisterPage
        onRegistered={setUser}
        onSwitchToLogin={() => setAuthView("login")}
      />
    );
  }

  return (
    <div className="auth-screen">
      <main className="auth-card">
        <h1 className="auth-title">Hi, {user.name}</h1>
        <p className="auth-subtitle">Logged in as {user.role}.</p>
        <button
          className="btn-primary"
          onClick={() => {
            localStorage.clear();
            setUser(null);
            setAuthView("login");
          }}
        >
          Log out
        </button>
      </main>
    </div>
  );
}
