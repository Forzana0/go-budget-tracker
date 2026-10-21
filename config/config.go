package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config містить усі налаштування застосунку.
// Теги env вказують бібліотеці caarlos0/env, яку змінну середовища зчитувати.
type Config struct {
	ServerPort  string `env:"SERVER_PORT" envDefault:":8080"`
	DatabaseURL string `env:"DATABASE_URL,notEmpty"`
}

// Load завантажує змінні з .env файлу (якщо він є) та парсить їх у структуру Config.
func Load() (*Config, error) {
	// Спробуємо завантажити .env файл.
	// Якщо файл не знайдено (наприклад, у продакшен середовищі Docker/K8s),
	// ми не викидаємо паніку, а просто логуємо це, бо змінні можуть прийти з самої ОС.
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: No .env file found, reading configuration directly from system environment")
	}

	cfg := &Config{}

	// env.Parse аналізує змінні оточення системи та записує їх у поля структури за тегами.
	// Якщо поле має тег `notEmpty`, але змінна порожня або відсутня — повернеться зрозуміла помилка.
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
