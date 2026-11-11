import { useState } from "react";
import { useNavigate } from "react-router-dom";
import httpClient, { getErrorMessage } from "@/shared/api/httpClient";
import { useAuthStore } from "@/features/auth/auth.store";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      const response = await httpClient.post("/auth/login", {
        email: email.trim(),
        password,
      });

      const { token, user } = response.data;

      if (user && user.role !== "admin") {
        setError("Esta cuenta no tiene permisos de administrador.");
        return;
      }

      login(token);
      navigate("/dashboard", { replace: true });
    } catch (err: unknown) {
      setError(
        getErrorMessage(err, "Credenciales inválidas o error de servidor"),
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-wrapper">
      <form className="login-card" onSubmit={handleSubmit}>
        <h1 className="brand">Kuska Haku</h1>
        <p className="subtitle">Panel administrativo</p>

        <label htmlFor="email">Correo electrónico</label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          autoComplete="username"
          required
        />

        <label htmlFor="password">Contraseña</label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          autoComplete="current-password"
          required
        />

        {error && <p className="error">{error}</p>}

        <button type="submit" className="primary" disabled={isLoading}>
          {isLoading ? "Ingresando..." : "Ingresar"}
        </button>
      </form>
    </div>
  );
}
