import { useState, useEffect } from "react";
import httpClient, { getErrorMessage } from "@/shared/api/httpClient";
import type { FleetData } from "./fleet.types";
import type { RouteData } from "@/features/transport_routes";
import type { UserData } from "@/features/auth";
import { ConfirmDialog } from "@/shared/ui/ConfirmDialog";

export default function FleetPage() {
  const [routes, setRoutes] = useState<RouteData[]>([]);
  const [fleet, setFleet] = useState<FleetData[]>([]);
  const [drivers, setDrivers] = useState<UserData[]>([]);
  const [busId, setBusId] = useState("");
  const [routeId, setRouteId] = useState("");
  const [driverId, setDriverId] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [busToDelete, setBusToDelete] = useState<FleetData | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadInitialData = async () => {
      try {
        const [routesRes, fleetRes, driversRes] = await Promise.all([
          httpClient.get<RouteData[]>("/admin/routes"),
          httpClient.get<FleetData[]>("/admin/fleet"),
          httpClient.get<UserData[]>("/admin/users?role=driver"),
        ]);
        if (isMounted) {
          const loadedRoutes = routesRes.data || [];
          setRoutes(loadedRoutes);
          setFleet(fleetRes.data || []);
          setDrivers(driversRes.data || []);
          if (loadedRoutes.length > 0) {
            setRouteId((current) => current || loadedRoutes[0].id);
          }
        }
      } catch (err) {
        if (isMounted) {
          setError(getErrorMessage(err, "Error cargando la flota"));
        }
      }
    };

    void loadInitialData();

    return () => {
      isMounted = false;
    };
  }, []);

  const refreshFleet = async () => {
    const fleetRes = await httpClient.get<FleetData[]>("/admin/fleet");
    setFleet(fleetRes.data || []);
  };

  const handleSaveBus = async () => {
    setError(null);
    if (!busId.trim()) {
      setError("El código del bus (ej. bus-01) es obligatorio.");
      return;
    }
    if (!routeId) {
      setError(
        "Debes crear y seleccionar al menos una ruta antes de registrar un bus.",
      );
      return;
    }

    setIsSaving(true);
    try {
      const trimmedDriver = driverId.trim();
      await httpClient.post("/admin/fleet", {
        busId: busId.trim(),
        routeId,
        driverId: trimmedDriver ? trimmedDriver : null,
        status: "active",
      });

      setBusId("");
      setDriverId("");
      await refreshFleet();
    } catch (err) {
      setError(getErrorMessage(err, "Error guardando el bus"));
    } finally {
      setIsSaving(false);
    }
  };

  const handleAssignDriver = async (bus: FleetData, newDriverId: string) => {
    setError(null);
    try {
      await httpClient.put(`/admin/fleet/${bus.id}`, {
        busId: bus.busId,
        routeId: bus.routeId,
        driverId: newDriverId ? newDriverId : null,
        status: bus.status || "active",
      });
    } catch (err) {
      setError(getErrorMessage(err, "Error asignando el conductor"));
    }
    // Se recarga siempre: si falló por 409 Conflict, el selector vuelve al valor real en BD.
    await refreshFleet();
  };

  const handleConfirmDelete = async () => {
    if (!busToDelete) return;
    setError(null);
    try {
      await httpClient.delete(`/admin/fleet/${busToDelete.id}`);
      await refreshFleet();
    } catch (err) {
      setError(getErrorMessage(err, "Error eliminando el bus"));
    } finally {
      setBusToDelete(null);
    }
  };

  const getRouteName = (id: string) =>
    routes.find((r) => r.id === id)?.name ?? "—";

  const getDriverLabel = (d: UserData) => d.fullName || d.full_name || d.email;

  const renderDriverSelector = (
    value: string,
    onChange: (id: string) => void,
  ) => (
    <select value={value} onChange={(e) => onChange(e.target.value)}>
      <option value="">Sin conductor</option>
      {drivers.map((d) => (
        <option key={d.id} value={d.id}>
          {getDriverLabel(d)}
        </option>
      ))}
    </select>
  );

  return (
    <section>
      <h2>Gestión de flota</h2>

      <div className="form-row">
        <input
          placeholder="Código (ej. bus-01)"
          value={busId}
          onChange={(e) => setBusId(e.target.value)}
        />
        <select value={routeId} onChange={(e) => setRouteId(e.target.value)}>
          {routes.length === 0 && (
            <option value="">Sin rutas disponibles</option>
          )}
          {routes.map((r) => (
            <option key={r.id} value={r.id}>
              {r.name}
            </option>
          ))}
        </select>
        {renderDriverSelector(driverId, setDriverId)}
        <button
          type="button"
          className="primary"
          onClick={handleSaveBus}
          disabled={isSaving}
        >
          {isSaving ? "Registrando..." : "Registrar bus"}
        </button>
      </div>

      {error && <p className="error">{error}</p>}

      <table className="data-table">
        <thead>
          <tr>
            <th>Código (Bus ID)</th>
            <th>Ruta</th>
            <th>Conductor</th>
            <th>Estado</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {fleet.length === 0 && (
            <tr>
              <td colSpan={5} className="empty">
                Todavía no hay buses registrados.
              </td>
            </tr>
          )}
          {fleet.map((b) => (
            <tr key={b.id}>
              <td>{b.busId}</td>
              <td>{getRouteName(b.routeId)}</td>
              <td>
                {renderDriverSelector(
                  b.driverId ?? "",
                  (id) => void handleAssignDriver(b, id),
                )}
              </td>
              <td>
                <span
                  className={
                    b.status === "active" ? "badge badge-new" : "badge"
                  }
                >
                  {b.status === "active" ? "Activo" : b.status}
                </span>
              </td>
              <td>
                <button
                  type="button"
                  className="danger-link"
                  onClick={() => setBusToDelete(b)}
                >
                  Eliminar
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {busToDelete && (
        <ConfirmDialog
          title="Eliminar bus"
          message={`¿Seguro que deseas eliminar el bus "${busToDelete.busId}"? Esta acción no se puede deshacer.`}
          onConfirm={handleConfirmDelete}
          onCancel={() => setBusToDelete(null)}
        />
      )}
    </section>
  );
}
