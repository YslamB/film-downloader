package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type EnvVar struct {
	Key          string
	Required     bool
	DefaultValue string
	Description  string
}

type ConfigLoader struct {
	prefix   string
	required []EnvVar
	optional []EnvVar
}

func NewConfigLoader(prefix string) *ConfigLoader {
	return &ConfigLoader{
		prefix:   prefix,
		required: make([]EnvVar, 0),
		optional: make([]EnvVar, 0),
	}
}

func (cl *ConfigLoader) AddRequired(key, description string) *ConfigLoader {
	cl.required = append(cl.required, EnvVar{
		Key:         key,
		Required:    true,
		Description: description,
	})
	return cl
}

func (cl *ConfigLoader) AddOptional(key, defaultValue, description string) *ConfigLoader {
	cl.optional = append(cl.optional, EnvVar{
		Key:          key,
		Required:     false,
		DefaultValue: defaultValue,
		Description:  description,
	})
	return cl
}

func LoadString(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

func LoadStringWithDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

func LoadInt(key string) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}

	parsedValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		log.Fatalf("Invalid integer value for %s: %v", key, err)
	}

	return int(parsedValue)
}

func LoadIntWithDefault(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s: %v, using default: %d", key, err, defaultValue)
		return defaultValue
	}

	return int(parsedValue)
}

func LoadBool(key string) bool {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Invalid boolean value for %s: %v", key, err)
	}

	return parsedValue
}

func LoadBoolWithDefault(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: Invalid boolean value for %s: %v, using default: %v", key, err, defaultValue)
		return defaultValue
	}

	return parsedValue
}

func (cl *ConfigLoader) ValidateConfiguration() error {
	var missingVars []string

	for _, envVar := range cl.required {
		key := cl.prefix + envVar.Key
		if value, exists := os.LookupEnv(key); !exists || value == "" {
			missingVars = append(missingVars, key)
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("Missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	return nil
}

func (cl *ConfigLoader) LoadConfiguration() map[string]string {
	cl.ValidateConfiguration()

	config := make(map[string]string)

	for _, envVar := range cl.required {
		key := cl.prefix + envVar.Key
		config[envVar.Key] = LoadString(key)
	}

	for _, envVar := range cl.optional {
		key := cl.prefix + envVar.Key
		config[envVar.Key] = LoadStringWithDefault(key, envVar.DefaultValue)
	}

	return config
}

func (cl *ConfigLoader) PrintConfiguration(config map[string]string, sensitiveKeys []string) {
	log.Println("Configuration loaded:")

	sensitiveMap := make(map[string]bool)
	for _, key := range sensitiveKeys {
		sensitiveMap[key] = true
	}

	for key, value := range config {
		if sensitiveMap[key] {
			log.Printf("  %s: ***masked***", key)
		} else {
			log.Printf("  %s: %s", key, value)
		}
	}
}

type ConfigValidator struct {
	validators map[string]func(string) error
}

func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		validators: make(map[string]func(string) error),
	}
}

func (cv *ConfigValidator) AddValidator(key string, validator func(string) error) {
	cv.validators[key] = validator
}

func (cv *ConfigValidator) ValidateConfig(config map[string]string) error {
	for key, validator := range cv.validators {
		if value, exists := config[key]; exists {
			if err := validator(value); err != nil {
				return WrapErrorf(err, "validation failed for %s", key)
			}
		}
	}
	return nil
}

func MustLoadConfig(loader *ConfigLoader) map[string]string {
	config := loader.LoadConfiguration()
	return config
}
