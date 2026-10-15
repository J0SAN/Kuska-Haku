import { useState, useCallback, useEffect } from "react";
import {
  GoogleMap,
  useJsApiLoader,
  Polyline,
  Marker,
} from "@react-google-maps/api";
import apiClient from "../../api/client";

// Coordenadas base (Centro de Cajamarca, Perú)
const CAJAMARCA_CENTER = { lat: -7.162, lng: -78.5 };
const MAP_CONTAINER_STYLE = {
  width: "100%",
  height: "420px",
  borderRadius: "8px",
};

// Interfaces basadas en nuestro Dominio de Go
interface Coordinate {
  latitude: number;
  longitude: number;
}

interface RouteData {
  id: string;
  name: string;
  path: Coordinate[];
}

export default function RutasPage() {
  // 1. Inicializar el SDK de Google Maps
  const { isLoaded } = useJsApiLoader({
    id: "google-map-script",
    // Vite inyecta las variables de entorno con import.meta.env
    googleMapsApiKey: import.meta.env.VITE_GOOGLE_MAPS_API_KEY || "",
  });

  // 2. Estados locales
  const [routes, setRoutes] = useState<RouteData[]>([]);
  const [newPath, setNewPath] = useState<google.maps.LatLngLiteral[]>([]);
  const [routeName, setRouteName] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  // 3. Cargar rutas al montar el componente (Patrón seguro de React)
  useEffect(() => {
    let isMounted = true;

    const loadRoutes = async () => {
      try {
        const response = await apiClient.get("/admin/routes");
        if (isMounted) {
          setRoutes(response.data);
        }
      } catch (error) {
        console.error("Error al cargar rutas desde la base de datos:", error);
      }
    };

    loadRoutes();

    // Función de limpieza: si el componente se desmonta, cancelamos la actualización de estado
    return () => {
      isMounted = false;
    };
  }, []);

  // 4. Capturar clics en el mapa para dibujar
  const handleMapClick = useCallback((e: google.maps.MapMouseEvent) => {
    if (e.latLng) {
      setNewPath((current) => [
        ...current,
        { lat: e.latLng!.lat(), lng: e.latLng!.lng() },
      ]);
    }
  }, []);

  // 5. Guardar la nueva ruta en PostgreSQL/PostGIS
  const handleSaveRoute = async () => {
    if (newPath.length < 2) {
      alert("Debes marcar al menos 2 puntos en el mapa para trazar una ruta.");
      return;
    }
    if (!routeName.trim()) {
      alert("Por favor, ingresa un nombre para la ruta.");
      return;
    }

    setIsSaving(true);
    try {
      const payloadPath = newPath.map((p) => ({
        latitude: p.lat,
        longitude: p.lng,
      }));

      // Guardar en la base de datos
      await apiClient.post("/admin/routes", {
        name: routeName,
        path: payloadPath,
      });

      setNewPath([]);
      setRouteName("");

      // Recargar las rutas directamente desde la base de datos para obtener el UUID generado
      const response = await apiClient.get("/admin/routes");
      setRoutes(response.data);
    } catch (error) {
      console.error(error);
      alert("Error al guardar la ruta. Revisa la consola.");
    } finally {
      setIsSaving(false);
    }
  };

  if (!isLoaded) return <div style={{ padding: "2rem" }}>Cargando Mapa</div>;

  return (
    <div
      style={{
        maxWidth: "1200px",
        margin: "0 auto",
        padding: "20px",
        fontFamily: "sans-serif",
      }}
    >
      <h2 style={{ color: "#FFFFFF" }}>Gestión de Rutas y Trazado</h2>

      {/* Panel de Controles */}
      <div
        style={{
          marginBottom: "20px",
          padding: "15px",
          background: "#f3f4f6",
          borderRadius: "8px",
          border: "1px solid #e5e7eb",
        }}
      >
        <input
          type="text"
          placeholder="Nombre de la ruta (ej. Ruta A)"
          value={routeName}
          onChange={(e) => setRouteName(e.target.value)}
          style={{
            padding: "10px",
            marginRight: "10px",
            width: "250px",
            borderRadius: "4px",
            border: "1px solid #d1d5db",
          }}
        />
        <button
          onClick={handleSaveRoute}
          disabled={isSaving || newPath.length < 2}
          style={{
            padding: "10px 20px",
            background: isSaving || newPath.length < 2 ? "#9ca3af" : "#2563eb",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: isSaving || newPath.length < 2 ? "not-allowed" : "pointer",
            fontWeight: "bold",
          }}
        >
          {isSaving ? "Guardando..." : "Guardar Trazado"}
        </button>

        {newPath.length > 0 && (
          <button
            onClick={() => setNewPath([])}
            style={{
              padding: "10px 20px",
              marginLeft: "10px",
              background: "#ef4444",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
            }}
          >
            Limpiar Puntos
          </button>
        )}

        <p style={{ margin: "10px 0 0 0", fontSize: "14px", color: "#4b5563" }}>
          * Instrucciones: Haz clic en el mapa para marcar los vértices de la
          ruta.
        </p>
      </div>

      {/* Visor de Mapas */}
      <div style={{ border: "2px solid #e5e7eb", borderRadius: "8px" }}>
        <GoogleMap
          mapContainerStyle={MAP_CONTAINER_STYLE}
          center={CAJAMARCA_CENTER}
          zoom={14}
          onClick={handleMapClick}
          options={{
            streetViewControl: false,
            mapTypeControl: false,
          }}
        >
          {/* 1. Dibujar Rutas Históricas (ya guardadas en base de datos) */}
          {routes.map((route) => (
            <Polyline
              key={route.id}
              path={route.path.map((p) => ({
                lat: p.latitude,
                lng: p.longitude,
              }))}
              options={{
                strokeColor: "#9ca3af",
                strokeOpacity: 0.8,
                strokeWeight: 4,
              }}
            />
          ))}

          {/* 2. Dibujar la ruta que el usuario está trazando en este momento */}
          {newPath.length > 0 && (
            <Polyline
              path={newPath}
              options={{
                strokeColor: "#2563eb",
                strokeOpacity: 1,
                strokeWeight: 5,
              }}
            />
          )}

          {/* Puntos (Marcadores) del nuevo trazado */}
          {newPath.map((point, index) => (
            <Marker
              key={index}
              position={point}
              icon={{
                path: window.google.maps.SymbolPath.CIRCLE,
                scale: 5,
                fillColor: "#2563eb",
                fillOpacity: 1,
                strokeWeight: 2,
                strokeColor: "#ffffff",
              }}
            />
          ))}
        </GoogleMap>
      </div>
    </div>
  );
}
