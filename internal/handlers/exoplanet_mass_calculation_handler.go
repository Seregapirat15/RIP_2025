package handlers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"

	"lab4/internal/database"
	"lab4/internal/models"
)

// getInstrumentAccuracy рассчитывает точность инструмента на основе его характеристик
func getInstrumentAccuracy(instrument *models.AstronomicalTelescopeInstrument) float64 {
	// Базовая точность на основе точности измерения скорости
	baseAccuracy := 1.0 - (instrument.VelocityPrecision / 10.0) // Нормализация к 0-1
	
	// Бонус за высокое спектральное разрешение
	resolutionBonus := 0.0
	if instrument.SpectralResolution > 100000 {
		resolutionBonus = 0.05 // +5% за очень высокое разрешение
	} else if instrument.SpectralResolution > 50000 {
		resolutionBonus = 0.03 // +3% за высокое разрешение
	}
	
	// Бонус за широкий диапазон длин волн
	wavelengthRange := instrument.WavelengthRangeMax - instrument.WavelengthRangeMin
	wavelengthBonus := 0.0
	if wavelengthRange > 1000 {
		wavelengthBonus = 0.02 // +2% за широкий диапазон
	}
	
	// Итоговая точность
	finalAccuracy := baseAccuracy + resolutionBonus + wavelengthBonus
	
	// Ограничиваем максимальную точность 99%
	if finalAccuracy > 0.99 {
		finalAccuracy = 0.99
	}
	
	// Преобразуем в проценты
	return finalAccuracy * 100
}

// calculateExoplanetMass рассчитывает массу экзопланеты с учетом характеристик инструмента
func calculateExoplanetMass(instrument *models.AstronomicalTelescopeInstrument, 
	velocityAmplitude, orbitalPeriod, starMass, inclination, eccentricity float64) float64 {
	
	// Константы
	const G = 6.67430e-11 // Гравитационная постоянная в м³/(кг·с²)
	const solarMass = 1.989e30 // Масса Солнца в кг
	const dayInSeconds = 86400 // Секунд в дне
	
	// Преобразуем единицы
	starMassKg := starMass * solarMass
	periodSeconds := orbitalPeriod * dayInSeconds
	
	// Формула расчета массы экзопланеты
	// M_planet = K * sqrt(1-e²) * (M* + M_planet)^(2/3) * (P/2πG)^(1/3) / sin(i)
	
	// Упрощенная формула для малых планет (M_planet << M*)
	sinInclination := math.Sin(inclination * math.Pi / 180.0)
	if sinInclination == 0 {
		sinInclination = 0.1 // Избегаем деления на ноль
	}
	
	// Корректируем амплитуду скорости с учетом точности инструмента
	correctedVelocity := velocityAmplitude * (1.0 - instrument.VelocityPrecision/100.0)
	
	// Расчет массы
	massFactor := correctedVelocity * math.Sqrt(1-eccentricity*eccentricity)
	periodFactor := math.Pow(periodSeconds/(2*math.Pi*G), 1.0/3.0)
	starMassFactor := math.Pow(starMassKg, 2.0/3.0)
	
	planetMass := massFactor * periodFactor * starMassFactor / sinInclination
	
	// Преобразуем в массы Солнца
	planetMassSolar := planetMass / solarMass
	
	return planetMassSolar
}

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
	
	// Рассчитываем общий результат
	var totalMass float64
	var averageAccuracy float64
	var instrumentCount int
	
	if len(instrumentsWithDetails) > 0 {
		for _, instrument := range instrumentsWithDetails {
			calcData := instrument["CalculationData"].(models.ExoplanetMassCalculationInstrument)
			instrumentData := instrument["Instrument"].(*models.AstronomicalTelescopeInstrument)
			
			// Рассчитываем массу с учетом характеристик инструмента
			calculatedMass := calculateExoplanetMass(
				instrumentData,
				calcData.VelocityAmplitude,
				calcData.OrbitalPeriod,
				calcData.StarMass,
				calcData.Inclination,
				0.0, // Эксцентриситет по умолчанию
			)
			
			totalMass += calculatedMass
			instrumentCount++
		}
		
		// Рассчитываем общую точность на основе характеристик инструментов
		if len(instrumentsWithDetails) > 0 {
			totalAccuracy := 0.0
			for _, instrument := range instrumentsWithDetails {
				instrumentData := instrument["Instrument"].(*models.AstronomicalTelescopeInstrument)
				// Используем реальную точность инструмента
				instrumentAccuracy := getInstrumentAccuracy(instrumentData)
				totalAccuracy += instrumentAccuracy
			}
			averageAccuracy = totalAccuracy / float64(len(instrumentsWithDetails))
		}
	}

	templateData := map[string]interface{}{
		"Calculation": calculation,
		"Instruments": instrumentsWithDetails,
		"CalculationResult": map[string]interface{}{
			"TotalMass":        totalMass,
			"InstrumentCount":  instrumentCount,
			"AverageAccuracy": averageAccuracy,
		},
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
	
	// Логическое удаление через курсор
	err = database.DeleteExoplanetMassCalculationWithCursor(id)
	if err != nil {
		http.Error(w, "Ошибка удаления заявки", http.StatusInternalServerError)
		return
	}
	
	// Перенаправление на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
