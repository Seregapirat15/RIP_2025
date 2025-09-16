package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Структуры данных
type Instrument struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	FullName        string  `json:"full_name"`
	Type            string  `json:"type"`
	Description     string  `json:"description"`
	Accuracy        float64 `json:"accuracy"`
	AccuracyUnit    string  `json:"accuracy_unit"`
	Location        string  `json:"location"`
	Status          string  `json:"status"`
	LaunchDate      string  `json:"launch_date"`
	MeasurementRange string `json:"measurement_range"`
	Resolution      string  `json:"resolution"`
	Calibration     string  `json:"calibration"`
	Stability       string  `json:"stability"`
	InstrumentType  string  `json:"instrument_type"`
	Image           string  `json:"image"`
}

type CalculationItem struct {
	InstrumentID      int     `json:"instrument_id"`
	ExoplanetName     string  `json:"exoplanet_name"`
	StarMass          float64 `json:"star_mass"`
	OrbitalPeriod     float64 `json:"orbital_period"`
	VelocityAmplitude float64 `json:"velocity_amplitude"`
	Inclination       float64 `json:"inclination"`
	Comment           string  `json:"comment"`
	Other             string  `json:"other"`
}

type Calculation struct {
	ID             int               `json:"id"`
	ResearcherName string            `json:"researcher_name"`
	Institution    string            `json:"institution"`
	CreatedAt      string            `json:"created_at"`
	Status         string            `json:"status"`
	Items          []CalculationItem `json:"items"`
	Result         string            `json:"result"`
}

type Data struct {
	Instruments  []Instrument  `json:"instruments"`
	Calculations []Calculation `json:"calculations"`
}

var data Data
var minioClient *minio.Client

// Конфигурация MinIO
const (
	minioEndpoint        = "localhost:9000"
	minioAccessKeyID     = "minioadmin123"
	minioSecretAccessKey = "minioadmin123456"
	minioUseSSL          = false
	minioBucketName      = "telescope-images"
)

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

// Загрузка данных из JSON
func loadData() {
	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal("Ошибка чтения data.json:", err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal("Ошибка парсинга JSON:", err)
	}
}

// Получение инструмента по ID
func getInstrumentByID(id int) *Instrument {
	for _, instrument := range data.Instruments {
		if instrument.ID == id {
			return &instrument
		}
	}
	return nil
}

// Главная страница - список инструментов
func instrumentsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/instruments.html"))
	
	// Получение параметра поиска
	searchQuery := r.URL.Query().Get("search")
	
	// Подсчет общего количества расчетов
	totalCalculations := len(data.Calculations)
	
	// Получение выбранных инструментов (первые 2 для примера)
	var selectedInstruments []Instrument
	if len(data.Instruments) >= 2 {
		selectedInstruments = data.Instruments[:2]
	} else {
		selectedInstruments = data.Instruments
	}
	
	// Фильтрация инструментов по поисковому запросу
	var filteredInstruments []Instrument
	if searchQuery != "" {
		for _, instrument := range data.Instruments {
			if strings.Contains(strings.ToLower(instrument.Name), strings.ToLower(searchQuery)) ||
			   strings.Contains(strings.ToLower(instrument.FullName), strings.ToLower(searchQuery)) {
				filteredInstruments = append(filteredInstruments, instrument)
			}
		}
	} else {
		filteredInstruments = data.Instruments
	}
	
	// Генерация URL изображений для отфильтрованных инструментов
	instrumentsWithImages := make([]map[string]interface{}, len(filteredInstruments))
	for i, instrument := range filteredInstruments {
		imageURL, err := getImageURL(instrument.Image)
		if err != nil {
			log.Printf("Ошибка получения URL изображения для %s: %v", instrument.Name, err)
			imageURL = ""
		}
		
		instrumentsWithImages[i] = map[string]interface{}{
			"Instrument": instrument,
			"ImageURL":   imageURL,
		}
	}
	
	templateData := map[string]interface{}{
		"Instruments":        instrumentsWithImages,
		"TotalCalculations":  totalCalculations,
		"SelectedInstruments": selectedInstruments,
		"SearchQuery":        searchQuery,
	}
	
	err := tmpl.Execute(w, templateData)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		return
	}
}

// Страница расчета массы экзопланеты
func calculationHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/calculation.html"))
	
	// Получение параметра поиска
	searchQuery := r.URL.Query().Get("search")
	
	// Получение выбранных инструментов для расчета (HARPS и James Webb)
	var selectedInstruments []Instrument
	for _, instrument := range data.Instruments {
		if instrument.Name == "HARPS" || instrument.Name == "James Webb" {
			selectedInstruments = append(selectedInstruments, instrument)
		}
	}
	
	// Фильтрация выбранных инструментов по поисковому запросу
	var filteredSelectedInstruments []Instrument
	if searchQuery != "" {
		for _, instrument := range selectedInstruments {
			if strings.Contains(strings.ToLower(instrument.Name), strings.ToLower(searchQuery)) ||
			   strings.Contains(strings.ToLower(instrument.FullName), strings.ToLower(searchQuery)) {
				filteredSelectedInstruments = append(filteredSelectedInstruments, instrument)
			}
		}
	} else {
		filteredSelectedInstruments = selectedInstruments
	}
	
	// Генерация URL изображений для отфильтрованных выбранных инструментов
	selectedInstrumentsWithImages := make([]map[string]interface{}, len(filteredSelectedInstruments))
	for i, instrument := range filteredSelectedInstruments {
		imageURL, err := getImageURL(instrument.Image)
		if err != nil {
			log.Printf("Ошибка получения URL изображения для %s: %v", instrument.Name, err)
			imageURL = ""
		}
		
		selectedInstrumentsWithImages[i] = map[string]interface{}{
			"Instrument": instrument,
			"ImageURL":   imageURL,
		}
	}
	
	totalCalculations := len(data.Calculations)
	
	templateData := map[string]interface{}{
		"SelectedInstruments": selectedInstrumentsWithImages,
		"TotalCalculations":   totalCalculations,
		"SearchQuery":         searchQuery,
	}
	
	err := tmpl.Execute(w, templateData)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		return
	}
}

// Страница детальной информации об инструменте
func instrumentDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/instrument/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Поиск инструмента по ID
	instrument := getInstrumentByID(id)
	if instrument == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Генерация URL изображения
	imageURL, err := getImageURL(instrument.Image)
	if err != nil {
		log.Printf("Ошибка получения URL изображения для %s: %v", instrument.Name, err)
		imageURL = ""
	}
	
	tmpl := template.Must(template.ParseFiles("templates/instrument_detail.html"))
	
	totalCalculations := len(data.Calculations)
	templateData := map[string]interface{}{
		"Instrument":        instrument,
		"TotalCalculations": totalCalculations,
		"ImageURL":          imageURL,
	}
	
	err = tmpl.Execute(w, templateData)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Инициализация MinIO
	initMinIO()
	
	// Загрузка данных из JSON
	loadData()
	
	// Настройка статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// Маршруты (более специфичные маршруты должны быть первыми)
	http.HandleFunc("/calculation", calculationHandler)
	http.HandleFunc("/instrument/", instrumentDetailHandler)
	http.HandleFunc("/", instrumentsHandler)
	
	fmt.Println("Сервер запущен на http://localhost:8081")
	fmt.Println("MinIO подключен к:", minioEndpoint)
	fmt.Println("Bucket:", minioBucketName)
	fmt.Println("Доступные страницы:")
	fmt.Println("  GET / - Список инструментов")
	fmt.Println("  GET /calculation - Расчет массы экзопланеты")
	fmt.Println("  GET /instrument/{id} - Детали инструмента")
	
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
