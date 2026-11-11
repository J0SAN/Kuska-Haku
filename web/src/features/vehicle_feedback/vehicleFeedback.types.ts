export interface VehicleQRSummary {
  fleetId: string;
  busCode: string;
  routeId: string;
  routeName: string;
  driverId?: string;
  driverName?: string;
  qrCodeId?: string;
  token?: string;
  createdAt?: string;
}

export interface VehicleQR {
  id: string;
  fleetId: string;
  token: string;
  createdAt: string;
  revokedAt?: string;
}

export type FeedbackStatus = "new" | "read" | "hidden";

export interface VehicleFeedback {
  id: string;
  fleetId?: string;
  qrCodeId?: string;
  busCode: string;
  routeName: string;
  driverId?: string;
  driverName?: string;
  rating: number;
  comment: string;
  serviceDate: string;
  busLive: boolean;
  status: FeedbackStatus;
  createdAt: string;
}

export interface VehicleRatingSummary {
  fleetId?: string;
  busCode: string;
  routeName: string;
  average: number;
  total: number;
}

export interface PublicVehicle {
  busCode: string;
  routeName: string;
}

export interface VehicleFeedbackFilter {
  fleetId?: string;
  status?: FeedbackStatus | "";
}

export interface SubmitFeedbackInput {
  rating: number;
  comment: string;
  deviceId: string;
}
