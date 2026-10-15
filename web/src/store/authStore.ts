import { create } from "zustand";
import { jwtDecode } from "jwt-decode";

// Interfaz basada en el registro de usuarios de nuestro backend Go
interface JwtPayload {
  exp: number; // Fecha de expiración del token
  sub: string; // ID del usuario
  role: string; // Rol del usuario (necesitamos que sea "admin")
  email: string;
}

interface AuthState {
  token: string | null;
  user: JwtPayload | null;
  isAuthenticated: boolean;
  isAdmin: boolean;
  login: (token: string) => void;
  logout: () => void;
}

// Función auxiliar para inicializar el estado leyendo el localStorage
const initializeState = () => {
  const token = localStorage.getItem("kuskahaku_token");
  if (token) {
    try {
      const decodedUser = jwtDecode<JwtPayload>(token);
      // Validar si el token ya expiró
      const currentTime = Date.now() / 1000;
      if (decodedUser.exp < currentTime) {
        localStorage.removeItem("kuskahaku_token");
        return { token: null, user: null };
      }
      return { token, user: decodedUser };
    } catch (error) {
      console.error("Error decodificando el token:", error);
      localStorage.removeItem("kuskahaku_token");
    }
  }
  return { token: null, user: null };
};

export const useAuthStore = create<AuthState>((set) => {
  const initialState = initializeState();

  return {
    token: initialState.token,
    user: initialState.user,
    isAuthenticated: !!initialState.token,
    isAdmin: initialState.user?.role === "admin",

    // Se llamará cuando el endpoint /api/v1/auth/login responda con éxito
    login: (newToken: string) => {
      localStorage.setItem("kuskahaku_token", newToken);
      const decodedUser = jwtDecode<JwtPayload>(newToken);

      set({
        token: newToken,
        user: decodedUser,
        isAuthenticated: true,
        isAdmin: decodedUser.role === "admin",
      });
    },

    // Se llamará al presionar "Cerrar sesión" o cuando el interceptor de Axios detecte un 401
    logout: () => {
      localStorage.removeItem("kuskahaku_token");
      set({
        token: null,
        user: null,
        isAuthenticated: false,
        isAdmin: false,
      });
    },
  };
});
