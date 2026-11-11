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
