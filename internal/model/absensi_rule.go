package model

type AbsensiLateRule struct {
	ID               int     `db:"id" json:"id"`
	Code             string  `db:"code" json:"code"`
	MinMinutes       int     `db:"min_minutes" json:"min_minutes"`
	MaxMinutes       int     `db:"max_minutes" json:"max_minutes"`
	DeductionBase    string  `db:"deduction_base" json:"deduction_base"`
	DeductionPercent float64 `db:"deduction_percent" json:"deduction_percent"`
	TenantID         int     `db:"tenant_id" json:"tenant_id"`
}
