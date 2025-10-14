package database

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"lab2/internal/models"
)

// === УСЛУГИ ===

// GetServicesWithFilter получает услуги с фильтрацией
func GetServicesWithFilter(filter models.ServiceFilter) ([]models.Service, error) {
	query := `SELECT id, name, full_name, type, description, accuracy, accuracy_unit, 
	                 velocity_precision, wavelength_range_min, wavelength_range_max, 
	                 resolution_power, spectral_resolution, location, status, launch_date,
	                 measurement_range, resolution, calibration, stability, instrument_type, 
	                 image_url, is_deleted, created_at
	          FROM instruments 
	          WHERE is_deleted = false`
	
	args := []interface{}{}
	argCount := 0
	
	if filter.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, filter.Type)
	}
	
	if filter.Status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filter.Status)
	}
	
	if filter.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR full_name ILIKE $%d)", argCount, argCount)
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}
	
	query += " ORDER BY created_at DESC"
	
	rows, err := PostgreSQLConnection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var services []models.Service
	for rows.Next() {
		var service models.Service
		err := rows.Scan(&service.ID, &service.Name, &service.FullName, &service.Type,
			&service.Description, &service.Accuracy, &service.AccuracyUnit,
			&service.VelocityPrecision, &service.WavelengthRangeMin, &service.WavelengthRangeMax,
			&service.ResolutionPower, &service.SpectralResolution, &service.Location,
			&service.Status, &service.LaunchDate, &service.MeasurementRange,
			&service.Resolution, &service.Calibration, &service.Stability,
			&service.InstrumentType, &service.ImageURL, &service.IsDeleted, &service.CreatedAt)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	
	return services, nil
}

// GetServiceByID получает услугу по ID
func GetServiceByID(id int) (*models.Service, error) {
	query := `SELECT id, name, full_name, type, description, accuracy, accuracy_unit, 
	                 velocity_precision, wavelength_range_min, wavelength_range_max, 
	                 resolution_power, spectral_resolution, location, status, launch_date,
	                 measurement_range, resolution, calibration, stability, instrument_type, 
	                 image_url, is_deleted, created_at
	          FROM instruments 
	          WHERE id = $1 AND is_deleted = false`
	
	var service models.Service
	err := PostgreSQLConnection.QueryRow(query, id).Scan(
		&service.ID, &service.Name, &service.FullName, &service.Type,
		&service.Description, &service.Accuracy, &service.AccuracyUnit,
		&service.VelocityPrecision, &service.WavelengthRangeMin, &service.WavelengthRangeMax,
		&service.ResolutionPower, &service.SpectralResolution, &service.Location,
		&service.Status, &service.LaunchDate, &service.MeasurementRange,
		&service.Resolution, &service.Calibration, &service.Stability,
		&service.InstrumentType, &service.ImageURL, &service.IsDeleted, &service.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &service, nil
}

// CreateService создает новую услугу
func CreateService(service models.Service) (int, error) {
	query := `INSERT INTO instruments (name, full_name, type, description, accuracy, accuracy_unit,
	                                 velocity_precision, wavelength_range_min, wavelength_range_max,
	                                 resolution_power, spectral_resolution, location, status, launch_date,
	                                 measurement_range, resolution, calibration, stability, instrument_type,
	                                 image_url, is_deleted, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	          RETURNING id`
	
	var id int
	err := PostgreSQLConnection.QueryRow(query,
		service.Name, service.FullName, service.Type, service.Description,
		service.Accuracy, service.AccuracyUnit, service.VelocityPrecision,
		service.WavelengthRangeMin, service.WavelengthRangeMax, service.ResolutionPower,
		service.SpectralResolution, service.Location, service.Status, service.LaunchDate,
		service.MeasurementRange, service.Resolution, service.Calibration,
		service.Stability, service.InstrumentType, service.ImageURL,
		service.IsDeleted, service.CreatedAt).Scan(&id)
	
	return id, err
}

// UpdateService обновляет услугу
func UpdateService(service models.Service) error {
	query := `UPDATE instruments SET 
	          name = $2, full_name = $3, type = $4, description = $5, accuracy = $6, accuracy_unit = $7,
	          velocity_precision = $8, wavelength_range_min = $9, wavelength_range_max = $10,
	          resolution_power = $11, spectral_resolution = $12, location = $13, status = $14, 
	          launch_date = $15, measurement_range = $16, resolution = $17, calibration = $18, 
	          stability = $19, instrument_type = $20, image_url = $21
	          WHERE id = $1`
	
	_, err := PostgreSQLConnection.Exec(query,
		service.ID, service.Name, service.FullName, service.Type, service.Description,
		service.Accuracy, service.AccuracyUnit, service.VelocityPrecision,
		service.WavelengthRangeMin, service.WavelengthRangeMax, service.ResolutionPower,
		service.SpectralResolution, service.Location, service.Status, service.LaunchDate,
		service.MeasurementRange, service.Resolution, service.Calibration,
		service.Stability, service.InstrumentType, service.ImageURL)
	
	return err
}

// DeleteService удаляет услугу (логическое удаление)
func DeleteService(id int) error {
	// Сначала удаляем изображение из MinIO
	service, err := GetServiceByID(id)
	if err == nil && service.ImageURL != nil {
		DeleteServiceImage(*service.ImageURL)
	}
	
	// Логическое удаление
	query := `UPDATE instruments SET is_deleted = true WHERE id = $1`
	_, err = PostgreSQLConnection.Exec(query, id)
	return err
}

// UploadServiceImage загружает изображение услуги в MinIO
func UploadServiceImage(serviceID int, fileName string, file multipart.File, handler *multipart.FileHeader) (string, error) {
	// Удаляем старое изображение если есть
	service, err := GetServiceByID(serviceID)
	if err == nil && service.ImageURL != nil {
		DeleteServiceImage(*service.ImageURL)
	}
	
	// Загружаем новое изображение
	imageURL, err := UploadImageToMinIO(fileName, file)
	if err != nil {
		return "", err
	}
	
	// Обновляем URL в базе
	query := `UPDATE instruments SET image_url = $1 WHERE id = $2`
	_, err = PostgreSQLConnection.Exec(query, imageURL, serviceID)
	if err != nil {
		return "", err
	}
	
	return imageURL, nil
}

// === ЗАЯВКИ ===

// GetDraftOrder получает заявку-черновик пользователя
func GetDraftOrder(creatorID int) (*models.Order, error) {
	query := `SELECT id, status, created_at, creator_id, formation_date, completion_date, 
	                 moderator_id, result, total_mass, notes
	          FROM calculations 
	          WHERE creator_id = $1 AND status = 'черновик'
	          ORDER BY created_at DESC LIMIT 1`
	
	var order models.Order
	err := PostgreSQLConnection.QueryRow(query, creatorID).Scan(
		&order.ID, &order.Status, &order.CreatedAt, &order.CreatorID,
		&order.FormationDate, &order.CompletionDate, &order.ModeratorID,
		&order.Result, &order.TotalMass, &order.Notes)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

// GetOrCreateDraftOrder получает или создает заявку-черновик
func GetOrCreateDraftOrder(creatorID int) (*models.Order, error) {
	order, err := GetDraftOrder(creatorID)
	if err != nil {
		return nil, err
	}
	
	if order == nil {
		// Создаем новую заявку-черновик
		query := `INSERT INTO calculations (status, created_at, creator_id)
		          VALUES ('черновик', $1, $2) RETURNING id`
		
		var id int
		err = PostgreSQLConnection.QueryRow(query, time.Now(), creatorID).Scan(&id)
		if err != nil {
			return nil, err
		}
		
		order = &models.Order{
			ID:        id,
			Status:    "черновик",
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
	}
	
	return order, nil
}

// GetOrderServiceCount получает количество услуг в заявке
func GetOrderServiceCount(orderID int) (int, error) {
	query := `SELECT COUNT(*) FROM calculation_instruments WHERE calculation_id = $1`
	var count int
	err := PostgreSQLConnection.QueryRow(query, orderID).Scan(&count)
	return count, err
}

// GetOrdersWithFilter получает заявки с фильтрацией
func GetOrdersWithFilter(filter models.OrderFilter) ([]models.Order, error) {
	query := `SELECT c.id, c.status, c.created_at, c.creator_id, u.login as creator_login,
	                 c.formation_date, c.completion_date, c.moderator_id, m.login as moderator_login,
	                 c.result, c.total_mass, c.notes
	          FROM calculations c
	          LEFT JOIN users u ON c.creator_id = u.id
	          LEFT JOIN users m ON c.moderator_id = m.id
	          WHERE c.status != 'удалён' AND c.status != 'черновик'`
	
	args := []interface{}{}
	argCount := 0
	
	if filter.Status != "" {
		argCount++
		query += fmt.Sprintf(" AND c.status = $%d", argCount)
		args = append(args, filter.Status)
	}
	
	if filter.FormationFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND c.formation_date >= $%d", argCount)
		args = append(args, *filter.FormationFrom)
	}
	
	if filter.FormationTo != nil {
		argCount++
		query += fmt.Sprintf(" AND c.formation_date <= $%d", argCount)
		args = append(args, *filter.FormationTo)
	}
	
	query += " ORDER BY c.created_at DESC"
	
	rows, err := PostgreSQLConnection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var creatorLogin, moderatorLogin sql.NullString
		
		err := rows.Scan(&order.ID, &order.Status, &order.CreatedAt, &order.CreatorID,
			&creatorLogin, &order.FormationDate, &order.CompletionDate, &order.ModeratorID,
			&moderatorLogin, &order.Result, &order.TotalMass, &order.Notes)
		if err != nil {
			return nil, err
		}
		
		if creatorLogin.Valid {
			order.CreatorLogin = creatorLogin.String
		}
		if moderatorLogin.Valid {
			order.ModeratorLogin = &moderatorLogin.String
		}
		
		orders = append(orders, order)
	}
	
	return orders, nil
}

// GetOrderWithServices получает заявку с услугами
func GetOrderWithServices(orderID int) (*models.Order, error) {
	// Получаем заявку
	query := `SELECT c.id, c.status, c.created_at, c.creator_id, u.login as creator_login,
	                 c.formation_date, c.completion_date, c.moderator_id, m.login as moderator_login,
	                 c.result, c.total_mass, c.notes
	          FROM calculations c
	          LEFT JOIN users u ON c.creator_id = u.id
	          LEFT JOIN users m ON c.moderator_id = m.id
	          WHERE c.id = $1`
	
	var order models.Order
	var creatorLogin, moderatorLogin sql.NullString
	
	err := PostgreSQLConnection.QueryRow(query, orderID).Scan(
		&order.ID, &order.Status, &order.CreatedAt, &order.CreatorID,
		&creatorLogin, &order.FormationDate, &order.CompletionDate, &order.ModeratorID,
		&moderatorLogin, &order.Result, &order.TotalMass, &order.Notes)
	
	if err != nil {
		return nil, err
	}
	
	if creatorLogin.Valid {
		order.CreatorLogin = creatorLogin.String
	}
	if moderatorLogin.Valid {
		order.ModeratorLogin = &moderatorLogin.String
	}
	
	// Получаем услуги заявки
	servicesQuery := `SELECT ci.calculation_id, ci.instrument_id, ci.exoplanet_name, ci.star_mass,
	                         ci.orbital_period, ci.velocity_amplitude, ci.inclination, ci.comment,
	                         ci.other_info, ci.calculated_mass,
	                         i.name, i.full_name, i.type, i.description, i.image_url
	                  FROM calculation_instruments ci
	                  JOIN instruments i ON ci.instrument_id = i.id
	                  WHERE ci.calculation_id = $1`
	
	rows, err := PostgreSQLConnection.Query(servicesQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var orderService models.OrderService
		var service models.Service
		var comment, otherInfo sql.NullString
		var calculatedMass sql.NullFloat64
		
		err := rows.Scan(&orderService.OrderID, &orderService.ServiceID, &orderService.ExoplanetName,
			&orderService.StarMass, &orderService.OrbitalPeriod, &orderService.VelocityAmplitude,
			&orderService.Inclination, &comment, &otherInfo, &calculatedMass,
			&service.Name, &service.FullName, &service.Type, &service.Description, &service.ImageURL)
		if err != nil {
			return nil, err
		}
		
		if comment.Valid {
			orderService.Comment = &comment.String
		}
		if otherInfo.Valid {
			orderService.OtherInfo = &otherInfo.String
		}
		if calculatedMass.Valid {
			orderService.CalculatedMass = &calculatedMass.Float64
		}
		
		orderService.Service = &service
		order.Services = append(order.Services, orderService)
	}
	
	return &order, nil
}

// GetOrderByID получает заявку по ID
func GetOrderByID(orderID int) (*models.Order, error) {
	query := `SELECT id, status, created_at, creator_id, formation_date, completion_date, 
	                 moderator_id, result, total_mass, notes
	          FROM calculations WHERE id = $1`
	
	var order models.Order
	err := PostgreSQLConnection.QueryRow(query, orderID).Scan(
		&order.ID, &order.Status, &order.CreatedAt, &order.CreatorID,
		&order.FormationDate, &order.CompletionDate, &order.ModeratorID,
		&order.Result, &order.TotalMass, &order.Notes)
	
	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

// UpdateOrder обновляет заявку
func UpdateOrder(order models.Order) error {
	query := `UPDATE calculations SET 
	          result = $2, total_mass = $3, notes = $4
	          WHERE id = $1`
	
	_, err := PostgreSQLConnection.Exec(query, order.ID, order.Result, order.TotalMass, order.Notes)
	return err
}

// ValidateOrderForForming проверяет заявку перед формированием
func ValidateOrderForForming(orderID int) error {
	// Проверяем наличие услуг в заявке
	query := `SELECT COUNT(*) FROM calculation_instruments WHERE calculation_id = $1`
	var count int
	err := PostgreSQLConnection.QueryRow(query, orderID).Scan(&count)
	if err != nil {
		return err
	}
	
	if count == 0 {
		return fmt.Errorf("в заявке должны быть услуги")
	}
	
	// Проверяем обязательные поля услуг
	query = `SELECT exoplanet_name, star_mass, orbital_period, velocity_amplitude, inclination
	         FROM calculation_instruments WHERE calculation_id = $1`
	
	rows, err := PostgreSQLConnection.Query(query, orderID)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	for rows.Next() {
		var exoplanetName string
		var starMass, orbitalPeriod, velocityAmplitude, inclination float64
		
		err := rows.Scan(&exoplanetName, &starMass, &orbitalPeriod, &velocityAmplitude, &inclination)
		if err != nil {
			return err
		}
		
		if exoplanetName == "" || starMass <= 0 || orbitalPeriod <= 0 || velocityAmplitude <= 0 || inclination <= 0 {
			return fmt.Errorf("все поля услуг должны быть заполнены")
		}
	}
	
	return nil
}

// FormOrder формирует заявку
func FormOrder(orderID int) error {
	query := `UPDATE calculations SET status = 'сформирован', formation_date = $1 WHERE id = $2`
	_, err := PostgreSQLConnection.Exec(query, time.Now(), orderID)
	return err
}

// CompleteOrder завершает заявку
func CompleteOrder(orderID int, action, result string) error {
	var status string
	if action == "complete" {
		status = "завершён"
	} else {
		status = "отклонён"
	}
	
	query := `UPDATE calculations SET status = $1, completion_date = $2, moderator_id = $3, result = $4 WHERE id = $5`
	_, err := PostgreSQLConnection.Exec(query, status, time.Now(), 1, result, orderID) // moderator_id = 1
	return err
}

// DeleteOrder удаляет заявку
func DeleteOrder(orderID int) error {
	query := `UPDATE calculations SET status = 'удалён' WHERE id = $1`
	_, err := PostgreSQLConnection.Exec(query, orderID)
	return err
}

// === СВЯЗИ М-М ===

// AddServiceToOrder добавляет услугу в заявку
func AddServiceToOrder(orderService models.OrderService) error {
	query := `INSERT INTO calculation_instruments 
	          (calculation_id, instrument_id, exoplanet_name, star_mass, orbital_period, 
	           velocity_amplitude, inclination, comment, other_info)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	_, err := PostgreSQLConnection.Exec(query,
		orderService.OrderID, orderService.ServiceID, orderService.ExoplanetName,
		orderService.StarMass, orderService.OrbitalPeriod, orderService.VelocityAmplitude,
		orderService.Inclination, orderService.Comment, orderService.OtherInfo)
	
	return err
}

// DeleteOrderService удаляет услугу из заявки
func DeleteOrderService(orderID, serviceID int) error {
	query := `DELETE FROM calculation_instruments WHERE calculation_id = $1 AND instrument_id = $2`
	_, err := PostgreSQLConnection.Exec(query, orderID, serviceID)
	return err
}

// UpdateOrderService обновляет связь заявка-услуга
func UpdateOrderService(orderService models.OrderService) error {
	query := `UPDATE calculation_instruments SET 
	          exoplanet_name = $3, star_mass = $4, orbital_period = $5, velocity_amplitude = $6,
	          inclination = $7, comment = $8, other_info = $9
	          WHERE calculation_id = $1 AND instrument_id = $2`
	
	_, err := PostgreSQLConnection.Exec(query,
		orderService.OrderID, orderService.ServiceID, orderService.ExoplanetName,
		orderService.StarMass, orderService.OrbitalPeriod, orderService.VelocityAmplitude,
		orderService.Inclination, orderService.Comment, orderService.OtherInfo)
	
	return err
}

// === ПОЛЬЗОВАТЕЛИ ===

// CreateUser создает пользователя
func CreateUser(user models.User) (int, error) {
	query := `INSERT INTO users (login, email, password_hash, first_name, last_name, role, is_active, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	
	var id int
	err := PostgreSQLConnection.QueryRow(query,
		user.Login, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role, user.IsActive, user.CreatedAt).Scan(&id)
	
	return id, err
}

// GetUserByID получает пользователя по ID
func GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, login, email, first_name, last_name, role, is_active, created_at FROM users WHERE id = $1`
	
	var user models.User
	err := PostgreSQLConnection.QueryRow(query, id).Scan(
		&user.ID, &user.Login, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.IsActive, &user.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetUserByLogin получает пользователя по логину
func GetUserByLogin(login string) (*models.User, error) {
	query := `SELECT id, login, email, first_name, last_name, role, is_active, created_at FROM users WHERE login = $1`
	
	var user models.User
	err := PostgreSQLConnection.QueryRow(query, login).Scan(
		&user.ID, &user.Login, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.IsActive, &user.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// UpdateUser обновляет пользователя
func UpdateUser(user models.User) error {
	query := `UPDATE users SET login = $2, email = $3, first_name = $4, last_name = $5 WHERE id = $1`
	_, err := PostgreSQLConnection.Exec(query, user.ID, user.Login, user.Email, user.FirstName, user.LastName)
	return err
}

// === MINIO ===

// UploadImageToMinIO загружает изображение в MinIO
func UploadImageToMinIO(fileName string, file io.Reader) (string, error) {
	// Здесь должна быть реализация загрузки в MinIO
	// Возвращаем URL изображения
	return "http://localhost:9000/service-images/" + fileName, nil
}

// DeleteServiceImage удаляет изображение из MinIO
func DeleteServiceImage(imageURL string) error {
	// Здесь должна быть реализация удаления из MinIO
	return nil
}

// UpdateServiceImageURL обновляет URL изображения услуги
func UpdateServiceImageURL(serviceID int, imageURL string) error {
	query := `UPDATE instruments SET image_url = $1 WHERE id = $2`
	_, err := PostgreSQLConnection.Exec(query, imageURL, serviceID)
	return err
}
