import { useState } from "react";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { AppShell } from "./components/AppShell";
import type { Tab } from "./components/AppShell";
import { MyListingsPage } from "./pages/donor/MyListingsPage";
import { CreateListingPage } from "./pages/donor/CreateListingPage";
import { MyRequestsPage } from "./pages/recipient/MyRequestsPage";
import { CreateRequestPage } from "./pages/recipient/CreateRequestPage";
import { NearbyPage } from "./pages/recipient/NearbyPage";
import type { User } from "./types/auth";

export default function App() {
  const [user, setUser] = useState<User | null>(() => {
    const stored = localStorage.getItem("user");
    return stored ? (JSON.parse(stored) as User) : null;
  });
  const [authView, setAuthView] = useState<"login" | "register">("login");
  const [view, setView] = useState("");

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
      : [
          { id: "requests", label: "My Requests" },
          { id: "create-request", label: "Create Request" },
          { id: "nearby", label: "Browse Nearby" },
        ];

  const activeView = tabs.some((tab) => tab.id === view)
    ? view
    : (tabs[0]?.id ?? "");

  function handleLogout() {
    localStorage.clear();
    setUser(null);
    setAuthView("login");
    setView("");
  }

  return (
    <AppShell
      user={user}
      tabs={tabs}
      current={activeView}
      onNavigate={setView}
      onLogout={handleLogout}
    >
      {activeView === "listings" && <MyListingsPage user={user} />}
      {activeView === "create-listing" && (
        <CreateListingPage onCreated={() => setView("listings")} />
      )}
      {activeView === "requests" && <MyRequestsPage user={user} />}
      {activeView === "create-request" && (
        <CreateRequestPage onCreated={() => setView("requests")} />
      )}
      {activeView === "nearby" && <NearbyPage user={user} />}
    </AppShell>
  );
}
