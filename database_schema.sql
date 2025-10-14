-- Создание базы данных для системы расчета массы экзопланет
-- Лабораторная работа 2
-- База данных уже создана через переменную окружения POSTGRES_DB

-- 1. Таблица астрономических инструментов (телескопы)
CREATE TABLE instruments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'ground' или 'space'
    description TEXT,
    accuracy DECIMAL(10,3),
    accuracy_unit VARCHAR(50),
    location VARCHAR(200),
    status VARCHAR(50) DEFAULT 'Активен',
    launch_date VARCHAR(10),
    measurement_range VARCHAR(100),
    resolution VARCHAR(100),
    calibration VARCHAR(200),
    stability VARCHAR(200),
    instrument_type VARCHAR(200),
    image_url VARCHAR(255), -- Nullable URL к изображению
    is_deleted BOOLEAN DEFAULT FALSE, -- статус удален/действует
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Таблица расчетов массы экзопланет (заявки)
CREATE TABLE calculations (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL, -- черновик, удалён, сформирован, завершён, отклонён
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- дата создания
    creator_id INTEGER NOT NULL DEFAULT 1, -- создатель (пока хардкод)
    
    -- Дополнительные поля (Nullable)
    formation_date TIMESTAMP, -- дата формирования (2 действия создателя)
    completion_date TIMESTAMP, -- дата завершения (2 действия модератора)
    moderator_id INTEGER DEFAULT 1, -- модератор (пока хардкод)
    
    -- Поля по предметной области
    researcher_name VARCHAR(100),
    institution VARCHAR(100),
    result TEXT,
    total_mass DECIMAL(15,6), -- рассчитывается при завершении заявки
    notes TEXT
);

-- 3. Таблица связи расчетов с инструментами (м-м)
CREATE TABLE calculation_instruments (
    calculation_id INTEGER NOT NULL REFERENCES calculations(id),
    instrument_id INTEGER NOT NULL REFERENCES instruments(id),
    
    -- Дополнительные поля м-м
    exoplanet_name VARCHAR(100) NOT NULL,
    star_mass DECIMAL(10,3) NOT NULL,
    orbital_period DECIMAL(10,3) NOT NULL,
    velocity_amplitude DECIMAL(10,3) NOT NULL,
    inclination DECIMAL(5,2) NOT NULL,
    comment TEXT,
    other_info VARCHAR(200),
    calculated_mass DECIMAL(15,6), -- рассчитывается при завершении
    
    -- Составной уникальный ключ
    PRIMARY KEY (calculation_id, instrument_id)
);

-- Создание индексов для оптимизации
CREATE INDEX idx_calculations_status ON calculations(status);
CREATE INDEX idx_calculations_creator ON calculations(creator_id);
CREATE INDEX idx_calculations_created_at ON calculations(created_at);
CREATE INDEX idx_instruments_deleted ON instruments(is_deleted);
CREATE INDEX idx_calculation_instruments_calc ON calculation_instruments(calculation_id);
CREATE INDEX idx_calculation_instruments_instr ON calculation_instruments(instrument_id);

-- Вставка тестовых данных

-- Инструменты
INSERT INTO instruments (name, full_name, type, description, accuracy, accuracy_unit, location, status, launch_date, measurement_range, resolution, calibration, stability, instrument_type, image_url) VALUES
('HARPS', 'HARPS(High Accuracy Radial velocity Planet Searcher)', 'ground', 'Один из самых продуктивных спектографов. Обнаружил множество экзопланет, включая планеты земного типа', 0.97, 'м/с (3.5 км/ч)', '3.6-м телескоп, Обсерватория Ла-Силья (ESO, Чили)', 'Активен', '2003', 'Видимый свет', 'Высокое (до 115 000)', 'Торий-аргонная лампа', 'Вакуумная камера с температурным контролем ±0.01 °С', 'Эшелле-спектрограф второго поколения', 'harps.jpg'),
('James Webb', 'James Webb Space Telescope (JWST)', 'space', 'Космический телескоп нового поколения для изучения экзопланет и их атмосфер', 0.1, 'м/с', 'Точка Лагранжа L2', 'Активен', '2021', 'Инфракрасный спектр', 'Очень высокое (до 2700)', 'Встроенные эталоны', 'Криогенное охлаждение до 7K', 'Космический инфракрасный телескоп', 'james_webb.jpg'),
('ESPRESSO', 'ESPRESSO (Echelle SPectrograph for Rocky Exoplanet and Stable Spectroscopic Observations)', 'ground', 'Сверхстабильный спектрограф для поиска землеподобных экзопланет', 0.1, 'м/с', 'VLT, Параналь, Чили', 'Активен', '2017', 'Видимый свет', 'Очень высокое (до 200 000)', 'Лазерный частотный гребенка', 'Температурная стабилизация ±0.001 °С', 'Эшелле-спектрограф третьего поколения', 'espresso.jpg'),
('SPIRou', 'SPIRou (SpectroPolarimètre InfraRouge)', 'ground', 'Спектрополяриметр для поиска экзопланет вокруг красных карликов', 1.0, 'м/с', 'CFHT, Мауна-Кеа, Гавайи', 'Активен', '2018', 'Ближний инфракрасный', 'Высокое (до 70 000)', 'Торий-аргонная лампа', 'Криогенное охлаждение', 'Инфракрасный спектрополяриметр', 'spirou.jpg');

-- Расчеты
INSERT INTO calculations (status, creator_id, researcher_name, institution, result) VALUES
('черновик', 1, 'Др. Анна Козлова', 'Институт астрономии РАН', 'Предварительный расчет массы завершен. Требуется дополнительная верификация.'),
('завершён', 1, 'Проф. Михаил Петров', 'МГУ', 'Расчет массы экзопланет завершен и подтвержден. Результаты опубликованы.');

-- Связи расчетов с инструментами
INSERT INTO calculation_instruments (calculation_id, instrument_id, exoplanet_name, star_mass, orbital_period, velocity_amplitude, inclination, comment, other_info) VALUES
(1, 1, 'Proxima Centauri b', 0.12, 11.18, 1.4, 90.0, 'Подтверждено HARPS', 'Жизнепригодная зона'),
(2, 2, 'WASP-18 b', 1.2, 0.94, 300.0, 86.0, 'Наблюдения JWST', 'Горячий Юпитер');

-- Обновление дат для завершенного расчета
UPDATE calculations SET 
    formation_date = '2024-02-15 10:30:00',
    completion_date = '2024-02-20 15:45:00',
    moderator_id = 1,
    total_mass = 10.3
WHERE id = 2;

-- Обновление рассчитанных масс в связях (результаты расчетов)
UPDATE calculation_instruments SET calculated_mass = 1.27 WHERE calculation_id = 1 AND instrument_id = 1;
UPDATE calculation_instruments SET calculated_mass = 10.3 WHERE calculation_id = 2 AND instrument_id = 2;

-- Добавим еще несколько инструментов в заявки для демонстрации
INSERT INTO calculation_instruments (calculation_id, instrument_id, exoplanet_name, star_mass, orbital_period, velocity_amplitude, inclination, comment, other_info, calculated_mass) VALUES
(1, 3, 'TRAPPIST-1 e', 0.08, 6.1, 0.6, 89.7, 'Наблюдения ESPRESSO', 'Потенциально обитаемая зона', 0.77),
(2, 4, 'GJ 1214 b', 0.15, 1.58, 12.5, 88.6, 'Наблюдения SPIRou', 'Суперземля', 6.55);