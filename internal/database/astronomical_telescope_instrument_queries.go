package database

import (
	"lab4/internal/models"
)

// GetAstronomicalTelescopeInstruments получает все активные астрономические инструменты
func GetAstronomicalTelescopeInstruments(searchQuery string) ([]models.AstronomicalTelescopeInstrument, error) {
	var query string
	var args []interface{}
	
	if searchQuery != "" {
		query = `SELECT id, name, full_name, type, description, accuracy, accuracy_unit, 
		                location, status, launch_date, measurement_range, resolution, 
		                calibration, stability, instrument_type, image_url, is_deleted, created_at
		         FROM instruments 
		         WHERE is_deleted = false AND (name ILIKE $1 OR full_name ILIKE $1)
		         ORDER BY name`
		args = []interface{}{"%" + searchQuery + "%"}
	} else {
		query = `SELECT id, name, full_name, type, description, accuracy, accuracy_unit, 
		                location, status, launch_date, measurement_range, resolution, 
		                calibration, stability, instrument_type, image_url, is_deleted, created_at
		         FROM instruments 
		         WHERE is_deleted = false
		         ORDER BY name`
		args = []interface{}{}
	}
	
	rows, err := PostgreSQLConnection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var instruments []models.AstronomicalTelescopeInstrument
	for rows.Next() {
		var instrument models.AstronomicalTelescopeInstrument
		err := rows.Scan(&instrument.ID, &instrument.Name, &instrument.FullName, 
			&instrument.Type, &instrument.Description, &instrument.Accuracy, 
			&instrument.AccuracyUnit, &instrument.Location, &instrument.Status, 
			&instrument.LaunchDate, &instrument.MeasurementRange, &instrument.Resolution, 
			&instrument.Calibration, &instrument.Stability, &instrument.InstrumentType, 
			&instrument.ImageURL, &instrument.IsDeleted, &instrument.CreatedAt)
		if err != nil {
			return nil, err
		}
		instruments = append(instruments, instrument)
	}
	
	return instruments, nil
}

// GetAstronomicalTelescopeInstrumentByID получает инструмент по ID
func GetAstronomicalTelescopeInstrumentByID(id int) (*models.AstronomicalTelescopeInstrument, error) {
	query := `SELECT id, name, full_name, type, description, accuracy, accuracy_unit, 
	                 location, status, launch_date, measurement_range, resolution, 
	                 calibration, stability, instrument_type, image_url, is_deleted, created_at
	          FROM instruments 
	          WHERE id = $1 AND is_deleted = false`
	
	var instrument models.AstronomicalTelescopeInstrument
	err := PostgreSQLConnection.QueryRow(query, id).Scan(&instrument.ID, &instrument.Name, &instrument.FullName, 
		&instrument.Type, &instrument.Description, &instrument.Accuracy, 
		&instrument.AccuracyUnit, &instrument.Location, &instrument.Status, 
		&instrument.LaunchDate, &instrument.MeasurementRange, &instrument.Resolution, 
		&instrument.Calibration, &instrument.Stability, &instrument.InstrumentType, 
		&instrument.ImageURL, &instrument.IsDeleted, &instrument.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &instrument, nil
}
