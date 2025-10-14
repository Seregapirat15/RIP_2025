package handlers

import "math"

// CalculateExoplanetMass рассчитывает массу экзопланеты по формуле радиальной скорости
// Использует формулу: M_p = (M_s * K * P^(1/3)) / (2π * G)^(2/3) * sin(i)
// где:
// M_p - масса планеты
// M_s - масса звезды (в солнечных массах)
// K - амплитуда радиальной скорости (м/с)
// P - орбитальный период (дни)
// i - наклон орбиты (градусы)
func CalculateExoplanetMass(starMass, velocityAmplitude, orbitalPeriod, inclination float64) float64 {
	// Константы
	const G = 6.67430e-11 // гравитационная постоянная (м³/кг/с²)
	const solarMass = 1.989e30 // масса Солнца в кг
	const dayInSeconds = 86400.0 // секунд в дне
	
	// Переводим в СИ
	starMassKg := starMass * solarMass // масса звезды в кг
	periodSeconds := orbitalPeriod * dayInSeconds // период в секундах
	inclinationRad := inclination * math.Pi / 180.0 // наклон в радианах
	
	// Формула расчета массы планеты
	// Упрощенная версия для демонстрации
	// В реальности формула более сложная и зависит от точности инструмента
	
	// Коэффициент для перевода в массы Юпитера
	jupiterMass := 1.898e27 // масса Юпитера в кг
	
	// Расчет массы планеты (упрощенная формула)
	planetMass := (starMassKg * velocityAmplitude * math.Pow(periodSeconds, 1.0/3.0)) / 
		(math.Pow(2*math.Pi*G, 2.0/3.0) * math.Sin(inclinationRad))
	
	// Переводим в массы Юпитера
	planetMassJupiter := planetMass / jupiterMass
	
	// Ограничиваем результат разумными пределами
	if planetMassJupiter < 0.001 {
		planetMassJupiter = 0.001
	}
	if planetMassJupiter > 100 {
		planetMassJupiter = 100
	}
	
	return planetMassJupiter
}

// GetAstronomicalTelescopeInstrumentAccuracy возвращает точность инструмента для расчета
func GetAstronomicalTelescopeInstrumentAccuracy(instrumentType string) float64 {
	// Точность инструментов (коэффициент погрешности)
	accuracyMap := map[string]float64{
		"HARPS":     0.97,  // высокая точность
		"James Webb": 0.1,   // очень высокая точность
		"ESPRESSO":  0.1,    // очень высокая точность
		"SPIRou":    1.0,    // средняя точность
	}
	
	if accuracy, exists := accuracyMap[instrumentType]; exists {
		return accuracy
	}
	return 1.0 // стандартная точность
}

// ApplyAstronomicalTelescopeInstrumentCorrection применяет коррекцию на основе точности инструмента
func ApplyAstronomicalTelescopeInstrumentCorrection(calculatedMass, accuracy float64) float64 {
	// Применяем коррекцию точности
	// Чем выше точность (меньше значение), тем меньше погрешность
	correctionFactor := 1.0 + (1.0 - accuracy) * 0.1
	return calculatedMass * correctionFactor
}

