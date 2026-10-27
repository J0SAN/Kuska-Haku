import { lazy, Suspense } from "react";
import {
  Routes,
  Route,
  Navigate,
  Outlet,
  useNavigate,
  useLocation,
} from "react-router-dom";
import { useAuthStore } from "./store/authStore";
import LoginPage from "./features/auth/LoginPage";

const DashboardPage = lazy(() => import("./features/dashboard/DashboardPage"));
const TransportRoutesPage = lazy(
  () => import("./features/transport_routes/TransportRoutesPage"),
);
const FleetPage = lazy(() => import("./features/fleet/FleetPage"));
const SuggestionsPage = lazy(
  () => import("./features/suggestions/SuggestionsPage"),
);

const NAV_ITEMS: [string, string][] = [
  ["/dashboard", "Dashboard"],
  ["/routes", "Rutas"],
  ["/fleet", "Flota"],
  ["/suggestions", "Sugerencias"],
];

function AdminLayout() {
  const { isAuthenticated, isAdmin, user, logout } = useAuthStore();
  const navigate = useNavigate();
  const location = useLocation();

  if (!isAuthenticated || !isAdmin) {
    return <Navigate to="/login" replace />;
  }

  const handleLogout = () => {
    logout();
    navigate("/login", { replace: true });
  };

  return (
    <div className="app">
      <header className="topbar">
        <span className="brand-small">Kuska Haku</span>
        <nav>
          {NAV_ITEMS.map(([path, label]) => (
            <button
              key={path}
              type="button"
              className={location.pathname === path ? "tab active" : "tab"}
              onClick={() => navigate(path)}
            >
              {label}
            </button>
          ))}
        </nav>
        <div className="user-area">
          <span>{user?.email ?? "Admin"}</span>
          <button type="button" className="secondary" onClick={handleLogout}>
            Salir
          </button>
        </div>
      </header>

      <main className="content">
        <Suspense fallback={<p className="hint">Cargando...</p>}>
          <Outlet />
        </Suspense>
      </main>
    </div>
  );
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />

      <Route element={<AdminLayout />}>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/routes" element={<TransportRoutesPage />} />
        <Route path="/fleet" element={<FleetPage />} />
        <Route path="/suggestions" element={<SuggestionsPage />} />
        <Route path="/rutas" element={<Navigate to="/routes" replace />} />
        <Route path="/flota" element={<Navigate to="/fleet" replace />} />
      </Route>

      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}
