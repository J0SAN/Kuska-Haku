import type { VehicleFeedbackFilter } from "../vehicleFeedback.types";

export const vehicleFeedbackKeys = {
  all: ["vehicle-feedback"] as const,
  qrList: () => [...vehicleFeedbackKeys.all, "qr-list"] as const,
  feedbackList: (filter: VehicleFeedbackFilter) =>
    [...vehicleFeedbackKeys.all, "feedback-list", filter] as const,
  summary: (days: number) =>
    [...vehicleFeedbackKeys.all, "summary", days] as const,
  publicVehicle: (token: string) =>
    [...vehicleFeedbackKeys.all, "public", token] as const,
};
