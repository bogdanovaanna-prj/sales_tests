package client

type ScheduleItem struct {
	Day  string `json:"day"`
	From string `json:"from"`
	To   string `json:"to"`
}

type Sale struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Timezone string         `json:"timezone"`
	Schedule []ScheduleItem `json:"schedule"`
}

type CreateSaleRequest struct {
	Name     string         `json:"name"`
	Timezone string         `json:"timezone"`
	Schedule []ScheduleItem `json:"schedule"`
}

type CreateSaleResponse struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Timezone string         `json:"timezone"`
	Schedule []ScheduleItem `json:"schedule"`
}

type IsActiveResponse struct {
	SaleID string `json:"sale_id"`
	Active bool   `json:"active"`
	At     string `json:"at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
