import { useState } from "react";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { AppShell } from "./components/AppShell";
import type { Tab } from "./components/AppShell";
import { ProfilePage } from "./pages/ProfilePage";
import { MyListingsPage } from "./pages/donor/MyListingsPage";
import { CreateListingPage } from "./pages/donor/CreateListingPage";
import { VerifyPickupPage } from "./pages/donor/VerifyPickupPage";
import { MyRequestsPage } from "./pages/recipient/MyRequestsPage";
import { CreateRequestPage } from "./pages/recipient/CreateRequestPage";
import { NearbyPage } from "./pages/recipient/NearbyPage";
import { MyClaimsPage } from "./pages/recipient/MyClaimsPage";
import { logout } from "./api/auth";
import type { User } from "./types/auth";

export default function App() {
  const [user, setUser] = useState<User | null>(() => {
    const stored = localStorage.getItem("user");
    return stored ? (JSON.parse(stored) as User) : null;
  });
  const [authView, setAuthView] = useState<"login" | "register">("login");
  const [notice, setNotice] = useState("");
  const [view, setView] = useState("");

  if (!user) {
    return authView === "login" ? (
      <LoginPage
        notice={notice}
        onLoggedIn={setUser}
        onSwitchToRegister={() => {
          setNotice("");
          setAuthView("register");
        }}
      />
    ) : (
      <RegisterPage
        onRegistered={() => {
          setNotice("Account created successfully. Please sign in.");
          setAuthView("login");
        }}
        onSwitchToLogin={() => {
          setNotice("");
          setAuthView("login");
        }}
      />
    );
  }

  const tabs: Tab[] =
    user.role === "donor"
      ? [
          { id: "listings", label: "My Listings" },
          { id: "create-listing", label: "Create Listing" },
          { id: "verify", label: "Verify Pickup" },
          { id: "profile", label: "Profile" },
        ]
      : [
          { id: "requests", label: "My Requests" },
          { id: "create-request", label: "Create Request" },
          { id: "nearby", label: "Browse Nearby" },
          { id: "claims", label: "My Claims" },
          { id: "profile", label: "Profile" },
        ];

  const activeView = tabs.some((tab) => tab.id === view)
    ? view
    : (tabs[0]?.id ?? "");

  function handleLogout() {
    const refresh = localStorage.getItem("refresh_token");
    if (refresh) {
      void logout(refresh).catch(() => undefined);
    }
    localStorage.clear();
    setUser(null);
    setAuthView("login");
    setNotice("");
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
      {activeView === "verify" && <VerifyPickupPage user={user} />}
      {activeView === "requests" && <MyRequestsPage user={user} />}
      {activeView === "create-request" && (
        <CreateRequestPage onCreated={() => setView("requests")} />
      )}
      {activeView === "nearby" && <NearbyPage user={user} />}
      {activeView === "claims" && <MyClaimsPage user={user} />}
      {activeView === "profile" && (
        <ProfilePage user={user} onUpdated={setUser} />
      )}
    </AppShell>
  );
}
