export interface FleetData {
  id: string;
  busId: string;
  routeId: string;
  driverId?: string | null;
  status: string;
  createdAt?: string;
  updatedAt?: string;
}
