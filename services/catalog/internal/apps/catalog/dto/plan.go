package dto

type CreatePlanDTO struct {
	Name          string  `json:"name" validate:"required"`
	PlanType      string  `json:"plan_type" validate:"required"` // subscription | package
	Price         float64 `json:"price"`
	PriceCurrency string  `json:"price_currency"`
	BillingCycle  string  `json:"billing_cycle"` // monthly | yearly (subscription only)
	SessionCount  *int    `json:"session_count"` // nil = unlimited; per-cycle for subscription, total for package
	ValidityDays  *int    `json:"validity_days"` // package only: days valid after purchase
}

type UpdatePlanDTO = CreatePlanDTO
