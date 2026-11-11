import httpClient from "@/shared/api/httpClient";
import type {
  FeedbackStatus,
  PublicVehicle,
  SubmitFeedbackInput,
  VehicleFeedback,
  VehicleFeedbackFilter,
  VehicleQR,
  VehicleQRSummary,
  VehicleRatingSummary,
} from "../vehicleFeedback.types";

export async function fetchVehicleQrList(): Promise<VehicleQRSummary[]> {
  const { data } =
    await httpClient.get<VehicleQRSummary[]>("/admin/vehicle-qr");
  return data ?? [];
}

export async function generateVehicleQr(
  fleetId: string,
): Promise<VehicleQR> {
  const { data } = await httpClient.post<VehicleQR>(
    `/admin/vehicle-qr/${fleetId}`,
  );
  return data;
}

export async function rotateVehicleQr(fleetId: string): Promise<VehicleQR> {
  const { data } = await httpClient.post<VehicleQR>(
    `/admin/vehicle-qr/${fleetId}/rotate`,
  );
  return data;
}

export async function revokeVehicleQr(fleetId: string): Promise<void> {
  await httpClient.delete(`/admin/vehicle-qr/${fleetId}`);
}

export async function fetchVehicleFeedback(
  filter: VehicleFeedbackFilter,
): Promise<VehicleFeedback[]> {
  const { data } = await httpClient.get<VehicleFeedback[]>(
    "/admin/vehicle-feedback",
    {
      params: {
        fleetId: filter.fleetId || undefined,
        status: filter.status || undefined,
      },
    },
  );
  return data ?? [];
}

export async function fetchFeedbackSummary(
  days: number,
): Promise<VehicleRatingSummary[]> {
  const { data } = await httpClient.get<VehicleRatingSummary[]>(
    "/admin/vehicle-feedback/summary",
    { params: { days } },
  );
  return data ?? [];
}

export async function updateFeedbackStatus(
  id: string,
  status: FeedbackStatus,
): Promise<void> {
  await httpClient.put(`/admin/vehicle-feedback/${id}/status`, { status });
}

export async function fetchPublicVehicle(
  token: string,
): Promise<PublicVehicle> {
  const { data } = await httpClient.get<PublicVehicle>(
    `/public/vehicle-feedback/${token}`,
  );
  return data;
}

export async function submitPublicFeedback(
  token: string,
  input: SubmitFeedbackInput,
): Promise<void> {
  await httpClient.post(`/public/vehicle-feedback/${token}`, input);
}
