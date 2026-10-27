package domain

// DayStats representa el total consolidado de un día. Los días sin registros
// se rellenan con 0 en SQL para evitar huecos en la gráfica.
type DayStats struct {
	Date       string  `json:"date"`
	Passengers int     `json:"passengers"`
	Amount     float64 `json:"amount"`
}

type BusStats struct {
	BusCode    string  `json:"busCode"`
	Passengers int     `json:"passengers"`
	Amount     float64 `json:"amount"`
}

type Dashboard struct {
	Date            string     `json:"date"`
	Days            int        `json:"days"`
	ActiveBuses     int        `json:"activeBuses"`
	ActiveRoutes    int        `json:"activeRoutes"`
	PassengersToday int        `json:"passengersToday"`
	IncomeToday     float64    `json:"incomeToday"`
	NewSuggestions  int        `json:"newSuggestions"`
	Series          []DayStats `json:"series"`
	TopBuses        []BusStats `json:"topBuses"`
}
