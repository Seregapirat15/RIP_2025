package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"lab4/internal/database"
)

// AstronomicalTelescopeInstrumentsHandler обрабатывает главную страницу со списком инструментов
func AstronomicalTelescopeInstrumentsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/instruments.html"))
	
	// Получение параметра поиска
	searchQuery := r.URL.Query().Get("search")
	
	// Получение инструментов из БД
	instruments, err := database.GetAstronomicalTelescopeInstruments(searchQuery)
	if err != nil {
		log.Printf("Ошибка получения астрономических инструментов: %v", err)
		http.Error(w, "Ошибка получения астрономических инструментов", http.StatusInternalServerError)
		return
	}
	
	log.Printf("Получено инструментов: %d", len(instruments))
	
	// Получение текущей заявки пользователя (для демонстрации используем userID = 1)
	currentCalc, err := database.GetCurrentExoplanetMassCalculation(1)
	if err != nil {
		log.Printf("Ошибка получения заявки: %v", err)
		// Продолжаем без заявки
		currentCalc = nil
	}
	
	var totalCalculations int
	if currentCalc != nil {
		calcInstruments, err := database.GetExoplanetMassCalculationInstruments(currentCalc.ID)
		if err == nil {
			totalCalculations = len(calcInstruments)
		}
	}
	
	// Генерация URL изображений
	instrumentsWithImages := make([]map[string]interface{}, len(instruments))
	for i, instrument := range instruments {
		var imageURL string
		if instrument.ImageURL != nil {
			imageURL, err = getTelescopeImageURL(*instrument.ImageURL)
			if err != nil {
				log.Printf("Ошибка получения URL изображения для %s: %v", instrument.Name, err)
				imageURL = ""
			}
		}
		
		instrumentsWithImages[i] = map[string]interface{}{
			"Instrument": instrument,
			"ImageURL":  imageURL,
		}
	}
	
	templateData := map[string]interface{}{
		"Instruments":        instrumentsWithImages,
		"TotalCalculations":  totalCalculations,
		"CurrentCalculation": currentCalc,
		"SearchQuery":        searchQuery,
	}
	
	err = tmpl.Execute(w, templateData)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		return
	}
}

// AstronomicalTelescopeInstrumentDetailHandler обрабатывает страницу детальной информации об инструменте
func AstronomicalTelescopeInstrumentDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL
	idStr := r.URL.Path[len("/tool/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Получение инструмента из БД
	instrument, err := database.GetAstronomicalTelescopeInstrumentByID(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Генерация URL изображения
	var imageURL string
	if instrument.ImageURL != nil {
		imageURL, err = getTelescopeImageURL(*instrument.ImageURL)
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

// getTelescopeImageURL получает URL изображения телескопа
func getTelescopeImageURL(objectName string) (string, error) {
	return database.GetTelescopeImagePresignedURL(objectName)
}
