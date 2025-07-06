// internal/config/config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `json:"server"`
	Ethereum EthereumConfig `json:"ethereum"`
	Merkle   MerkleConfig   `json:"merkle"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	CORS         bool   `json:"cors"`
}

// EthereumConfig holds Ethereum-related configuration
type EthereumConfig struct {
	RPCURL          string `json:"rpc_url"`
	PrivateKey      string `json:"private_key"`
	ContractAddress string `json:"contract_address"`
	TokenAddress    string `json:"token_address"`
	GasLimit        uint64 `json:"gas_limit"`
	GasPrice        int64  `json:"gas_price"`
}

// MerkleConfig holds Merkle tree configuration
type MerkleConfig struct {
	MaxClaims    int    `json:"max_claims"`
	WorkerCount  int    `json:"worker_count"`
	BatchSize    int    `json:"batch_size"`
	CacheEnabled bool   `json:"cache_enabled"`
	OutputFormat string `json:"output_format"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	File   string `json:"file"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
			CORS:         true,
		},
		Ethereum: EthereumConfig{
			RPCURL:   "http://localhost:8545",
			GasLimit: 3000000,
			GasPrice: 20000000000, // 20 gwei
		},
		Merkle: MerkleConfig{
			MaxClaims:    1000000,
			WorkerCount:  0, // 0 means use runtime.NumCPU()
			BatchSize:    1000,
			CacheEnabled: true,
			OutputFormat: "json",
		},
		Database: DatabaseConfig{
			Type:    "sqlite",
			Host:    "localhost",
			Port:    5432,
			Name:    "merkle_airdrop",
			SSLMode: "disable",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			File:   "airdrop.log",
		},
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		// If file doesn't exist, return default config
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(config *Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Merkle.MaxClaims <= 0 {
		return fmt.Errorf("max_claims must be positive")
	}

	if c.Merkle.BatchSize <= 0 {
		return fmt.Errorf("batch_size must be positive")
	}

	validFormats := map[string]bool{"json": true, "csv": true}
	if !validFormats[c.Merkle.OutputFormat] {
		return fmt.Errorf("invalid output format: %s", c.Merkle.OutputFormat)
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	switch c.Database.Type {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host, c.Database.Port, c.Database.User,
			c.Database.Password, c.Database.Name, c.Database.SSLMode)
	case "sqlite":
		return fmt.Sprintf("%s.db", c.Database.Name)
	default:
		return ""
	}
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
