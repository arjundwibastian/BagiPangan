import type { ReactNode } from "react";
import type { User } from "../types/auth";

export interface Tab {
  id: string;
  label: string;
}

interface AppShellProps {
  user: User;
  tabs: Tab[];
  current: string;
  onNavigate: (id: string) => void;
  onLogout: () => void;
  children: ReactNode;
}

export function AppShell({
  user,
  tabs,
  current,
  onNavigate,
  onLogout,
  children,
}: AppShellProps) {
  return (
    <div className="shell">
      <header className="shell-header">
        <h1 className="shell-brand">BagiPangan</h1>
        <span className="shell-spacer" />
        <span className="shell-user">
          {user.name} · {user.role}
        </span>
        <button type="button" className="btn-secondary" onClick={onLogout}>
          Log out
        </button>
      </header>

      <nav className="shell-nav">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            type="button"
            className={`nav-btn${current === tab.id ? " nav-btn--active" : ""}`}
            onClick={() => onNavigate(tab.id)}
          >
            {tab.label}
          </button>
        ))}
      </nav>

      <main className="shell-main">{children}</main>
    </div>
  );
}
