import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { dashboardKeys } from "@/features/dashboard/api/dashboard.keys";
import {
  fetchFeedbackSummary,
  fetchPublicVehicle,
  fetchVehicleFeedback,
  fetchVehicleQrList,
  generateVehicleQr,
  revokeVehicleQr,
  rotateVehicleQr,
  submitPublicFeedback,
  updateFeedbackStatus,
} from "../api/vehicleFeedback.api";
import { vehicleFeedbackKeys } from "../api/vehicleFeedback.keys";
import type {
  FeedbackStatus,
  SubmitFeedbackInput,
  VehicleFeedbackFilter,
} from "../vehicleFeedback.types";

export function useVehicleQrList() {
  return useQuery({
    queryKey: vehicleFeedbackKeys.qrList(),
    queryFn: fetchVehicleQrList,
  });
}

export function useGenerateVehicleQr() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: generateVehicleQr,
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: vehicleFeedbackKeys.qrList(),
      });
    },
  });
}

export function useRotateVehicleQr() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: rotateVehicleQr,
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: vehicleFeedbackKeys.qrList(),
      });
    },
  });
}

export function useRevokeVehicleQr() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: revokeVehicleQr,
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: vehicleFeedbackKeys.qrList(),
      });
    },
  });
}

export function useVehicleFeedback(filter: VehicleFeedbackFilter) {
  return useQuery({
    queryKey: vehicleFeedbackKeys.feedbackList(filter),
    queryFn: () => fetchVehicleFeedback(filter),
  });
}

export function useFeedbackSummary(days: number) {
  return useQuery({
    queryKey: vehicleFeedbackKeys.summary(days),
    queryFn: () => fetchFeedbackSummary(days),
  });
}

export function useUpdateFeedbackStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: FeedbackStatus }) =>
      updateFeedbackStatus(id, status),
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: [...vehicleFeedbackKeys.all, "feedback-list"],
      });
      void queryClient.invalidateQueries({
        queryKey: [...vehicleFeedbackKeys.all, "summary"],
      });
      void queryClient.invalidateQueries({ queryKey: dashboardKeys.all });
    },
  });
}

export function usePublicVehicle(token: string | undefined) {
  return useQuery({
    queryKey: vehicleFeedbackKeys.publicVehicle(token ?? ""),
    queryFn: () => fetchPublicVehicle(token ?? ""),
    enabled: Boolean(token),
    retry: false,
  });
}

export function useSubmitFeedback(token: string | undefined) {
  return useMutation({
    mutationFn: (input: SubmitFeedbackInput) =>
      submitPublicFeedback(token ?? "", input),
  });
}
