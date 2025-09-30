package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Модели данных для работы с БД
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	Institution  string    `json:"institution"`
	Email        string    `json:"email"`
	IsModerator  bool      `json:"is_moderator"`
	CreatedAt    time.Time `json:"created_at"`
}

type Instrument struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Type            string    `json:"type"`
	Description     string    `json:"description"`
	Accuracy        float64   `json:"accuracy"`
	AccuracyUnit    string    `json:"accuracy_unit"`
	Location        string    `json:"location"`
	Status          string    `json:"status"`
	LaunchDate      string    `json:"launch_date"`
	MeasurementRange string   `json:"measurement_range"`
	Resolution      string    `json:"resolution"`
	Calibration     string    `json:"calibration"`
	Stability       string    `json:"stability"`
	InstrumentType  string    `json:"instrument_type"`
	ImageURL        *string   `json:"image_url"` // Nullable
	IsDeleted       bool      `json:"is_deleted"`
	CreatedAt       time.Time `json:"created_at"`
}

type Calculation struct {
	ID              int        `json:"id"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatorID       int        `json:"creator_id"`
	FormationDate   *time.Time `json:"formation_date"`
	CompletionDate  *time.Time `json:"completion_date"`
	ModeratorID     *int       `json:"moderator_id"`
	ResearcherName  string     `json:"researcher_name"`
	Institution     string     `json:"institution"`
	Result          string     `json:"result"`
	TotalMass       *float64   `json:"total_mass"`
	Notes           string     `json:"notes"`
}

type CalculationInstrument struct {
	CalculationID    int      `json:"calculation_id"`
	InstrumentID     int      `json:"instrument_id"`
	ExoplanetName    string   `json:"exoplanet_name"`
	StarMass         float64  `json:"star_mass"`
	OrbitalPeriod    float64  `json:"orbital_period"`
	VelocityAmplitude float64 `json:"velocity_amplitude"`
	Inclination      float64  `json:"inclination"`
	Comment          string   `json:"comment"`
	OtherInfo        string   `json:"other_info"`
	CalculatedMass   *float64 `json:"calculated_mass"`
}

// Глобальные переменные
var db *sql.DB
var minioClient *minio.Client

// Конфигурация
const (
	// PostgreSQL
	dbHost     = "postgres" // Имя сервиса в Docker Compose
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres123"
	dbName     = "exoplanet_calculations"
	
	// MinIO
	minioEndpoint        = "minio:9000" // Имя сервиса в Docker Compose
	minioAccessKeyID     = "minioadmin123"
	minioSecretAccessKey = "minioadmin123456"
	minioUseSSL          = false
	minioBucketName      = "telescope-images"
)

// Инициализация подключения к PostgreSQL
func initDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	
	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка ping БД:", err)
	}
	
	log.Println("Подключение к PostgreSQL установлено")
}

// Инициализация MinIO клиента
func initMinIO() {
	var err error
	minioClient, err = minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKeyID, minioSecretAccessKey, ""),
		Secure: minioUseSSL,
	})
	if err != nil {
		log.Fatal("Ошибка инициализации MinIO клиента:", err)
	}

	// Создание bucket если не существует
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, minioBucketName)
	if err != nil {
		log.Fatal("Ошибка проверки существования bucket:", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, minioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal("Ошибка создания bucket:", err)
		}
		log.Printf("Bucket '%s' создан успешно", minioBucketName)
	}
}

// Генерация presigned URL для изображения
func getImageURL(objectName string) (string, error) {
	if objectName == "" {
		return "", nil
	}

	ctx := context.Background()
	presignedURL, err := minioClient.PresignedGetObject(ctx, minioBucketName, objectName, time.Hour*24, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// ORM операции

// Получение всех активных инструментов
func getInstruments(searchQuery string) ([]Instrument, error) {
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
	
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var instruments []Instrument
	for rows.Next() {
		var instrument Instrument
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

// Получение инструмента по ID
func getInstrumentByID(id int) (*Instrument, error) {
	query := `SELECT id, name, full_name, type, description, accuracy, accuracy_unit, 
	                 location, status, launch_date, measurement_range, resolution, 
	                 calibration, stability, instrument_type, image_url, is_deleted, created_at
	          FROM instruments 
	          WHERE id = $1 AND is_deleted = false`
	
	var instrument Instrument
	err := db.QueryRow(query, id).Scan(&instrument.ID, &instrument.Name, &instrument.FullName, 
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

// Получение текущей заявки пользователя (статус "черновик")
func getCurrentCalculation(userID int) (*Calculation, error) {
	query := `SELECT id, status, created_at, creator_id, formation_date, completion_date, 
	                 moderator_id, researcher_name, institution, result, total_mass, notes
	          FROM calculations 
	          WHERE creator_id = $1 AND status = 'черновик'
	          ORDER BY created_at DESC
	          LIMIT 1`
	
	var calculation Calculation
	err := db.QueryRow(query, userID).Scan(&calculation.ID, &calculation.Status, 
		&calculation.CreatedAt, &calculation.CreatorID, &calculation.FormationDate, 
		&calculation.CompletionDate, &calculation.ModeratorID, &calculation.ResearcherName, 
		&calculation.Institution, &calculation.Result, &calculation.TotalMass, &calculation.Notes)
	
	if err == sql.ErrNoRows {
		return nil, nil // Нет текущей заявки
	}
	if err != nil {
		return nil, err
	}
	
	return &calculation, nil
}

// Получение заявки по ID
func getCalculationByID(id int) (*Calculation, error) {
	query := `SELECT id, status, created_at, creator_id, formation_date, completion_date, 
	                 moderator_id, researcher_name, institution, result, total_mass, notes
	          FROM calculations 
	          WHERE id = $1 AND status != 'удалён'`
	
	var calculation Calculation
	err := db.QueryRow(query, id).Scan(&calculation.ID, &calculation.Status, 
		&calculation.CreatedAt, &calculation.CreatorID, &calculation.FormationDate, 
		&calculation.CompletionDate, &calculation.ModeratorID, &calculation.ResearcherName, 
		&calculation.Institution, &calculation.Result, &calculation.TotalMass, &calculation.Notes)
	
	if err != nil {
		return nil, err
	}
	
	return &calculation, nil
}

// Получение инструментов заявки
func getCalculationInstruments(calculationID int) ([]CalculationInstrument, error) {
	query := `SELECT calculation_id, instrument_id, exoplanet_name, star_mass, 
	                 orbital_period, velocity_amplitude, inclination, comment, other_info, calculated_mass
	          FROM calculation_instruments 
	          WHERE calculation_id = $1
	          ORDER BY instrument_id`
	
	rows, err := db.Query(query, calculationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var instruments []CalculationInstrument
	for rows.Next() {
		var instrument CalculationInstrument
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

// Создание новой заявки или добавление инструмента в существующую
func addInstrumentToCalculation(userID, instrumentID int, exoplanetName string, starMass, orbitalPeriod, velocityAmplitude, inclination float64, comment, otherInfo string) error {
	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Получаем или создаем текущую заявку
	var calculationID int
	currentCalc, err := getCurrentCalculation(userID)
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

// Логическое удаление заявки через SQL UPDATE (без ORM)
func deleteCalculationSQL(calculationID int) error {
	query := `UPDATE calculations SET status = 'удалён' WHERE id = $1`
	_, err := db.Exec(query, calculationID)
	return err
}

// HTTP обработчики

// Главная страница - список инструментов
func instrumentsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/instruments.html"))
	
	// Получение параметра поиска
	searchQuery := r.URL.Query().Get("search")
	
	// Получение инструментов из БД
	instruments, err := getInstruments(searchQuery)
	if err != nil {
		http.Error(w, "Ошибка получения инструментов", http.StatusInternalServerError)
		return
	}
	
	// Получение текущей заявки пользователя (для демонстрации используем userID = 1)
	currentCalc, err := getCurrentCalculation(1)
	if err != nil {
		http.Error(w, "Ошибка получения заявки", http.StatusInternalServerError)
		return
	}
	
	var totalCalculations int
	if currentCalc != nil {
		calcInstruments, err := getCalculationInstruments(currentCalc.ID)
		if err == nil {
			totalCalculations = len(calcInstruments)
		}
	}
	
	// Генерация URL изображений
	instrumentsWithImages := make([]map[string]interface{}, len(instruments))
	for i, instrument := range instruments {
		var imageURL string
		if instrument.ImageURL != nil {
			imageURL, err = getImageURL(*instrument.ImageURL)
			if err != nil {
				log.Printf("Ошибка получения URL изображения для %s: %v", instrument.Name, err)
				imageURL = ""
			}
		}
		
		instrumentsWithImages[i] = map[string]interface{}{
			"Instrument": instrument,
			"ImageURL":   imageURL,
		}
	}
	
	templateData := map[string]interface{}{
		"Instruments":       instrumentsWithImages,
		"TotalCalculations": totalCalculations,
		"CurrentCalculation": currentCalc,
		"SearchQuery":       searchQuery,
	}
	
	err = tmpl.Execute(w, templateData)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		return
	}
}

// Страница детальной информации об инструменте
func instrumentDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Получение инструмента из БД
	instrument, err := getInstrumentByID(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Генерация URL изображения
	var imageURL string
	if instrument.ImageURL != nil {
		imageURL, err = getImageURL(*instrument.ImageURL)
		if err != nil {
			log.Printf("Ошибка получения URL изображения для %s: %v", instrument.Name, err)
			imageURL = ""
		}
	}
	
	tmpl := template.Must(template.ParseFiles("templates/instrument_detail.html"))
	
	templateData := map[string]interface{}{
		"Instrument": instrument,
		"ImageURL":   imageURL,
	}
	
	err = tmpl.Execute(w, templateData)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		return
	}
}

// Страница просмотра заявки
func calculationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Получение заявки из БД
	calculation, err := getCalculationByID(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Получение инструментов заявки
	instruments, err := getCalculationInstruments(id)
	if err != nil {
		http.Error(w, "Ошибка получения инструментов заявки", http.StatusInternalServerError)
		return
	}
	
	// Получение полной информации об инструментах
	var instrumentsWithDetails []map[string]interface{}
	for _, calcInstrument := range instruments {
		instrument, err := getInstrumentByID(calcInstrument.InstrumentID)
		if err != nil {
			continue
		}
		
		var imageURL string
		if instrument.ImageURL != nil {
			imageURL, err = getImageURL(*instrument.ImageURL)
			if err != nil {
				imageURL = ""
			}
		}
		
		instrumentsWithDetails = append(instrumentsWithDetails, map[string]interface{}{
			"Instrument":        instrument,
			"CalculationData":   calcInstrument,
			"ImageURL":          imageURL,
		})
	}
	
	tmpl := template.Must(template.ParseFiles("templates/calculation.html"))
	
	templateData := map[string]interface{}{
		"Calculation": calculation,
		"Instruments": instrumentsWithDetails,
	}
	
	err = tmpl.Execute(w, templateData)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		return
	}
}

// POST: Добавление инструмента в заявку
func addToCalculationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	
	// Парсинг формы
	instrumentID, err := strconv.Atoi(r.FormValue("instrument_id"))
	if err != nil {
		http.Error(w, "Неверный ID инструмента", http.StatusBadRequest)
		return
	}
	
	exoplanetName := r.FormValue("exoplanet_name")
	starMass, err := strconv.ParseFloat(r.FormValue("star_mass"), 64)
	if err != nil {
		http.Error(w, "Неверная масса звезды", http.StatusBadRequest)
		return
	}
	
	orbitalPeriod, err := strconv.ParseFloat(r.FormValue("orbital_period"), 64)
	if err != nil {
		http.Error(w, "Неверный орбитальный период", http.StatusBadRequest)
		return
	}
	
	velocityAmplitude, err := strconv.ParseFloat(r.FormValue("velocity_amplitude"), 64)
	if err != nil {
		http.Error(w, "Неверная амплитуда скорости", http.StatusBadRequest)
		return
	}
	
	inclination, err := strconv.ParseFloat(r.FormValue("inclination"), 64)
	if err != nil {
		http.Error(w, "Неверный наклон орбиты", http.StatusBadRequest)
		return
	}
	
	comment := r.FormValue("comment")
	otherInfo := r.FormValue("other_info")
	
	// Добавление в заявку (для демонстрации используем userID = 1)
	err = addInstrumentToCalculation(1, instrumentID, exoplanetName, starMass, 
		orbitalPeriod, velocityAmplitude, inclination, comment, otherInfo)
	if err != nil {
		http.Error(w, "Ошибка добавления в заявку", http.StatusInternalServerError)
		return
	}
	
	// Перенаправление на страницу заявки
	http.Redirect(w, r, "/calculation/1", http.StatusSeeOther)
}

// POST: Логическое удаление заявки
func deleteCalculationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	// Логическое удаление через SQL UPDATE
	err = deleteCalculationSQL(id)
	if err != nil {
		http.Error(w, "Ошибка удаления заявки", http.StatusInternalServerError)
		return
	}
	
	// Перенаправление на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	// Инициализация подключений
	initDB()
	initMinIO()
	
	// Настройка роутера
	r := mux.NewRouter()
	
	// Статические файлы
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// Маршруты
	r.HandleFunc("/", instrumentsHandler).Methods("GET")
	r.HandleFunc("/instrument/{id}", instrumentDetailHandler).Methods("GET")
	r.HandleFunc("/calculation/{id}", calculationHandler).Methods("GET")
	r.HandleFunc("/add-to-calculation", addToCalculationHandler).Methods("POST")
	r.HandleFunc("/delete-calculation/{id}", deleteCalculationHandler).Methods("POST")
	
	fmt.Println("Сервер запущен на http://localhost:8081")
	fmt.Println("PostgreSQL подключен к:", dbHost+":"+strconv.Itoa(dbPort))
	fmt.Println("MinIO подключен к:", minioEndpoint)
	fmt.Println("Adminer доступен на: http://localhost:8080")
	fmt.Println("Доступные страницы:")
	fmt.Println("  GET / - Список инструментов")
	fmt.Println("  GET /instrument/{id} - Детали инструмента")
	fmt.Println("  GET /calculation/{id} - Просмотр заявки")
	fmt.Println("  POST /add-to-calculation - Добавить инструмент в заявку")
	fmt.Println("  POST /delete-calculation/{id} - Удалить заявку")
	
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}