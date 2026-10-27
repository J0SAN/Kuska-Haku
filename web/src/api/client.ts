import axios from "axios";

export interface Coordinate {
  latitude: number;
  longitude: number;
}

export interface RouteData {
  id: string;
  name: string;
  path: Coordinate[];
  createdAt?: string;
  updatedAt?: string;
}

export interface StopData {
  id: string;
  routeId: string;
  name: string;
  location: Coordinate;
  stopOrder: number;
  createdAt?: string;
}

export interface FleetData {
  id: string;
  busId: string;
  routeId: string;
  driverId?: string | null;
  status: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface DayStats {
  date: string;
  passengers: number;
  amount: number;
}

export interface BusStats {
  busCode: string;
  passengers: number;
  amount: number;
}

export interface DashboardData {
  date: string;
  days: number;
  activeBuses: number;
  activeRoutes: number;
  passengersToday: number;
  incomeToday: number;
  newSuggestions: number;
  series: DayStats[];
  topBuses: BusStats[];
}

export interface SuggestionData {
  id: string;
  userId: string;
  userName: string;
  message: string;
  status: "new" | "read";
  createdAt: string;
}

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("kuskahaku_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem("kuskahaku_token");
      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    }
    return Promise.reject(error);
  },
);

export function getErrorMessage(err: unknown, fallback: string): string {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data;
    if (typeof data === "string" && data.trim() !== "") return data;
    if (data?.message) return data.message;
    if (data?.error) return data.error;
  }
  return fallback;
}

export default apiClient;
