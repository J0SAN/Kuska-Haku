import { useState, useEffect } from "react";
import apiClient from "../../api/client";

interface RouteData {
  id: string;
  name: string;
}

interface FleetData {
  id: string;
  busId: string;
  routeId: string;
  status: string;
}

export default function FlotaPage() {
  const [routes, setRoutes] = useState<RouteData[]>([]);
  const [fleet, setFleet] = useState<FleetData[]>([]);

  // Estados del formulario
  const [busId, setBusId] = useState("");
  const [selectedRouteId, setSelectedRouteId] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  // Cargar datos iniciales
  useEffect(() => {
    let isMounted = true;

    const loadData = async () => {
      try {
        // Ejecutamos ambas peticiones en paralelo para mayor velocidad
        const [routesRes, fleetRes] = await Promise.all([
          apiClient.get("/admin/routes"),
          apiClient.get("/admin/fleet"),
        ]);

        if (isMounted) {
          setRoutes(routesRes.data || []);
          setFleet(fleetRes.data || []);
        }
      } catch (error) {
        console.error("Error cargando datos:", error);
      }
    };

    loadData();
    return () => {
      isMounted = false;
    };
  }, []);

  const handleRegisterBus = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!busId.trim() || !selectedRouteId) {
      alert("Debes ingresar el ID del bus y seleccionar una ruta.");
      return;
    }

    setIsSaving(true);
    try {
      await apiClient.post("/admin/fleet", {
        busId: busId,
        routeId: selectedRouteId,
        status: "active",
      });

      setBusId("");
      setSelectedRouteId("");

      // Recargar la tabla de flota
      const response = await apiClient.get("/admin/fleet");
      setFleet(response.data || []);

      alert("Vehículo registrado exitosamente");
    } catch (error: any) {
      console.error(error);
      alert(error.response?.data?.message || "Error al registrar el vehículo.");
    } finally {
      setIsSaving(false);
    }
  };

  // Función auxiliar para mostrar el nombre de la ruta en la tabla en lugar del UUID
  const getRouteName = (routeId: string) => {
    const route = routes.find((r) => r.id === routeId);
    return route ? route.name : "Ruta Desconocida";
  };

  return (
    <div
      style={{
        maxWidth: "1000px",
        margin: "0 auto",
        padding: "20px",
        fontFamily: "sans-serif",
      }}
    >
      <h2 style={{ color: "#FFFFFF", marginBottom: "20px" }}>
        Gestión de Flota Operativa
      </h2>

      {/* Formulario de Registro */}
      <div
        style={{
          padding: "20px",
          background: "#f9fafb",
          borderRadius: "8px",
          border: "1px solid #e5e7eb",
          marginBottom: "30px",
        }}
      >
        <h3 style={{ marginTop: 0, fontSize: "1.1rem" }}>
          Vincular Nuevo Vehículo
        </h3>
        <form
          onSubmit={handleRegisterBus}
          style={{ display: "flex", gap: "15px", alignItems: "flex-end" }}
        >
          <div style={{ flex: 1 }}>
            <label
              style={{
                display: "block",
                marginBottom: "5px",
                fontSize: "14px",
                fontWeight: "bold",
              }}
            >
              Identificador del Bus (MQTT ID)
            </label>
            <input
              type="text"
              placeholder="ej. bus-01"
              value={busId}
              onChange={(e) => setBusId(e.target.value)}
              style={{
                width: "100%",
                padding: "10px",
                borderRadius: "4px",
                border: "1px solid #d1d5db",
                boxSizing: "border-box",
              }}
              required
            />
          </div>

          <div style={{ flex: 1 }}>
            <label
              style={{
                display: "block",
                marginBottom: "5px",
                fontSize: "14px",
                fontWeight: "bold",
              }}
            >
              Ruta Asignada
            </label>
            <select
              value={selectedRouteId}
              onChange={(e) => setSelectedRouteId(e.target.value)}
              style={{
                width: "100%",
                padding: "10px",
                borderRadius: "4px",
                border: "1px solid #d1d5db",
                boxSizing: "border-box",
              }}
              required
            >
              <option value="">-- Selecciona una ruta --</option>
              {routes.map((route) => (
                <option key={route.id} value={route.id}>
                  {route.name}
                </option>
              ))}
            </select>
          </div>

          <button
            type="submit"
            disabled={isSaving}
            style={{
              padding: "10px 24px",
              background: isSaving ? "#9ca3af" : "#10b981",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: isSaving ? "not-allowed" : "pointer",
              fontWeight: "bold",
              height: "40px",
            }}
          >
            {isSaving ? "Guardando..." : "Registrar"}
          </button>
        </form>
      </div>

      {/* Tabla de Vehículos Activos */}
      <h3 style={{ fontSize: "1.1rem" }}>Vehículos Registrados</h3>
      <table
        style={{ width: "100%", borderCollapse: "collapse", textAlign: "left" }}
      >
        <thead>
          <tr
            style={{ background: "#f3f4f6", borderBottom: "2px solid #e5e7eb" }}
          >
            <th style={{ padding: "12px" }}>ID del Bus</th>
            <th style={{ padding: "12px" }}>Ruta Asignada</th>
            <th style={{ padding: "12px" }}>Estado</th>
          </tr>
        </thead>
        <tbody>
          {fleet.length === 0 ? (
            <tr>
              <td
                colSpan={3}
                style={{
                  padding: "20px",
                  textAlign: "center",
                  color: "#6b7280",
                }}
              >
                No hay vehículos registrados en el sistema.
              </td>
            </tr>
          ) : (
            fleet.map((vehicle) => (
              <tr
                key={vehicle.id}
                style={{ borderBottom: "1px solid #e5e7eb" }}
              >
                <td style={{ padding: "12px", fontWeight: "bold" }}>
                  {vehicle.busId}
                </td>
                <td style={{ padding: "12px" }}>
                  {getRouteName(vehicle.routeId)}
                </td>
                <td style={{ padding: "12px" }}>
                  <span
                    style={{
                      padding: "4px 8px",
                      borderRadius: "9999px",
                      fontSize: "12px",
                      background:
                        vehicle.status === "active" ? "#d1fae5" : "#fee2e2",
                      color:
                        vehicle.status === "active" ? "#065f46" : "#991b1b",
                    }}
                  >
                    {vehicle.status}
                  </span>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
