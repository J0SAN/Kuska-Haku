import { useState, useCallback, useEffect } from "react";
import {
  GoogleMap,
  useJsApiLoader,
  Polyline,
  Marker,
} from "@react-google-maps/api";
import apiClient, { getErrorMessage, type RouteData } from "../../api/client";
import { ConfirmDialog } from "../ConfirmDialog";

const CAJAMARCA_CENTER = { lat: -7.1611, lng: -78.5127 };

export default function TransportRoutesPage() {
  const { isLoaded } = useJsApiLoader({
    id: "google-map-script",
    googleMapsApiKey: import.meta.env.VITE_GOOGLE_MAPS_API_KEY || "",
  });

  const [routes, setRoutes] = useState<RouteData[]>([]);
  const [newPath, setNewPath] = useState<google.maps.LatLngLiteral[]>([]);
  const [routeName, setRouteName] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [routeToDelete, setRouteToDelete] = useState<RouteData | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadRoutes = async () => {
      try {
        const response = await apiClient.get<RouteData[]>("/admin/routes");
        if (isMounted) {
          setRoutes(response.data || []);
        }
      } catch (err) {
        if (isMounted) {
          setError(getErrorMessage(err, "Error cargando rutas"));
        }
      }
    };

    void loadRoutes();

    return () => {
      isMounted = false;
    };
  }, []);

  const refreshRoutes = async () => {
    const response = await apiClient.get<RouteData[]>("/admin/routes");
    setRoutes(response.data || []);
  };

  const handleMapClick = useCallback((e: google.maps.MapMouseEvent) => {
    if (e.latLng) {
      setNewPath((current) => [
        ...current,
        { lat: e.latLng!.lat(), lng: e.latLng!.lng() },
      ]);
    }
  }, []);

  const handleSaveRoute = async () => {
    setError(null);
    if (!routeName.trim()) {
      setError("El nombre de la ruta es obligatorio.");
      return;
    }
    if (newPath.length < 2) {
      setError(
        "El trazado debe tener al menos 2 puntos. Haz clic en el mapa para dibujarlo.",
      );
      return;
    }

    setIsSaving(true);
    try {
      const payloadPath = newPath.map((p) => ({
        latitude: p.lat,
        longitude: p.lng,
      }));

      await apiClient.post("/admin/routes", {
        name: routeName.trim(),
        path: payloadPath,
      });

      setRouteName("");
      setNewPath([]);
      await refreshRoutes();
    } catch (err) {
      setError(getErrorMessage(err, "Error guardando la ruta"));
    } finally {
      setIsSaving(false);
    }
  };

  const handleConfirmDelete = async () => {
    if (!routeToDelete) return;
    setError(null);
    try {
      await apiClient.delete(`/admin/routes/${routeToDelete.id}`);
      await refreshRoutes();
    } catch (err) {
      setError(getErrorMessage(err, "Error eliminando la ruta"));
    } finally {
      setRouteToDelete(null);
    }
  };

  return (
    <section>
      <h2>Trazado de rutas</h2>
      <p className="hint">
        Haz clic sobre el mapa para ir marcando el recorrido. Se necesitan al
        menos 2 puntos.
      </p>

      <div className="form-row">
        <input
          type="text"
          placeholder="Nombre o código de la ruta (ej. R-01 Centro)"
          value={routeName}
          onChange={(e) => setRouteName(e.target.value)}
        />
        <button
          type="button"
          className="primary"
          onClick={handleSaveRoute}
          disabled={isSaving}
        >
          {isSaving ? "Guardando..." : "Guardar ruta"}
        </button>
        <button
          type="button"
          className="secondary"
          onClick={() => setNewPath([])}
        >
          Limpiar trazado
        </button>
      </div>

      <p className="hint">Puntos marcados: {newPath.length}</p>
      {error && <p className="error">{error}</p>}

      <div className="map-container">
        {isLoaded ? (
          <GoogleMap
            mapContainerStyle={{ width: "100%", height: "420px" }}
            center={CAJAMARCA_CENTER}
            zoom={14}
            onClick={handleMapClick}
          >
            {newPath.map((point, index) => (
              <Marker key={index} position={point} label={String(index + 1)} />
            ))}
            {newPath.length >= 2 && (
              <Polyline
                path={newPath}
                options={{ strokeColor: "#2E7D32", strokeWeight: 4 }}
              />
            )}
            {routes.map((route) => (
              <Polyline
                key={route.id}
                path={(route.path || []).map((p) => ({
                  lat: p.latitude,
                  lng: p.longitude,
                }))}
                options={{
                  strokeColor: "#81C784",
                  strokeWeight: 3,
                  strokeOpacity: 0.7,
                }}
              />
            ))}
          </GoogleMap>
        ) : (
          <p className="hint">
            Cargando mapa... (revisa VITE_GOOGLE_MAPS_API_KEY en .env.local)
          </p>
        )}
      </div>

      <h3>Rutas registradas</h3>
      <table className="data-table">
        <thead>
          <tr>
            <th>Nombre / Código</th>
            <th>Puntos</th>
            <th>Estado</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {routes.length === 0 && (
            <tr>
              <td colSpan={4} className="empty">
                Todavía no hay rutas registradas.
              </td>
            </tr>
          )}
          {routes.map((r) => (
            <tr key={r.id}>
              <td>{r.name}</td>
              <td>{r.path?.length ?? 0}</td>
              <td>Activa</td>
              <td>
                <button
                  type="button"
                  className="danger-link"
                  onClick={() => setRouteToDelete(r)}
                >
                  Eliminar
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {routeToDelete && (
        <ConfirmDialog
          title="Eliminar ruta"
          message={`¿Seguro que deseas eliminar la ruta "${routeToDelete.name}"? Esta acción no se puede deshacer.`}
          onConfirm={handleConfirmDelete}
          onCancel={() => setRouteToDelete(null)}
        />
      )}
    </section>
  );
}
