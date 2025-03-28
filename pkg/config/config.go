package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for our application.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	Logging  LoggingConfig
	JWT      JWTConfig
}

// ServerConfig holds all server related configuration.
type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	Host         string
}

// DatabaseConfig holds all database related configuration.
type DatabaseConfig struct {
	InfluxURL      string
	InfluxOrg      string
	InfluxBucket   string
	InfluxToken    string
	InfluxPassword string
}

// RedisConfig holds all Redis related configuration.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// KafkaConfig holds all Kafka related configuration.
type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

// LoggingConfig holds all logging related configuration.
type LoggingConfig struct {
	Level string
	File  string
}

// JWTConfig holds all JWT related configuration.
type JWTConfig struct {
	Secret     string
	ExpireMins int
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	// Create a new viper instance
	v := viper.New()

	// Set default configurations
	setDefaults(v)

	// Set config file name and path
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	v.AddConfigPath(".")

	// Enable reading from environment variables
	// Environment variables take precedence over config files
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error if desired
		fmt.Println("Config file not found, using defaults and environment variables")
	}

	// Unmarshal the config into our Config struct
	config := &Config{}
	err = v.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server configs
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.readTimeout", 10)  // seconds
	v.SetDefault("server.writeTimeout", 10) // seconds
	v.SetDefault("server.host", "0.0.0.0")

	// Database configs
	v.SetDefault("database.influxURL", "http://localhost:8086")
	v.SetDefault("database.influxOrg", "analytics")
	v.SetDefault("database.influxBucket", "metrics")

	// Redis configs
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	// Kafka configs
	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.topic", "analytics")
	v.SetDefault("kafka.groupID", "analytics-group")

	// Logging configs
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.file", "")

	// JWT configs
	v.SetDefault("jwt.expireMins", 60)
}
