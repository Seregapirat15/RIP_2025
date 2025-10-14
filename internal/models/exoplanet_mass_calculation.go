package models

import "time"

// ExoplanetMassCalculation представляет расчет массы экзопланеты
type ExoplanetMassCalculation struct {
	ID              int        `json:"id" db:"id"`
	Status          string     `json:"status" db:"status"` // черновик, удалён, сформирован, завершён, отклонён
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	CreatorID       int        `json:"creator_id" db:"creator_id"`
	FormationDate   *time.Time `json:"formation_date" db:"formation_date"`
	CompletionDate  *time.Time `json:"completion_date" db:"completion_date"`
	ModeratorID     *int       `json:"moderator_id" db:"moderator_id"`
	Result          *string    `json:"result" db:"result"`
	TotalMass       *float64   `json:"total_mass" db:"total_mass"`
	Notes           *string    `json:"notes" db:"notes"`
}

// ExoplanetMassCalculationInstrument представляет связь между расчетом и инструментом
type ExoplanetMassCalculationInstrument struct {
	CalculationID    int      `json:"calculation_id" db:"calculation_id"`
	InstrumentID     int      `json:"instrument_id" db:"instrument_id"`
	ExoplanetName    string   `json:"exoplanet_name" db:"exoplanet_name"`
	StarMass         float64  `json:"star_mass" db:"star_mass"`
	OrbitalPeriod    float64  `json:"orbital_period" db:"orbital_period"`
	VelocityAmplitude float64 `json:"velocity_amplitude" db:"velocity_amplitude"`
	Inclination      float64  `json:"inclination" db:"inclination"`
	Comment          *string  `json:"comment" db:"comment"`
	OtherInfo        *string  `json:"other_info" db:"other_info"`
	CalculatedMass   *float64 `json:"calculated_mass" db:"calculated_mass"`
}
