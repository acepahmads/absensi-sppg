package model

import (
	"time"
)

// GORM models (dari ERD)
//id	user_id	last_status	last_received	parameter_flat	parameter_zero	parameter_anomaly	ispu	status_mutu_air	validate	accuweather_values	created_at	status	note

type ComplainHandling struct {
	UserID            int       `gorm:"column:user_id" json:"user_id" db:"user_id"`
	LastStatus        *string   `gorm:"column:last_status" json:"last_status" db:"last_status"`
	LastReceived      *string   `gorm:"column:last_received" json:"last_received" db:"last_received"`
	ParameterFlat     *string   `gorm:"column:parameter_flat" json:"parameter_flat" db:"parameter_flat"`
	ParameterZero     *string   `gorm:"column:parameter_zero" json:"parameter_zero" db:"parameter_zero"`
	ParameterAnomaly  *string   `gorm:"column:parameter_anomaly" json:"parameter_anomaly" db:"parameter_anomaly"`
	ISPU              *string   `gorm:"column:ispu" json:"ispu" db:"ispu"`
	StatusMutuAir     *string   `gorm:"column:status_mutu_air" json:"status_mutu_air" db:"status_mutu_air"`
	Validate          *string   `gorm:"column:validate" json:"validate" db:"validate"`
	AccuweatherValues *string   `gorm:"column:accuweather_values" json:"accuweather_values" db:"accuweather_values"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at" db:"created_at"`
	Status            *string   `gorm:"column:status" json:"status" db:"status"`
	Note              *string   `gorm:"column:note" json:"note" db:"note"`
	SiteName          *string   `gorm:"column:site_name" json:"site_name" db:"klh_name"`
	Type              *string   `gorm:"column:type" json:"type" db:"type"`
}
