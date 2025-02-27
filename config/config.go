package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config - структура конфигурации приложения
type Config struct {
	App struct {
		Env  string `mapstructure:"env" validate:"required"`  // default, production, development
		Name string `mapstructure:"name" validate:"required"` // finance_bot
	}
	DB    DatabaseConfig `mapstructure:"db"`    // DatabaseConfig - структура конфигурации базы данных
	Redis RedisConfig    `mapstructure:"redis"` // RedisConfig - структура конфигурации Redis
}

// DatabaseConfig - структура конфигурации базы данных
type DatabaseConfig struct {
	URL       string `mapstructure:"url" validate:"required"`       // URL подключения к БД
	MaxConns  int    `mapstructure:"max_conns" validate:"required"` // Максимальное количество соединений
	IdleConns int    `mapstructure:"idle_conns"`                    // Количество простаивающих соединений
	Timeout   int    `mapstructure:"timeout"`                       // Таймаут подключения
}

// RedisConfig - структура конфигурации Redis
type RedisConfig struct {
	Addr     string `mapstructure:"addr" validate:"required"` // Адрес подключения к Redis
	Password string `mapstructure:"password"`                 // Пароль
	DB       int    `mapstructure:"db"`                       // Номер БД
	PoolSize int    `mapstructure:"pool_size"`                // Размер пула соединений
	Timeout  int    `mapstructure:"timeout"`                  // Таймаут подключения
}

// LoadConfig загружает конфигурацию с приоритетом:
//
// 1. Переменные окружения
// 2. Файл конфигурации в папке ./config/ с именем {env}.yaml
// 3. Значения по умолчанию
//
// Параметр path - путь к папке с конфигами
func LoadConfig(path string) (*Config, error) {
	// Переменные окружения
	viper.AutomaticEnv()

	// Значения по умолчанию
	viper.SetDefault("app.env", "default")
	viper.SetDefault("app.name", "finance_bot")
	viper.SetDefault("db.max_conns", 10)
	viper.SetDefault("db.idle_conns", 5)

	// Получаем окружение (из ENV или default)
	env := viper.GetString("app.env")
	fmt.Println("Loading configuration for environment:", env)

	viper.SetConfigName(env) // Например, local.yaml, production.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path) // Папка с конфигами (например, ./config/)

	// Читаем конфиг из файла (если найден)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found: %s.yaml, using ENV or defaults", env)
	} else {
		log.Println("Config loaded from:", viper.ConfigFileUsed())
	}

	// Парсим конфиг
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("не удалось распарсить конфиг: %w", err)
	}

	// Валидируем конфиг
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("конфиг не прошел валидацию: %w", err)
	}

	return &config, nil
}

// Validate проверяет конфигурацию на валидность
func (c *Config) Validate() error {
	if c.App.Env == "" {
		return fmt.Errorf("env не может быть пустым")
	}
	if c.App.Name == "" {
		return fmt.Errorf("name не может быть пустым")
	}
	if c.DB.URL == "" {
		return fmt.Errorf("db.url не может быть пустым")
	}
	if c.DB.MaxConns <= 0 {
		return fmt.Errorf("db.max_conns не может быть меньше или равно 0")
	}
	if c.DB.IdleConns < 0 {
		return fmt.Errorf("db.idle_conns не может быть меньше 0")
	}
	if c.DB.Timeout < 0 {
		return fmt.Errorf("db.timeout не может быть меньше 0")
	}
	if c.DB.Timeout == 0 {
		log.Println("db.timeout не указан, установлено значение по умолчанию: 30")
		c.DB.Timeout = 30
	}
	if c.Redis.Addr == "" {
		return fmt.Errorf("redis.addr не может быть пустым")
	}
	if c.Redis.PoolSize <= 0 {
		return fmt.Errorf("redis.pool_size не может быть меньше или равно 0")
	}
	if c.Redis.Timeout < 0 {
		log.Println("redis.timeout не может быть меньше 0, установлено значение по умолчанию: 30")
		return fmt.Errorf("redis.timeout не может быть меньше 0")
	}
	if c.Redis.Timeout == 0 {
		c.Redis.Timeout = 30
	}

	return nil
}
