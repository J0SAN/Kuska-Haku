import axios from "axios";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import apiClient from "../../api/client";
import { useAuthStore } from "../../store/authStore";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();
  // Extraemos la función login de nuestro estado global
  const login = useAuthStore((state) => state.login);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      // Hacemos el POST al endpoint exacto que definimos en Go (Fase 2)
      const response = await apiClient.post("/auth/login", {
        email,
        password,
      });

      // Go nos devuelve el token JWT
      const { token } = response.data;

      // Zustand guarda el token en localStorage y valida el rol
      login(token);

      // Si todo sale bien, lo empujamos a la ruta protegida
      navigate("/rutas", { replace: true });
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        // Si es un error HTTP devuelto por Go (ej. 401)
        setError(err.response?.data?.message || 'Credenciales inválidas o error de servidor');
      } else {
        // Si es un error interno de React o de red (se cayó el servidor)
        setError('Ocurrió un error inesperado');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      style={{
        maxWidth: "400px",
        margin: "4rem auto",
        fontFamily: "sans-serif",
      }}
    >
      <h2>Kuska Haku</h2>

      {error && (
        <div
          style={{
            background: "#fee2e2",
            color: "#991b1b",
            padding: "10px",
            marginBottom: "1rem",
            borderRadius: "4px",
          }}
        >
          {error}
        </div>
      )}

      <form
        onSubmit={handleSubmit}
        style={{ display: "flex", flexDirection: "column", gap: "1rem" }}
      >
        <div>
          <label style={{ display: "block", marginBottom: "0.5rem" }}>
            Correo Electrónico
          </label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            style={{ width: "100%", padding: "8px", boxSizing: "border-box" }}
          />
        </div>

        <div>
          <label style={{ display: "block", marginBottom: "0.5rem" }}>
            Contraseña
          </label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            style={{ width: "100%", padding: "8px", boxSizing: "border-box" }}
          />
        </div>

        <button
          type="submit"
          disabled={isLoading}
          style={{
            padding: "10px",
            background: "#2563eb",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: isLoading ? "not-allowed" : "pointer",
          }}
        >
          {isLoading ? "Ingresando..." : "Iniciar Sesión"}
        </button>
      </form>
    </div>
  );
}
