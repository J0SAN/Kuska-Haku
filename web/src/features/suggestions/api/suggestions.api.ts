import httpClient from "@/shared/api/httpClient";
import type { SuggestionData } from "../suggestions.types";

// Sugerencias generales del servicio (no están asociadas a ningún vehículo).
export async function fetchSuggestions(
  onlyNew: boolean,
): Promise<SuggestionData[]> {
  const { data } = await httpClient.get<SuggestionData[]>(
    "/admin/suggestions",
    { params: onlyNew ? { only: "new" } : undefined },
  );
  return data ?? [];
}

export async function markSuggestionRead(id: string): Promise<void> {
  await httpClient.put(`/admin/suggestions/${id}/read`);
}
