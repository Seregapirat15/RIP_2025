# Тестирование API endpoints

Write-Host "=== Тестирование API корзины ==="
try {
    $cart = Invoke-RestMethod -Uri "http://localhost:8081/api/orders/cart" -Method GET
    Write-Host "Корзина: $($cart | ConvertTo-Json)"
} catch {
    Write-Host "Ошибка получения корзины: $($_.Exception.Message)"
}

Write-Host "`n=== Тестирование добавления услуги ==="
$body = @{
    service_id = 1
    exoplanet_name = "Test Planet"
    star_mass = 1.0
    orbital_period = 365.0
    velocity_amplitude = 10.0
    inclination = 90.0
    eccentricity = 0.0
    comment = "Test comment"
    other_info = ""
} | ConvertTo-Json

try {
    $result = Invoke-RestMethod -Uri "http://localhost:8081/api/orders/services" -Method POST -Body $body -ContentType "application/json"
    Write-Host "Успешно добавлено: $($result | ConvertTo-Json)"
} catch {
    Write-Host "Ошибка добавления: $($_.Exception.Message)"
    Write-Host "Статус код: $($_.Exception.Response.StatusCode)"
}

Write-Host "`n=== Повторное добавление той же услуги (должна быть ошибка 409) ==="
try {
    $result = Invoke-RestMethod -Uri "http://localhost:8081/api/orders/services" -Method POST -Body $body -ContentType "application/json"
    Write-Host "Неожиданно успешно: $($result | ConvertTo-Json)"
} catch {
    Write-Host "Ожидаемая ошибка дубликата: $($_.Exception.Message)"
    Write-Host "Статус код: $($_.Exception.Response.StatusCode)"
}

Write-Host "`n=== Проверка обновленной корзины ==="
try {
    $cart = Invoke-RestMethod -Uri "http://localhost:8081/api/orders/cart" -Method GET
    Write-Host "Обновленная корзина: $($cart | ConvertTo-Json)"
} catch {
    Write-Host "Ошибка получения корзины: $($_.Exception.Message)"
}
