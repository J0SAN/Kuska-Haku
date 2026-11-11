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
