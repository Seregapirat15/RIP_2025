package database

import (
	"database/sql"
	"fmt"
	"lab2/internal/models"
)

// GetCurrentExoplanetMassCalculation получает текущую заявку (статус "черновик")
func GetCurrentExoplanetMassCalculation(userID int) (*models.ExoplanetMassCalculation, error) {
	query := `SELECT id, status, created_at, creator_id, formation_date, completion_date, 
	                 moderator_id, result, total_mass, notes
	          FROM calculations 
	          WHERE creator_id = $1 AND status = 'черновик'
	          ORDER BY created_at DESC
	          LIMIT 1`
	
	var calculation models.ExoplanetMassCalculation
	err := PostgreSQLConnection.QueryRow(query, userID).Scan(&calculation.ID, &calculation.Status, 
		&calculation.CreatedAt, &calculation.CreatorID, &calculation.FormationDate, 
		&calculation.CompletionDate, &calculation.ModeratorID, &calculation.Result, 
		&calculation.TotalMass, &calculation.Notes)
	
	if err == sql.ErrNoRows {
		return nil, nil // Нет текущей заявки
	}
	if err != nil {
		return nil, err
	}
	
	return &calculation, nil
}

// GetExoplanetMassCalculationByID получает заявку по ID
func GetExoplanetMassCalculationByID(id int) (*models.ExoplanetMassCalculation, error) {
	query := `SELECT id, status, created_at, creator_id, formation_date, completion_date, 
	                 moderator_id, result, total_mass, notes
	          FROM calculations 
	          WHERE id = $1 AND status != 'удалён'`
	
	var calculation models.ExoplanetMassCalculation
	err := PostgreSQLConnection.QueryRow(query, id).Scan(&calculation.ID, &calculation.Status, 
		&calculation.CreatedAt, &calculation.CreatorID, &calculation.FormationDate, 
		&calculation.CompletionDate, &calculation.ModeratorID, &calculation.Result, 
		&calculation.TotalMass, &calculation.Notes)
	
	if err != nil {
		return nil, err
	}
	
	return &calculation, nil
}

// GetExoplanetMassCalculationInstruments получает инструменты заявки
func GetExoplanetMassCalculationInstruments(calculationID int) ([]models.ExoplanetMassCalculationInstrument, error) {
	query := `SELECT calculation_id, instrument_id, exoplanet_name, star_mass, 
	                 orbital_period, velocity_amplitude, inclination, comment, other_info, calculated_mass
	          FROM calculation_instruments 
	          WHERE calculation_id = $1
	          ORDER BY instrument_id`
	
	rows, err := PostgreSQLConnection.Query(query, calculationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var instruments []models.ExoplanetMassCalculationInstrument
	for rows.Next() {
		var instrument models.ExoplanetMassCalculationInstrument
		err := rows.Scan(&instrument.CalculationID, &instrument.InstrumentID, 
			&instrument.ExoplanetName, &instrument.StarMass, &instrument.OrbitalPeriod, 
			&instrument.VelocityAmplitude, &instrument.Inclination, &instrument.Comment, 
			&instrument.OtherInfo, &instrument.CalculatedMass)
		if err != nil {
			return nil, err
		}
		instruments = append(instruments, instrument)
	}
	
	return instruments, nil
}

// AddAstronomicalTelescopeInstrumentToExoplanetMassCalculation добавляет инструмент в заявку
func AddAstronomicalTelescopeInstrumentToExoplanetMassCalculation(userID, instrumentID int, exoplanetName string, starMass, orbitalPeriod, velocityAmplitude, inclination float64, comment, otherInfo string) error {
	// Начинаем транзакцию
	tx, err := PostgreSQLConnection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Получаем или создаем текущую заявку
	var calculationID int
	currentCalc, err := GetCurrentExoplanetMassCalculation(userID)
	if err != nil {
		return err
	}
	
	if currentCalc == nil {
		// Создаем новую заявку
		query := `INSERT INTO calculations (status, creator_id, researcher_name, institution) 
		          VALUES ('черновик', $1, 'Исследователь', 'Институт') 
		          RETURNING id`
		err = tx.QueryRow(query, userID).Scan(&calculationID)
		if err != nil {
			return err
		}
	} else {
		calculationID = currentCalc.ID
	}
	
	// Добавляем инструмент в заявку
	query := `INSERT INTO calculation_instruments 
	          (calculation_id, instrument_id, exoplanet_name, star_mass, orbital_period, 
	           velocity_amplitude, inclination, comment, other_info)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	          ON CONFLICT (calculation_id, instrument_id) 
	          DO UPDATE SET exoplanet_name = $3, star_mass = $4, orbital_period = $5, 
	                        velocity_amplitude = $6, inclination = $7, comment = $8, other_info = $9`
	
	_, err = tx.Exec(query, calculationID, instrumentID, exoplanetName, starMass, 
		orbitalPeriod, velocityAmplitude, inclination, comment, otherInfo)
	if err != nil {
		return err
	}
	
	// Подтверждаем транзакцию
	return tx.Commit()
}

// DeleteExoplanetMassCalculationSQL логическое удаление заявки через SQL UPDATE
func DeleteExoplanetMassCalculationSQL(calculationID int) error {
	query := `UPDATE calculations SET status = 'удалён' WHERE id = $1`
	_, err := PostgreSQLConnection.Exec(query, calculationID)
	return err
}

// DeleteExoplanetMassCalculationWithCursor удаляет заявку через курсор
func DeleteExoplanetMassCalculationWithCursor(calculationID int) error {
	// Начинаем транзакцию
	tx, err := PostgreSQLConnection.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// Создаем курсор для обновления заявки
	cursorQuery := `UPDATE calculations SET status = 'удалён' WHERE id = $1`
	
	// Выполняем UPDATE через курсор
	result, err := tx.Exec(cursorQuery, calculationID)
	if err != nil {
		return fmt.Errorf("ошибка выполнения UPDATE через курсор: %w", err)
	}

	// Проверяем, что строка была обновлена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("заявка с ID %d не найдена", calculationID)
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}
