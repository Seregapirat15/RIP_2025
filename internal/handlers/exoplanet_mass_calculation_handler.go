package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"lab2/internal/database"
)

// ExoplanetMassCalculationHandler обрабатывает страницу просмотра заявки
func ExoplanetMassCalculationHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL
	idStr := r.URL.Path[len("/calculation/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Получение заявки из БД
	calculation, err := database.GetExoplanetMassCalculationByID(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	// Получение инструментов заявки
	instruments, err := database.GetExoplanetMassCalculationInstruments(id)
	if err != nil {
		http.Error(w, "Ошибка получения инструментов заявки", http.StatusInternalServerError)
		return
	}
	
	// Получение полной информации об инструментах
	var instrumentsWithDetails []map[string]interface{}
	for _, calcInstrument := range instruments {
		instrument, err := database.GetAstronomicalTelescopeInstrumentByID(calcInstrument.InstrumentID)
		if err != nil {
			continue
		}
		
		var imageURL string
		if instrument.ImageURL != nil {
			imageURL, err = getTelescopeImageURL(*instrument.ImageURL)
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

// AddAstronomicalTelescopeInstrumentToExoplanetMassCalculationHandler обрабатывает добавление инструмента в заявку
func AddAstronomicalTelescopeInstrumentToExoplanetMassCalculationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	
	// Парсинг формы
	instrumentID, err := strconv.Atoi(r.FormValue("tool_id"))
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
	
	// Добавление в заявку (для демонстрации используем userID = 1)
	err = database.AddAstronomicalTelescopeInstrumentToExoplanetMassCalculation(1, instrumentID, exoplanetName, starMass, 
		orbitalPeriod, velocityAmplitude, inclination, "", "")
	if err != nil {
		http.Error(w, "Ошибка добавления в заявку", http.StatusInternalServerError)
		return
	}
	
	// Перенаправление на страницу заявки
	http.Redirect(w, r, "/calculation/1", http.StatusSeeOther)
}

// DeleteExoplanetMassCalculationHandler обрабатывает логическое удаление заявки
func DeleteExoplanetMassCalculationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	
	// Извлечение ID из URL
	idStr := r.URL.Path[len("/delete-calculation/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	// Логическое удаление через SQL UPDATE
	err = database.DeleteExoplanetMassCalculationSQL(id)
	if err != nil {
		http.Error(w, "Ошибка удаления заявки", http.StatusInternalServerError)
		return
	}
	
	// Перенаправление на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
