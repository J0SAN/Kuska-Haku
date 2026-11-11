import { QueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";

// Defaults del panel: sin reintentos en errores de cliente (4xx) y caché
// corta para que los datos operativos se refresquen al navegar.
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      refetchOnWindowFocus: false,
      retry: (failureCount, error) => {
        if (isAxiosError(error)) {
          const status = error.response?.status ?? 0;
          if (status >= 400 && status < 500) {
            return false;
          }
        }
        return failureCount < 2;
      },
    },
  },
});
