package models

import "time"

// Service представляет услугу (астрономический инструмент)
type Service struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	FullName        string    `json:"full_name" db:"full_name"`
	Type            string    `json:"type" db:"type"` // 'ground' или 'space'
	Description     string    `json:"description" db:"description"`
	Accuracy        float64   `json:"accuracy" db:"accuracy"`
	AccuracyUnit    string    `json:"accuracy_unit" db:"accuracy_unit"`
	VelocityPrecision float64 `json:"velocity_precision" db:"velocity_precision"`
	WavelengthRangeMin float64 `json:"wavelength_range_min" db:"wavelength_range_min"`
	WavelengthRangeMax float64 `json:"wavelength_range_max" db:"wavelength_range_max"`
	ResolutionPower  int       `json:"resolution_power" db:"resolution_power"`
	SpectralResolution float64 `json:"spectral_resolution" db:"spectral_resolution"`
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

// Order представляет заявку на расчет массы экзопланеты
type Order struct {
	ID              int        `json:"id" db:"id"`
	Status          string     `json:"status" db:"status"` // черновик, сформирован, завершён, отклонён, удалён
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	CreatorID       int        `json:"creator_id" db:"creator_id"`
	CreatorLogin    string     `json:"creator_login" db:"creator_login"`
	FormationDate   *time.Time `json:"formation_date" db:"formation_date"`
	CompletionDate  *time.Time `json:"completion_date" db:"completion_date"`
	ModeratorID     *int       `json:"moderator_id" db:"moderator_id"`
	ModeratorLogin  *string    `json:"moderator_login" db:"moderator_login"`
	Result          *string    `json:"result" db:"result"`
	TotalMass       *float64   `json:"total_mass" db:"total_mass"`
	Notes           *string    `json:"notes" db:"notes"`
	Services        []OrderService `json:"services,omitempty"` // Услуги в заявке
}

// OrderService представляет связь заявки с услугой (м-м)
type OrderService struct {
	OrderID           int      `json:"order_id" db:"order_id"`
	ServiceID         int      `json:"service_id" db:"service_id"`
	ExoplanetName     string   `json:"exoplanet_name" db:"exoplanet_name"`
	StarMass          float64  `json:"star_mass" db:"star_mass"`
	OrbitalPeriod     float64  `json:"orbital_period" db:"orbital_period"`
	VelocityAmplitude float64  `json:"velocity_amplitude" db:"velocity_amplitude"`
	Inclination       float64  `json:"inclination" db:"inclination"`
	Eccentricity      float64  `json:"eccentricity" db:"eccentricity"`
	Comment           *string  `json:"comment" db:"comment"`
	OtherInfo         *string  `json:"other_info" db:"other_info"`
	CalculatedMass    *float64 `json:"calculated_mass" db:"calculated_mass"`
	Service           *Service `json:"service,omitempty"` // Детали услуги
}

// User представляет пользователя системы
type User struct {
	ID        int       `json:"id" db:"id"`
	Login     string    `json:"login" db:"login"`
	Email     string    `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Role      string    `json:"role" db:"role"` // user, moderator, admin
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CartIcon представляет иконку корзины
type CartIcon struct {
	OrderID     int `json:"order_id"`
	ServiceCount int `json:"service_count"`
}

// ServiceFilter представляет фильтр для услуг
type ServiceFilter struct {
	Type     string `json:"type"`     // ground, space
	Status   string `json:"status"`   // Активен, Неактивен
	Search   string `json:"search"`   // Поиск по названию
}

// OrderFilter представляет фильтр для заявок
type OrderFilter struct {
	Status        string    `json:"status"`         // сформирован, завершён, отклонён
	FormationFrom *time.Time `json:"formation_from"` // Дата формирования от
	FormationTo   *time.Time `json:"formation_to"`   // Дата формирования до
}

