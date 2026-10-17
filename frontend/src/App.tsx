import { useState } from "react";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { AppShell } from "./components/AppShell";
import type { Tab } from "./components/AppShell";
import { MyListingsPage } from "./pages/donor/MyListingsPage";
import { CreateListingPage } from "./pages/donor/CreateListingPage";
import type { User } from "./types/auth";

export default function App() {
  const [user, setUser] = useState<User | null>(() => {
    const stored = localStorage.getItem("user");
    return stored ? (JSON.parse(stored) as User) : null;
  });
  const [authView, setAuthView] = useState<"login" | "register">("login");
  const [view, setView] = useState("listings");

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

  const tabs: Tab[] =
    user.role === "donor"
      ? [
          { id: "listings", label: "My Listings" },
          { id: "create-listing", label: "Create Listing" },
        ]
      : [];

  function handleLogout() {
    localStorage.clear();
    setUser(null);
    setAuthView("login");
    setView("listings");
  }

  return (
    <AppShell
      user={user}
      tabs={tabs}
      current={view}
      onNavigate={setView}
      onLogout={handleLogout}
    >
      {user.role === "donor" && view === "listings" && (
        <MyListingsPage user={user} />
      )}
      {user.role === "donor" && view === "create-listing" && (
        <CreateListingPage onCreated={() => setView("listings")} />
      )}
      {tabs.length === 0 && (
        <section className="page">
          <div className="page-header">
            <h2 className="page-title">Welcome, {user.name}</h2>
            <p className="page-subtitle">
              Recipient features arrive in the next batch.
            </p>
          </div>
        </section>
      )}
    </AppShell>
  );
}
