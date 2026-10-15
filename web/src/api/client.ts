import axios from "axios";

// 1. Instancia base de Axios
// Vite usa import.meta.env para las variables de entorno
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

// 2. Interceptor de Peticiones (Request)
// Se ejecuta ANTES de que la petición salga hacia el backend en Go
apiClient.interceptors.request.use(
  (config) => {
    // Leeremos el token del localStorage (donde lo guardará Zustand más adelante)
    const token = localStorage.getItem("kuskahaku_token");

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// 3. Interceptor de Respuestas (Response)
// Se ejecuta CUANDO el backend responde, antes de que llegue a nuestros componentes
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Si Go nos devuelve un 401 (Unauthorized), significa que el token expiró o es falso
    if (error.response && error.response.status === 401) {
      console.warn("Sesión expirada o inválida. Cerrando sesión...");
      localStorage.removeItem("kuskahaku_token");
      // Redirigir agresivamente al login
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);

export default apiClient;
