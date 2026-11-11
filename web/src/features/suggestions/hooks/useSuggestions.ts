import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { dashboardKeys } from "@/features/dashboard/api/dashboard.keys";
import { fetchSuggestions, markSuggestionRead } from "../api/suggestions.api";
import { suggestionsKeys } from "../api/suggestions.keys";

export function useSuggestions(onlyNew: boolean) {
  return useQuery({
    queryKey: suggestionsKeys.list(onlyNew),
    queryFn: () => fetchSuggestions(onlyNew),
  });
}

export function useMarkSuggestionRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: markSuggestionRead,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: suggestionsKeys.all });
      void queryClient.invalidateQueries({ queryKey: dashboardKeys.all });
    },
  });
}
