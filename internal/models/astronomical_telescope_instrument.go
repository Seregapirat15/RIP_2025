package models

import "time"

// AstronomicalTelescopeInstrument представляет астрономический телескоп/инструмент
type AstronomicalTelescopeInstrument struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	FullName        string    `json:"full_name" db:"full_name"`
	Type            string    `json:"type" db:"type"` // 'ground' или 'space'
	Description     string    `json:"description" db:"description"`
	Accuracy        float64   `json:"accuracy" db:"accuracy"`
	AccuracyUnit    string    `json:"accuracy_unit" db:"accuracy_unit"`
	VelocityPrecision float64 `json:"velocity_precision" db:"velocity_precision"` // Точность измерения скорости (м/с)
	WavelengthRangeMin float64 `json:"wavelength_range_min" db:"wavelength_range_min"` // Минимальная длина волны (нм)
	WavelengthRangeMax float64 `json:"wavelength_range_max" db:"wavelength_range_max"` // Максимальная длина волны (нм)
	ResolutionPower  int       `json:"resolution_power" db:"resolution_power"` // Разрешающая способность
	SpectralResolution float64 `json:"spectral_resolution" db:"spectral_resolution"` // Спектральное разрешение
	Location        string    `json:"location" db:"location"`
	Status          string    `json:"status" db:"status"`
	LaunchDate      string    `json:"launch_date" db:"launch_date"`
	MeasurementRange string   `json:"measurement_range" db:"measurement_range"`
	Resolution      string    `json:"resolution" db:"resolution"`
	Calibration     string    `json:"calibration" db:"calibration"`
	Stability       string    `json:"stability" db:"stability"`
	InstrumentType  string    `json:"instrument_type" db:"instrument_type"`
	ImageURL        *string   `json:"image_url" db:"image_url"`
	IsDeleted       bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

