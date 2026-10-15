import { Routes, Route, Navigate, Outlet } from "react-router-dom";
import { useAuthStore } from "./store/authStore";

import LoginPage from "./features/auth/LoginPage";
import TransportRoutesPage from "./features/transport_routes/TransportRoutesPage";
import FlotaPage from "./features/fleet/FleetPage";

// Guardián de Rutas Protegidas (RBAC)
function RequireAuth() {
  const { isAuthenticated, isAdmin } = useAuthStore();

  // Si no hay token o el rol no es admin, redirigir al login
  if (!isAuthenticated || !isAdmin) {
    return <Navigate to="/login" replace />;
  }

  // Si es válido, renderizar las vistas hijas (Rutas, Flota, etc.)
  return <Outlet />;
}

export default function App() {
  return (
    <Routes>
      {/* Ruta Pública */}
      <Route path="/login" element={<LoginPage />} />

      {/* Rutas Protegidas */}
      <Route element={<RequireAuth />}>
        <Route path="/rutas" element={<TransportRoutesPage />} />
        <Route path="/flota" element={<FlotaPage />} />
      </Route>

      {/* Redirección por defecto */}
      <Route path="*" element={<Navigate to="/rutas" replace />} />
    </Routes>
  );
}
